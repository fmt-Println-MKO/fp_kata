package utils

import "math/rand"

func GenerateNewId() int {
	return rand.Intn(100) + 1
}
