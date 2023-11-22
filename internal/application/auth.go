package application

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JOSEService struct {
	TokenExp  int
	SecretKey string
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func (jose JOSEService) IssueToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jose.TokenExp) * time.Second)),
		},
		UserID: userID.String(),
	})

	tokenString, err := token.SignedString([]byte(jose.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (jose JOSEService) ParseUserID(tokenString string) *uuid.UUID {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jose.SecretKey), nil
		},
	)
	if err != nil {
		return nil
	}

	if !token.Valid {
		return nil
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil
	}
	return &userID
}

func (jose JOSEService) Hash(password string) string {
	var passwordBytes = []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	return string(hashedPassword)
}

func (jose JOSEService) VerifyPassword(hashedPassword string, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))
	return err == nil
}
