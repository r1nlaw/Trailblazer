package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Maker interface {
	CreateToken(userID int64) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	UserID    int64     `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("invalid key size: must be at least 32 bytes")
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

func (j *JWTMaker) CreateToken(userID int64) (string, error) {
	payload := jwt.MapClaims{
		"user_id":    userID,
		"expired_at": time.Now().Add(time.Hour * 24).Unix(),
		"created_at": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTMaker) VerifyToken(tokenStr string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	}

	token, err := jwt.Parse(tokenStr, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return convertClaimsToPayload(claims)
}

func convertClaimsToPayload(claims jwt.MapClaims) (*Payload, error) {
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id claim")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid exp claim")
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid iat claim")
	}

	return &Payload{
		UserID:    int64(userID),
		ExpiredAt: time.Unix(int64(exp), 0),
		IssuedAt:  time.Unix(int64(iat), 0),
	}, nil
}
