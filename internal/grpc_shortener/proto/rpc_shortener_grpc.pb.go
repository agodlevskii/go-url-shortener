// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: proto/rpc_shortener.proto

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ShortServiceClient is the client API for ShortService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortServiceClient interface {
	GetFullUrl(ctx context.Context, in *FullUrlRequest, opts ...grpc.CallOption) (*FullUrlResponse, error)
	GetUserUrls(ctx context.Context, in *UserUrlsGetRequest, opts ...grpc.CallOption) (*UserUrlsGetResponse, error)
	DeleteUserUrls(ctx context.Context, in *UserUrlsDeleteRequest, opts ...grpc.CallOption) (*UserUrlsDeleteResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Shorten(ctx context.Context, in *ShortRequest, opts ...grpc.CallOption) (*ShortResponse, error)
	ShortenBatch(ctx context.Context, in *ShortBatchRequest, opts ...grpc.CallOption) (*ShortBatchResponse, error)
	Statistics(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error)
}

type shortServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShortServiceClient(cc grpc.ClientConnInterface) ShortServiceClient {
	return &shortServiceClient{cc}
}

func (c *shortServiceClient) GetFullUrl(ctx context.Context, in *FullUrlRequest, opts ...grpc.CallOption) (*FullUrlResponse, error) {
	out := new(FullUrlResponse)
	err := c.cc.Invoke(ctx, "/grpc_shortener.ShortService/GetFullUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortServiceClient) GetUserUrls(ctx context.Context, in *UserUrlsGetRequest, opts ...grpc.CallOption) (*UserUrlsGetResponse, error) {
	out := new(UserUrlsGetResponse)
	err := c.cc.Invoke(ctx, "/grpc_shortener.ShortService/GetUserUrls", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortServiceClient) DeleteUserUrls(ctx context.Context, in *UserUrlsDeleteRequest, opts ...grpc.CallOption) (*UserUrlsDeleteResponse, error) {
	out := new(UserUrlsDeleteResponse)
	err := c.cc.Invoke(ctx, "/grpc_shortener.ShortService/DeleteUserUrls", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/grpc_shortener.ShortService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortServiceClient) Shorten(ctx context.Context, in *ShortRequest, opts ...grpc.CallOption) (*ShortResponse, error) {
	out := new(ShortResponse)
	err := c.cc.Invoke(ctx, "/grpc_shortener.ShortService/Shorten", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortServiceClient) ShortenBatch(ctx context.Context, in *ShortBatchRequest, opts ...grpc.CallOption) (*ShortBatchResponse, error) {
	out := new(ShortBatchResponse)
	err := c.cc.Invoke(ctx, "/grpc_shortener.ShortService/ShortenBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortServiceClient) Statistics(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, "/grpc_shortener.ShortService/Statistics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortServiceServer is the server API for ShortService service.
// All implementations must embed UnimplementedShortServiceServer
// for forward compatibility
type ShortServiceServer interface {
	GetFullUrl(context.Context, *FullUrlRequest) (*FullUrlResponse, error)
	GetUserUrls(context.Context, *UserUrlsGetRequest) (*UserUrlsGetResponse, error)
	DeleteUserUrls(context.Context, *UserUrlsDeleteRequest) (*UserUrlsDeleteResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Shorten(context.Context, *ShortRequest) (*ShortResponse, error)
	ShortenBatch(context.Context, *ShortBatchRequest) (*ShortBatchResponse, error)
	Statistics(context.Context, *StatsRequest) (*StatsResponse, error)
	mustEmbedUnimplementedShortServiceServer()
}

// UnimplementedShortServiceServer must be embedded to have forward compatible implementations.
type UnimplementedShortServiceServer struct {
}

func (UnimplementedShortServiceServer) GetFullUrl(context.Context, *FullUrlRequest) (*FullUrlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFullUrl not implemented")
}
func (UnimplementedShortServiceServer) GetUserUrls(context.Context, *UserUrlsGetRequest) (*UserUrlsGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserUrls not implemented")
}
func (UnimplementedShortServiceServer) DeleteUserUrls(context.Context, *UserUrlsDeleteRequest) (*UserUrlsDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserUrls not implemented")
}
func (UnimplementedShortServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedShortServiceServer) Shorten(context.Context, *ShortRequest) (*ShortResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shorten not implemented")
}
func (UnimplementedShortServiceServer) ShortenBatch(context.Context, *ShortBatchRequest) (*ShortBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenBatch not implemented")
}
func (UnimplementedShortServiceServer) Statistics(context.Context, *StatsRequest) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Statistics not implemented")
}
func (UnimplementedShortServiceServer) mustEmbedUnimplementedShortServiceServer() {}

// UnsafeShortServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortServiceServer will
// result in compilation errors.
type UnsafeShortServiceServer interface {
	mustEmbedUnimplementedShortServiceServer()
}

func RegisterShortServiceServer(s grpc.ServiceRegistrar, srv ShortServiceServer) {
	s.RegisterService(&ShortService_ServiceDesc, srv)
}

func _ShortService_GetFullUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FullUrlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortServiceServer).GetFullUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_shortener.ShortService/GetFullUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortServiceServer).GetFullUrl(ctx, req.(*FullUrlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortService_GetUserUrls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserUrlsGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortServiceServer).GetUserUrls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_shortener.ShortService/GetUserUrls",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortServiceServer).GetUserUrls(ctx, req.(*UserUrlsGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortService_DeleteUserUrls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserUrlsDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortServiceServer).DeleteUserUrls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_shortener.ShortService/DeleteUserUrls",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortServiceServer).DeleteUserUrls(ctx, req.(*UserUrlsDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_shortener.ShortService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortService_Shorten_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortServiceServer).Shorten(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_shortener.ShortService/Shorten",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortServiceServer).Shorten(ctx, req.(*ShortRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortService_ShortenBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortServiceServer).ShortenBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_shortener.ShortService/ShortenBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortServiceServer).ShortenBatch(ctx, req.(*ShortBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortService_Statistics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortServiceServer).Statistics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc_shortener.ShortService/Statistics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortServiceServer).Statistics(ctx, req.(*StatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShortService_ServiceDesc is the grpc.ServiceDesc for ShortService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShortService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc_shortener.ShortService",
	HandlerType: (*ShortServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFullUrl",
			Handler:    _ShortService_GetFullUrl_Handler,
		},
		{
			MethodName: "GetUserUrls",
			Handler:    _ShortService_GetUserUrls_Handler,
		},
		{
			MethodName: "DeleteUserUrls",
			Handler:    _ShortService_DeleteUserUrls_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _ShortService_Ping_Handler,
		},
		{
			MethodName: "Shorten",
			Handler:    _ShortService_Shorten_Handler,
		},
		{
			MethodName: "ShortenBatch",
			Handler:    _ShortService_ShortenBatch_Handler,
		},
		{
			MethodName: "Statistics",
			Handler:    _ShortService_Statistics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/rpc_shortener.proto",
}
