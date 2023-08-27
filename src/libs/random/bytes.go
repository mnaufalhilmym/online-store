package random

import (
	"crypto/rand"
)

func Bytes(n int) ([]byte, error) {
	if n <= 0 {
		n = 1
	}

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return b, nil
}
