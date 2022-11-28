package main

import (
	"context"
	"errors"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go-url-shortener/internal/apperrors"
	"go-url-shortener/internal/config"
	pb "go-url-shortener/internal/grpc_shortener/proto"
	"go-url-shortener/internal/services"
	"go-url-shortener/internal/storage"
)

type ShortenerServer struct {
	pb.UnimplementedShortServiceServer
	db  storage.Storager
	cfg *config.Config
}

func main() {
	cfg := config.New(config.WithEnv(), config.WithFlags(), config.WithFile())
	repo, err := getRepo(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	addr := strings.Replace(cfg.GetServerAddr(), ":8080", ":3200", 1)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterShortServiceServer(s, &ShortenerServer{
		UnimplementedShortServiceServer: pb.UnimplementedShortServiceServer{},
		db:                              repo,
		cfg:                             cfg,
	})

	if err = s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func (s *ShortenerServer) GetFullURL(ctx context.Context, in *pb.FullUrlRequest) (*pb.FullUrlResponse, error) {
	url, err := services.GetFullURL(ctx, s.db, in.GetId())
	return &pb.FullUrlResponse{Url: "http://" + s.cfg.GetBaseURL() + "/" + url.URL}, err
}

func (s *ShortenerServer) GetUserURLs(ctx context.Context, in *pb.UserUrlsGetRequest) (*pb.UserUrlsGetResponse, error) {
	urls, err := services.GetUserURLs(ctx, s.db, in.GetUid(), s.cfg.GetBaseURL())
	if err != nil {
		return &pb.UserUrlsGetResponse{UserUrls: []*pb.UserUrl{}}, err
	}

	urlRes := make([]*pb.UserUrl, 0, len(urls))
	for _, url := range urls {
		urlRes = append(urlRes, &pb.UserUrl{
			Short:    url.Short,
			Original: url.Original,
		})
	}
	return &pb.UserUrlsGetResponse{UserUrls: urlRes}, nil
}

func (s *ShortenerServer) DeleteUserURLs(ctx context.Context, in *pb.UserUrlsDeleteRequest) *pb.UserUrlsDeleteResponse {
	ps := s.cfg.GetPoolSize()
	pool := make(chan func(), ps)
	for i := 0; i < ps; i++ {
		go func() {
			for f := range pool {
				f()
			}
		}()
	}

	go func() {
		pool <- func() {
			if err := services.DeleteUserURLs(ctx, s.db, in.GetUid(), in.GetIds()); err != nil {
				log.Error(err)
			}
		}
	}()

	return &pb.UserUrlsDeleteResponse{Result: "Request accepted"}
}

func (s *ShortenerServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	if p := services.Ping(ctx, s.db); p {
		return &pb.PingResponse{Result: "DB is up and running"}, nil
	}
	return &pb.PingResponse{Result: "DB is unavailable"}, nil
}

func (s *ShortenerServer) Shorten(ctx context.Context, in *pb.ShortRequest) (*pb.ShortResponse, error) {
	url, _, err := services.GetShortURL(ctx, s.db, in.GetUri(), in.GetUid(), s.cfg.GetBaseURL())
	return &pb.ShortResponse{Result: url}, err
}

func (s *ShortenerServer) ShortenBatch(ctx context.Context, in *pb.ShortBatchRequest) (*pb.ShortBatchResponse, error) {
	reqData := make([]services.BatchOriginalData, 0, len(in.BatchOriginal))
	for _, req := range in.BatchOriginal {
		reqData = append(reqData, services.BatchOriginalData{
			CorrelationID: req.CorrelationId,
			OriginalURL:   req.OriginalUrl,
		})
	}

	urls, err := services.GetShortURLsFromBatch(ctx, s.db, reqData, in.GetUid(), s.cfg.GetBaseURL())
	if err != nil {
		return nil, err
	}

	resData := make([]*pb.BatchShort, 0, len(urls))
	for _, url := range urls {
		resData = append(resData, &pb.BatchShort{
			CorrelationId: url.CorrelationID,
			ShortUrl:      url.ShortURL,
		})
	}

	return &pb.ShortBatchResponse{BatchShort: resData}, nil
}

func (s *ShortenerServer) Statistics(ctx context.Context, in *pb.StatsRequest) (*pb.StatsResponse, error) {
	ip, ok := getIPFromContext(ctx)
	if !ok {
		return nil, errors.New(apperrors.UntrustedIP)
	}

	stats, err := services.GetStatistics(ctx, s.db, ip, s.cfg.GetTrustedSubnet())
	if err != nil {
		return nil, err
	}
	return &pb.StatsResponse{Stats: &pb.Stats{
		Urls:  uint32(stats.URLs),
		Users: uint32(stats.Users),
	}}, nil
}

func getIPFromContext(ctx context.Context) (string, bool) {
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		xForwardFor := headers.Get("X-Real-IP")
		if len(xForwardFor) > 0 && xForwardFor[0] != "" {
			ips := strings.Split(xForwardFor[0], ",")
			if len(ips) > 0 {
				clientIP := ips[0]
				return clientIP, true
			}
		}
	}

	return "", false
}

func getRepo(ctx context.Context, cfg *config.Config) (storage.Storager, error) {
	if cfg.GetDBURL() != "" {
		return storage.NewDBRepo(ctx, cfg.GetDBURL())
	}
	if cfg.GetStorageFileName() != "" {
		return storage.NewFileRepo(cfg.GetStorageFileName())
	}
	return storage.NewMemoryRepo(), nil
}
