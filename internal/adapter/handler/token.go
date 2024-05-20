package handler

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

var secretKey = []byte("bananchiki") // TODO: config.yaml

var testUsers = []User{{role: admin, Username: "admin", Password: "admin"},
	{role: user, Username: "user", Password: "user"}} // TODO: config.yaml

const (
	admin = iota
	user
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	role     int
}

func createToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"role":     user.role,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
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

	fmt.Println(token.Claims)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("unable to extract claims")
	}

	return claims["role"] == admin, nil
}

func auth(user User) (int, error) {
	for _, u := range testUsers {
		if u.Username == user.Username && u.Password == user.Password {
			return u.role, nil
		}
	}
	return -1, fmt.Errorf("invalid user")
}
