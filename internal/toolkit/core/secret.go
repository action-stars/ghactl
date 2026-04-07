package core

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/box"
)

// ErrInvalidPublicKey is returned when the provided public key is not 32 bytes long.
var ErrInvalidPublicKey = errors.New("invalid public key: must be 32 bytes")

// SetSecret sends an mask command to the workflow writer.
func SetSecret(w io.Writer, value string) error {
	c, err := NewCommand(MaskCmd, nil, value)
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
}

// EncryptSecret encrypts a secret value using the recipient's public key with a NaCl sealed box.
// This is compatible with GitHub's secret encryption which uses libsodium sealed boxes.
func EncryptSecret(publicKey []byte, value string) (string, error) {
	if len(publicKey) != 32 {
		return "", ErrInvalidPublicKey
	}

	var recipientPublicKey [32]byte
	copy(recipientPublicKey[:], publicKey)

	cipherText, err := box.SealAnonymous(nil, []byte(value), &recipientPublicKey, rand.Reader)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}
