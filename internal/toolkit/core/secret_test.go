package core

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestSetSecret(t *testing.T) {
	t.Run("writes the mask command", func(t *testing.T) {
		is := is.New(t)
		value := "value"

		var b bytes.Buffer
		err := SetSecret(&b, value)

		is.NoErr(err)                                                   // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::%s\n", MaskCmd, value)) // should be equal
	})
}

func TestEncryptSecret(t *testing.T) {
	t.Run("encrypts a secret", func(t *testing.T) {
		is := is.New(t)
		k := "secret"
		v := "Hello world!"

		cipherText, err := EncryptSecret(k, v)

		is.NoErr(err) // should not error

		if err == nil {
			plainText, _ := DecryptSecret(k, cipherText)

			is.Equal(string(plainText), v) // should be equal
		}
	})
}

func TestDecryptSecret(t *testing.T) {
	t.Run("decrypts a secret", func(t *testing.T) {
		is := is.New(t)
		k := "secret"
		v := "Hello world!"

		cipherText, err := EncryptSecret(k, v)
		if err != nil {
			t.FailNow()
		}

		plainText, err := DecryptSecret(k, cipherText)

		is.NoErr(err)                  // should not error
		is.Equal(string(plainText), v) // should be equal
	})
}
