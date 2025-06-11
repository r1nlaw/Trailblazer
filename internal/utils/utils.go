package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"trailblazer/internal/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
)

func LocationFromPoint(p string) models.Location {
	p = strings.TrimPrefix(p, "POINT(")
	p = strings.TrimSuffix(p, ")")
	numbers := strings.Fields(p)
	lon, _ := strconv.ParseFloat(numbers[0], 64)
	lat, _ := strconv.ParseFloat(numbers[1], 64)
	loc := models.Location{
		Lat: lat,
		Lng: lon,
	}
	return loc
}

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) bool
}

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (b *BcryptHasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %v", err)
	}
	return string(hashedPassword), nil
}

func (b *BcryptHasher) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

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
	expirationTime := time.Now().Add(24 * time.Hour)

	payload := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
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
