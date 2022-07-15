package crypt

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/go-crypt/x/scrypt"
)

const (
	hashScryptDefaultRounds      = 16
	hashScryptDefaultBlockSize   = 8
	hashScryptDefaultParallelism = 1
)

type ScryptHash struct {
	rounds, blockSize int
	parallelism       uint8

	keySize, saltSize uint32

	salt string
}

// WithKeySize adjusts the key size of the resulting Scrypt hash. Default is 32.
func (b *ScryptHash) WithKeySize(size uint32) *ScryptHash {
	b.keySize = size

	return b
}

// WithSaltSize adjusts the salt size of the resulting Scrypt hash. Default is 16.
func (b *ScryptHash) WithSaltSize(size uint32) *ScryptHash {
	b.saltSize = size

	return b
}

// WithSalt sets the salt of the resulting Scrypt hash. Default is generated by crypto/rand.
func (b *ScryptHash) WithSalt(salt string) *ScryptHash {
	b.salt = salt

	return b
}

// WithLN sets the ln parameter (logN) of the resulting Scrypt hash. Default is 16.
func (b *ScryptHash) WithLN(rounds int) *ScryptHash {
	b.rounds = rounds

	return b
}

// WithR sets the r parameter (block size) of the resulting Scrypt hash. Default is 8.
func (b *ScryptHash) WithR(blockSize int) *ScryptHash {
	b.blockSize = blockSize

	return b
}

// WithP sets the p parameter (parallelism) of the resulting Scrypt hash. Default is 1.
func (b *ScryptHash) WithP(parallelism uint8) *ScryptHash {
	b.parallelism = parallelism

	return b
}

// Build checks the options are all configured correctly, setting defaults as necessary, calculates the password hash,
// and returns the Argon2id hash.
func (b ScryptHash) Build(password string) (h *ScryptDigest, err error) {
	h = &ScryptDigest{
		ln: b.rounds,
		r:  b.blockSize,
		p:  int(b.parallelism),
		k:  int(b.keySize),
	}

	if h.ln <= 0 {
		h.ln = hashScryptDefaultRounds
	}

	if h.r <= 0 {
		h.r = hashScryptDefaultBlockSize
	}

	if h.p <= 0 {
		h.p = hashScryptDefaultParallelism
	}

	if h.k == 0 {
		h.k = defaultKeySize
	}

	if b.salt != "" {
		if h.salt, err = b64rs.DecodeString(b.salt); err != nil {
			return nil, fmt.Errorf("error decoding Password salt from base64: %w", err)
		}
	} else {
		var (
			size = b.saltSize
		)

		if size <= 0 {
			size = defaultSaltSize
		}

		h.salt = make([]byte, size)

		if _, err = io.ReadFull(rand.Reader, h.salt); err != nil {
			return nil, fmt.Errorf("error reading random bytes for the salt: %w", err)
		}
	}

	if h.key, err = scrypt.Key([]byte(password), h.salt, h.n(), h.r, h.p, h.k); err != nil {
		return nil, fmt.Errorf("error calculating hash: %w", err)
	}

	return h, nil
}
