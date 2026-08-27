package client

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alttab8520/qqfarm-sdk/internal/game"
)

func TestPutSocialValidates(t *testing.T) {
	c := &Client{loggedIn: true}
	ctx := context.Background()
	if _, err := c.PutSocial(ctx, game.PutSocialIn{ItemID: 5005}); err == nil || !strings.Contains(err.Error(), "host_gid") {
		t.Fatalf("need host: %v", err)
	}
	if _, err := c.PutSocial(ctx, game.PutSocialIn{HostGID: 9}); err == nil || !strings.Contains(err.Error(), "item_id") {
		t.Fatalf("need item: %v", err)
	}
	if _, err := c.PutSocial(ctx, game.PutSocialIn{HostGID: 9, ItemID: 5006}); err == nil || !strings.Contains(err.Error(), "land_ids") {
		t.Fatalf("cloud needs land: %v", err)
	}
}

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
