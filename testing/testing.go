package testing

import (
	"math/rand"
	"time"
)

var (
	s       = rand.NewSource(time.Now().UnixNano())
	r       = rand.New(s)
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// RandStr generates a random string of n characters
func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
