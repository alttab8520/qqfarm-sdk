package crypto

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tetratelabs/wazero"
)

func TestNewGatewayToken(t *testing.T) {
	tok := NewGatewayToken()
	if !strings.HasSuffix(tok, "=") {
		t.Fatalf("token %q", tok)
	}
	n := len(tok) - 1
	if n < 64 || n > 127 {
		t.Fatalf("len %d", n)
	}
	for _, c := range tok[:n] {
		if !strings.ContainsRune(tokenAlphabet, c) {
			t.Fatalf("bad char %q", c)
		}
	}
}

func TestIdentitySeal(t *testing.T) {
	plain := []byte("hello")
	out, tok, err := Identity{}.Seal(plain)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "hello" {
		t.Fatalf("body %q", out)
	}
	if !strings.HasSuffix(tok, "=") {
		t.Fatalf("token %q", tok)
	}
}

func TestOpenMissingWASM(t *testing.T) {
	t.Setenv("FARM_WASM", filepath.Join(t.TempDir(), "missing.wasm"))
	_, err := Open()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "FARM_WASM") {
		t.Fatalf("%v", err)
	}
}

func TestHostModuleInstantiates(t *testing.T) {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	t.Cleanup(func() { _ = r.Close(ctx) })
	rt := &Runtime{ctx: ctx, dataDir: t.TempDir(), appID: "wx", device: "PC;", version: officialVersion}
	if err := rt.instantiateHost(r); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeRoundtrip(t *testing.T) {
	path := os.Getenv("FARM_WASM")
	if path == "" {
		for _, p := range []string{
			filepath.Join("data", "tsdk-v3.8.2.wasm"),
			"tsdk-v3.8.2.wasm",
		} {
			if st, err := os.Stat(p); err == nil && st.Size() > 0 {
				path = p
				break
			}
		}
	}
	if path == "" {
		t.Skip("没有 tsdk-v3.8.2.wasm，跳过真实加密测试")
	}
	rt, err := NewRuntime(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = rt.Close() })

	plain := []byte("farm-seal-test")
	sealed, tok, err := rt.Seal(plain)
	if err != nil {
		t.Fatal(err)
	}
	if len(sealed) != len(plain) {
		t.Fatalf("len %d != %d", len(sealed), len(plain))
	}
	if string(sealed) == string(plain) {
		t.Fatal("密文与明文相同")
	}
	if !strings.HasSuffix(tok, "=") {
		t.Fatalf("token %q", tok)
	}
	if err := rt.BindUser("oTestOpenId"); err != nil {
		t.Fatal(err)
	}
	_, cred, err := rt.Seal([]byte{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(cred) < 140 {
		t.Fatalf("初始化凭据太短: %d", len(cred))
	}
}
