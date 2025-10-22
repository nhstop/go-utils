package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func Hashpassword(password string, secret int) string {
	if password == "" {
		return ""
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), secret)
	return string(hashed)
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashText(value string, secret string) string {
	if value == "" {
		return ""
	}
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}
