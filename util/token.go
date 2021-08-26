package util

import (
	"crypto/sha256"
	"strings"
	"time"
)

// Token utils
type Token struct {
	seed string
}

// Create a Token util instance.
//
// seed {string} - seed value.
func New(seed string) *Token {
	return &Token{
		seed: seed,
	}
}

// Create a Token util instance.
// seed from current time.
func NewDateSeed() *Token {
	return &Token{
		seed: time.Now().String(),
	}
}

// Add seed value,
func (t *Token) AddSeed(seed string) *Token {
	return &Token{
		seed: strings.Join([]string{t.seed, seed}, ""),
	}
}

func (t *Token) hash() [32]byte {
	return sha256.Sum256([]byte(t.seed))
}

// Create token.
func (t *Token) Create() string {
	hash := t.hash()
	return string(hash[:])
}

// Create token.
// specify length
func (t *Token) CreateSpecifyLength(length int) string {
	hash := t.hash()
	return string(hash[0:length])
}
