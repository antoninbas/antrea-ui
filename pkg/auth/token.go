package auth

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	Issuer = "ui.antrea.io"
)

type Token struct {
	Raw string
	// ExpiresIn time.Duration
}

type TokenManager interface {
	GetToken() (*Token, error)
	VerifyToken(rawToken string) error
}

type tokenManager struct {
	SigningKeyID  string
	SigningKey    *rsa.PrivateKey
	SigningMethod jwt.SigningMethod
}

type JWTAccessClaims struct {
	jwt.StandardClaims
}

func NewTokenManager(keyID string, key *rsa.PrivateKey) *tokenManager {
	return &tokenManager{
		SigningKeyID:  keyID,
		SigningKey:    key,
		SigningMethod: jwt.SigningMethodRS512,
	}
}

func (m *tokenManager) GetToken() (*Token, error) {
	createdAt := time.Now()
	// expiresIn := 1 * time.Hour
	claims := &JWTAccessClaims{
		StandardClaims: jwt.StandardClaims{
			// ExpiresAt: createdAt.Add(expiresIn).Unix(),
			Issuer:   Issuer,
			IssuedAt: createdAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(m.SigningMethod, claims)
	if m.SigningKeyID != "" {
		token.Header["kid"] = m.SigningKeyID
	}

	access, err := token.SignedString(m.SigningKey)
	if err != nil {
		return nil, err
	}
	return &Token{
		Raw: access,
		// ExpiresIn: expiresIn,
	}, nil
}

func (m *tokenManager) VerifyToken(rawToken string) error {
	_, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		return &m.SigningKey.PublicKey, nil
	})
	if err != nil {
		return err
	}
	return nil
}
