package handler

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"time"
)

var secretKey = []byte("bananchiki")

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

const (
	admin = iota
	user
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	role     int
}

var testUsers = []User{{role: admin, Username: "admin", Password: "admin"},
	{role: user, Username: "user", Password: "user"}}

func auth(user User) (int, error) {
	for _, u := range testUsers {
		if u.Username == user.Username && u.Password == user.Password {
			return u.role, nil
		}
	}
	return -1, fmt.Errorf("invalid user")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u User
	json.NewDecoder(r.Body).Decode(&u)
	fmt.Printf("The user request value %v", u)

	role, err := auth(u)
	if err == nil {
		u.role = role
		tokenString, err := createToken(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Errorf("No username found")
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, tokenString)
		log.Println(tokenString)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid credentials")
	}
}

func ProtectedHandler(next func(http.ResponseWriter, *http.Request), checkForAdmin bool, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}

	isAdmin, err := verifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		fmt.Println(err)
		return
	}

	if checkForAdmin && !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Only admin can update comics set")
		return
	}

	fmt.Fprint(w, "Welcome to the the protected area")
	next(w, r)
}

func protectedGet(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ProtectedHandler(next, false, w, r)
	}
}

func protectedUpdate(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ProtectedHandler(next, true, w, r)
	}
}
