package web

import (
	"testing"
	"time"
)

func TestConfig_LoadValid(t *testing.T) {
	t.Parallel()
	cfg := Config{}
	cfg.Load("tests/config.json")

	cfg.Port = "a"
	if cfg.Validate() != ErrInvalidPort {
		t.Error("Port validation failed")
	}
	
	cfg.Port = "65536"
	if cfg.Validate() != ErrInvalidPort {
		t.Error("Port validation failed")
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Parallel()
	type fields struct {
		Host               string
		Port               string
		CertFile           string
		KeyFile            string
		HTTPSPort          string
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		InsecureSkipVerify bool
		ShutdownTimeout    time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "invalid port",
			fields: fields{
				Port: "-12",
			},
			wantErr: true,
		},
		{
			name: "valid port",
			fields: fields{
				Port: "9000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Host:               tt.fields.Host,
				Port:               tt.fields.Port,
				CertFile:           tt.fields.CertFile,
				KeyFile:            tt.fields.KeyFile,
				HTTPSPort:          tt.fields.HTTPSPort,
				ReadTimeout:        tt.fields.ReadTimeout,
				WriteTimeout:       tt.fields.WriteTimeout,
				InsecureSkipVerify: tt.fields.InsecureSkipVerify,
				ShutdownTimeout:    tt.fields.ShutdownTimeout,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
