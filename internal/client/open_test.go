package client

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func TestNewLoginRequiresWASM(t *testing.T) {
	t.Setenv("FARM_WASM", filepath.Join(t.TempDir(), "missing.wasm"))
	s := New()
	_, err := s.Login(context.Background(), game.LoginIn{Code: "x"})
	if err == nil {
		t.Fatal("expected missing runtime error")
	}
	if !strings.Contains(err.Error(), "FARM_WASM") {
		t.Fatalf("%v", err)
	}
}
