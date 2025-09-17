package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/busnosh/go-utils/pkg/utils"
)

func getAESKey(secretKey string) ([]byte, error) {
	key := utils.GetEnv(secretKey, "ENCRYPT_SECRET")
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("ENCRYPT_SECRET must be 16, 24, or 32 bytes")
	}
	return []byte(key), nil
}

func Encrypt(plaintext string, secretKey string) ([]byte, error) {
	key, err := getAESKey(secretKey)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)
	final := append(nonce, ciphertext...)
	return final, nil
}

func Decrypt(data []byte, secretKey string) (string, error) {
	key, err := getAESKey(secretKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("data too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
