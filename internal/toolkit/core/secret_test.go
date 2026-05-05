package core

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/matryer/is"
	"golang.org/x/crypto/nacl/box"
)

func TestSetSecret(t *testing.T) {
	t.Run("writes_the_mask_command", func(t *testing.T) {
		is := is.New(t)
		value := "value"

		var b bytes.Buffer
		err := SetSecret(&b, value)

		is.NoErr(err)                                                   // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::%s\n", MaskCmd, value)) // should be equal
	})
}

func TestEncryptSecret(t *testing.T) {
	t.Run("encrypts_a_secret", func(t *testing.T) {
		is := is.New(t)
		v := "Hello world!"

		publicKey, privateKey, err := box.GenerateKey(rand.Reader)
		is.NoErr(err) // should generate key pair

		cipherText, err := EncryptSecret(publicKey[:], v)

		is.NoErr(err) // should not error

		decoded, err := base64.StdEncoding.DecodeString(cipherText)
		is.NoErr(err) // should decode base64

		plainText, ok := box.OpenAnonymous(nil, decoded, publicKey, privateKey)
		is.True(ok) // should open anonymous box
		is.True(ok) // should be ok

		is.Equal(string(plainText), v) // should be equal
	})
}
