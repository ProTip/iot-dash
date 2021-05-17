package main

import (
	"crypto/rand"
	"encoding/base64"
)

/* Generates a secure base64 encoded token with 128 bits of entropy. */
func GetSecureToken() string {
	c := make([]byte, 16)
	_, err := rand.Read(c)
	if err != nil {
		/* Possible EAGAIN/entropy issue */
		panic("Error reading random bytes!")
	}

	return base64.StdEncoding.EncodeToString(c)
}
