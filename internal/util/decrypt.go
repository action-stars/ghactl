package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
)

// Decrypt decrypts an encrypted value using the secret used to encrypt it.
func Decrypt(secret string, cipherText []byte) ([]byte, error) {
	secretKey := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(secretKey[:])
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}
