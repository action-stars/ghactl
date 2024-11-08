package util

import (
	"math/rand"
)

// GenerateRandomString generates a random string of a length l using the charset [a-zA-Z0-9].
// The returned string will have a length of l and will not be terminated by a new line.
func GenerateRandomString(l int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, l)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
