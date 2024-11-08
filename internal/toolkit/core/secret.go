package core

import (
	"encoding/base64"
	"io"

	"github.com/action-stars/ghactl/internal/util"
)

// SetSecret sends an mask command to the workflow writer.
func SetSecret(w io.Writer, value string) error {
	c, err := NewCommand(MaskCmd, nil, value)
	if err != nil {
		return err
	}

	return IssueCommand(w, c)
}

// EncryptSecret encrypts a secret value so it isn't masked.
func EncryptSecret(key, value string) (string, error) {
	cipherText, err := util.Encrypt(key, []byte(value))
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(cipherText), nil
}

// DecryptSecret decrypts an encrypted secret value.
func DecryptSecret(key, secretValue string) (string, error) {
	cipherText, err := base64.RawStdEncoding.DecodeString(secretValue)
	if err != nil {
		return "", err
	}

	plainText, err := util.Decrypt(key, cipherText)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
