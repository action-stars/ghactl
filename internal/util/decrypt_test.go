package util

import (
	"testing"

	"github.com/matryer/is"
)

func TestDecrypt(t *testing.T) {
	t.Run("errors if value is too short", func(t *testing.T) {
		is := is.New(t)
		cipherText := [10]byte{}

		_, err := Decrypt("secret", cipherText[:])

		is.True(err != nil) // should error
	})

	t.Run("decrypts string", func(t *testing.T) {
		is := is.New(t)
		secret := "secret"
		input := "Hello world!"

		cipherText, err := Encrypt(secret, []byte(input))
		if err != nil {
			t.FailNow()
		}

		plainText, err := Decrypt(secret, cipherText)

		is.NoErr(err)                      // should not error
		is.Equal(string(plainText), input) // should match
	})
}
