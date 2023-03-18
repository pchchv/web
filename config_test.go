package web

import (
	"testing"
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
