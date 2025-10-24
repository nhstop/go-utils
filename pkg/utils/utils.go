package utils

import (
	"crypto/rand"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GenerateOTP(length int) (string, error) {
	max := big.NewInt(10)
	var sb strings.Builder

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		sb.WriteString(strconv.FormatInt(n.Int64(), 10))
	}
	return sb.String(), nil
}

func Map[T any, R any](arr []T, fn func(T) R) []R {
	result := make([]R, len(arr))
	for i, v := range arr {
		result[i] = fn(v)
	}
	return result
}
