package main

import (
	"github.com/stretchr/testify/assert"
	"go-url-shortener/internal/testhelp"
	"reflect"
	"testing"
)

func Test_getServerAddress(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
		args    struct {
			addr string
		}
	}{
		{
			name: "Missing env variable",
			want: "localhost:8080",
		},
		{
			name: "Existing env variable",
			args: struct{ addr string }{addr: "testserver.com"},
			want: "testserver.com",
		},
	}

	testhelp.RemoveEnvVar(addrKey)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.addr != "" {
				testhelp.SetEnvVar(addrKey, tt.args.addr)
			}

			got, err := getServerAddress()
			if (err != nil) != tt.wantErr {
				t.Errorf("getServerAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getServerAddress() got = %v, want %v", got, tt.want)
			}

			testhelp.RemoveEnvVar(addrKey)
		})
	}
}

func Test_getRepo(t *testing.T) {
	tests := []struct {
		name string
		want string
		args struct{ filename string }
	}{
		{
			name: "Missing env variable",
			want: "MemoRepo",
		},
		{
			name: "Existing env variable",
			args: struct{ filename string }{filename: "testfile"},
			want: "FileRepo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.filename != "" {
				testhelp.SetEnvVar(storageFileName, tt.args.filename)
			}

			got, _ := getRepo()
			assert.Equal(t, tt.want, reflect.TypeOf(got).Name())

			if tt.args.filename != "" {
				testhelp.RemoveEnvVar(storageFileName)
			}
		})
	}
}
