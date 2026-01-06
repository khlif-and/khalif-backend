package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func GenerateAccessToken(userID uint, username, role, secret string, expiryHours int) (string, error) {
	claims := JWTClaim{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "khalif-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID uint, secret string, expiryDays int) (string, string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["type"] = "refresh"
	claims["exp"] = time.Now().Add(time.Hour * 24 * time.Duration(expiryDays)).Unix()
	claims["jti"] = GenerateRandomString(32) // Add unique entropy

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	hasher := sha256.New()
	hasher.Write([]byte(signedToken))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	return signedToken, hashedToken, nil
}

func GenerateToken(userID uint, username, role, secret string, expHours int) (string, error) {
	return GenerateAccessToken(userID, username, role, secret, expHours)
}

func ValidateToken(tokenString string, secret string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
