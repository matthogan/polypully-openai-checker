package cache

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/codejago/polypully-openai-checker/internal/config"
	"hash"
)

const (
	// KeyHashAlgorithmSHA256 is the SHA256 hash algorithm
	KeyHashAlgorithmSHA256 = "sha256"
)

// NewKeySerializer creates a new serializer based on the configuration
// If the configuration specifies sha256, a SHA256 serializer is created,
// otherwise an MD5 serializer is created
func newKeySerializer(config *config.Config) keySerializer {
	if config.Cache.KeyHashAlgorithm == KeyHashAlgorithmSHA256 {
		return &sha256Serializer{}
	}
	return &md5Serializer{}
}

type keySerializer interface {
	serialise(key string) (string, error)
}

// serializer that uses the SHA256 hash algorithm

type sha256Serializer struct {
}

func (s *sha256Serializer) serialise(key string) (string, error) {
	return serialise(sha256.New(), key) // new one each time is simple for concurrency, however it's inefficient
}

// serializer that uses the MD5 hash algorithm

type md5Serializer struct {
}

func (s *md5Serializer) serialise(key string) (string, error) {
	return serialise(md5.New(), key)
}

func serialise[T hash.Hash](hasher T, key string) (string, error) {
	_, err := hasher.Write([]byte(key))
	if err != nil {
		return "", err
	}
	encoded := hex.EncodeToString(hasher.Sum(nil))
	return encoded, nil
}
