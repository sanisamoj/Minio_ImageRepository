package main

import (
	"fmt"
	"math/rand"
)

func gen6DigCod() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}