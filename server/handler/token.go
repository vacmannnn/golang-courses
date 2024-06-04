package handler

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

var secretKey = []byte("bananchiki")

func createToken(user userInfo, expTime int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Username": user.Username,
			"role":     user.role,
			"exp":      time.Now().Add(time.Minute * time.Duration(expTime)).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("unable to extract claims")
	}

	return claims["role"] == float64(admin), nil
}
