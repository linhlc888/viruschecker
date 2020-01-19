package main

import (
	"math/rand"
	"time"
)

var seededRand *rand.Rand

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

//RandomStringWithCharset returns rand string with length and charset
func RandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

//RandomString get random string from alphabet characters
func RandomString(length int) string {
	return RandomStringWithCharset(length, charset)
}
