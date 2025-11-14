package utils

import (
	"errors"
	"math/rand/v2"
	"strconv"
	"strings"
)

var (
	ErrInvalidLength = errors.New("invalid code length")
)

func GenerateNumericCode(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	var numericCode strings.Builder

	for range length {
		numericCode.WriteString(strconv.Itoa(rand.IntN(10)))
	}

	return numericCode.String(), nil
}
