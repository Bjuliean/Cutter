package random

import (
	"math/rand"
	"time"
)

func RandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("QWERTYUIOPASDFGHJKLZXCVBNM" +
		"qwertyuiopasdfghjklzxcvbnm" +
		"1234567890")

	b := make([]rune, size)

	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}
	return string(b)
}
