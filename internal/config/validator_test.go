package config

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Servers: []Server{{
					Name:     "web1",
					Host:     "192.168.1.10",
					User:     "admin",
					Password: "secret",
				}},
				Actions: []Action{{
					Name:    "deploy",
					Command: "echo hello",
				}},
			},
			wantErr: false,
		},
		{
			name: "no servers",
			config: &Config{
				Servers: []Server{},
				Actions: []Action{{
					Name:    "deploy",
					Command: "echo hello",
				}},
			},
			wantErr: true,
		},
		{
			name: "server without name",
			config: &Config{
				Servers: []Server{{
					Host:     "192.168.1.10",
					User:     "admin",
					Password: "secret",
				}},
			},
			wantErr: true,
		},
		{
			name: "server without host",
			config: &Config{
				Servers: []Server{{
					Name:     "web1",
					User:     "admin",
					Password: "secret",
				}},
			},
			wantErr: true,
		},
		{
			name: "server without user",
			config: &Config{
				Servers: []Server{{
					Name:     "web1",
					Host:     "192.168.1.10",
					Password: "secret",
				}},
			},
			wantErr: true,
		},
		{
			name: "server without authentication",
			config: &Config{
				Servers: []Server{{
					Name: "web1",
					Host: "192.168.1.10",
					User: "admin",
				}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
