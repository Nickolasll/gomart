package application

import (
	"strconv"
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

// Может нужен отдельный сервис для паролей?
// Я пока решил не плодить сущности
func (jose JOSEService) Hash(password string) string {
	var passwordBytes = []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	return string(hashedPassword)
}

func (jose JOSEService) VerifyPassword(hashedPassword string, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func IsValidNumber(input string) bool {
	number, err := strconv.Atoi(input)
	if err != nil {
		return false
	}
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
