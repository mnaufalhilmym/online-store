package random

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

func Numbers(n int) (string, error) {
	if n <= 0 {
		n = 1
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	min := 0
	if n-1 > 0 {
		min = int(math.Pow10(n - 1))
	}
	max := int(math.Pow10(n))

	return strconv.Itoa(rand.Intn(max-min) + min), nil
}
