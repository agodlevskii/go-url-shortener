package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_setBaseURL(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want string
	}{
		{
			name: "Config value is missing",
			want: "http://localhost:8080",
		},
		{
			name: "Config value is present",
			val:  "https://base.url",
			want: "https://base.url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val != "" {
				config.baseURL = tt.val
			}

			setBaseURL()
			assert.Equal(t, tt.want, config.baseURL)
		})
	}
}

func Test_setServerAddress(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want string
	}{
		{
			name: "Config value is missing",
			want: "localhost:8080",
		},
		{
			name: "Config value is present",
			val:  "base.url",
			want: "base.url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val != "" {
				config.addr = tt.val
			}

			setServerAddress()
			assert.Equal(t, tt.want, config.addr)
		})
	}
}

func Test_setFilename(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want string
	}{
		{
			name: "Config value is missing",
			want: "",
		},
		{
			name: "Config value is present",
			val:  "teststorage",
			want: "teststorage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val != "" {
				config.filename = tt.val
			}

			setFilename()
			assert.Equal(t, tt.want, config.filename)
		})
	}
}
