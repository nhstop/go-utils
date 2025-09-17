package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"github.com/busnosh/go-utils/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func GetBcryptCost(secretKey string) int {
	cost := utils.GetEnv(secretKey, "BCRYPT_COST")
	num, err := strconv.Atoi(cost)
	if err != nil {
		return bcrypt.DefaultCost
	}
	return num
}

func Hashpassword(Password string, secretKey string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(Password), GetBcryptCost(secretKey))
	return string(hashed)
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HashText(value string, secretKey string) string {
	secret := utils.GetEnv(secretKey, "HASH_SECRET")
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}
