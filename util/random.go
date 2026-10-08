package util

import (
	"crypto/rand"
	"math/big"
)

func RandomNumber6() (int64, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return 0, err
	}
	digit := num.Int64() + 100000
	return digit, nil
}
