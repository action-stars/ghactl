package util

import (
	"testing"

	"github.com/matryer/is"
)

func TestEncrypt(t *testing.T) {
	t.Run("encrypts string", func(t *testing.T) {
		is := is.New(t)
		secret := "secret"
		input := "Hello world!"

		cipherText, err := Encrypt(secret, []byte(input))

		is.NoErr(err) // should not error

		if err == nil {
			plainText, err := Decrypt(secret, cipherText)
			if err != nil {
				t.FailNow()
			}

			is.Equal(string(plainText), input) // should match
		}
	})
}
