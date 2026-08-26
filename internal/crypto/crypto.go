package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

const tokenAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Encryptor seals RPC bodies and issues gateway tokens.
// Seal is atomic so the one-shot post-login credential cannot be stolen by a concurrent call.
type Encryptor interface {
	Seal(body []byte) (sealed []byte, token string, err error)
	BindUser(openID string) error
	Close() error
}

type Identity struct{}

func (Identity) Seal(body []byte) ([]byte, string, error) {
	if body == nil {
		body = []byte{}
	}
	return body, NewGatewayToken(), nil
}
func (Identity) BindUser(string) error { return nil }
func (Identity) Close() error          { return nil }

func NewGatewayToken() string {
	n := 64 + intn(64)
	buf := make([]byte, n)
	for i := range buf {
		buf[i] = tokenAlphabet[intn(len(tokenAlphabet))]
	}
	return string(buf) + "="
}

func intn(max int) int {
	if max <= 0 {
		return 0
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(n.Int64())
}

// Open loads the official encrypt runtime. Missing wasm is a hard error —
// falling back to plaintext would talk to the live gate and fail opaquely.
func Open() (Encryptor, error) {
	path, err := findWASM()
	if err != nil {
		return nil, err
	}
	return NewRuntime(path)
}

func findWASM() (string, error) {
	if p := os.Getenv("FARM_WASM"); p != "" {
		if st, err := os.Stat(p); err == nil && st.Size() > 0 {
			return p, nil
		}
		return "", fmt.Errorf("FARM_WASM 无效: %s", p)
	}
	cwd, _ := os.Getwd()
	exe, _ := os.Executable()
	var candidates []string
	for _, dir := range []string{cwd, filepath.Dir(exe)} {
		if dir == "" {
			continue
		}
		candidates = append(candidates,
			filepath.Join(dir, "data", "tsdk-v3.9.0.wasm"),
			filepath.Join(dir, "data", "tsdk.wasm"),
			filepath.Join(dir, "tsdk-v3.9.0.wasm"),
			filepath.Join(dir, "tsdk.wasm"),
		)
	}
	for _, p := range candidates {
		if st, err := os.Stat(p); err == nil && st.Size() > 0 {
			return p, nil
		}
	}
	return "", fmt.Errorf("未找到加密运行时，请将 tsdk-v3.9.0.wasm 放到 data/ 或设置 FARM_WASM")
}
