package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomString generates a random string with n-length
func RandomString(n int) string {
	var sb strings.Builder
	len := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomNumber generates a int64 random number between 0 and max
func RandomNumber(max int64) int64 {
	return rand.Int63n(max)
}
