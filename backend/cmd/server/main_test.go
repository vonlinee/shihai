package main

import (
	"testing"

	"shihai/internal/config"
)

func TestApplyServerPortOverride(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		wantPort string
		wantErr  bool
	}{
		{name: "not provided", port: "", wantPort: "8080"},
		{name: "valid override", port: "9090", wantPort: "9090"},
		{name: "normalizes leading zeroes", port: "08080", wantPort: "8080"},
		{name: "zero", port: "0", wantPort: "8080", wantErr: true},
		{name: "negative", port: "-1", wantPort: "8080", wantErr: true},
		{name: "too large", port: "65536", wantPort: "8080", wantErr: true},
		{name: "not numeric", port: "http", wantPort: "8080", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Server: config.ServerConfig{Port: "8080"}}

			err := applyServerPortOverride(cfg, tt.port)
			if (err != nil) != tt.wantErr {
				t.Fatalf("applyServerPortOverride() error = %v, wantErr %v", err, tt.wantErr)
			}
			if cfg.Server.Port != tt.wantPort {
				t.Fatalf("server port = %q, want %q", cfg.Server.Port, tt.wantPort)
			}
		})
	}
}
