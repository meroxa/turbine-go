package internal

import (
	"math/rand"
)

func randID() int {
	return rand.Intn(10000)
}
