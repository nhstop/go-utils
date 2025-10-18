package jwt_utils

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims defines a generic JWT payload
type CustomClaims struct {
	Payload map[string]interface{}
	jwt.RegisteredClaims
}

// JWTManager is an abstract interface to generate and verify tokens
type JWTManager struct {
	SigningMethod jwt.SigningMethod
	SecretKey     interface{} // can be []byte for HMAC or *rsa.PrivateKey for RS256
	PublicKey     interface{} // optional, used for asymmetric verification
	Issuer        string
}

// GenerateJWT generates a JWT token with given payload and expiry
func (j *JWTManager) GenerateJWT(payload map[string]interface{}, expiry time.Duration) (string, error) {
	claims := CustomClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(j.SigningMethod, claims)
	return token.SignedString(j.SecretKey)
}

// VerifyJWT verifies the token using the appropriate key
func (j *JWTManager) VerifyJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Check algorithm matches
		if t.Method.Alg() != j.SigningMethod.Alg() {
			return nil, errors.New("unexpected signing method")
		}

		return j.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// LoadRSAPrivateKeyFromBytes loads an RSA private key from PEM bytes
func LoadRSAPrivateKeyFromBytes(pemBytes []byte) (*rsa.PrivateKey, error) {
	if len(pemBytes) == 0 {
		return nil, errors.New("empty private key bytes")
	}
	return jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
}

// LoadRSAPublicKeyFromBytes loads an RSA public key from PEM bytes
func LoadRSAPublicKeyFromBytes(pemBytes []byte) (*rsa.PublicKey, error) {
	if len(pemBytes) == 0 {
		return nil, errors.New("empty public key bytes")
	}
	return jwt.ParseRSAPublicKeyFromPEM(pemBytes)
}
