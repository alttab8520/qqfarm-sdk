package crypto

import (
	"crypto/rand"
	"math/big"
)

const tokenAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Encryptor seals an RPC body. Production traffic needs the official runtime;
// tests can use Identity.
type Encryptor interface {
	Encrypt(body []byte) ([]byte, error)
	Token() string
}

type Identity struct{}

func (Identity) Encrypt(body []byte) ([]byte, error) { return body, nil }

func (Identity) Token() string { return NewGatewayToken() }

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
