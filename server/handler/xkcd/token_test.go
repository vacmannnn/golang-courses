package xkcd

import (
	"testing"
)

func TestCreateToken(t *testing.T) {
	cases := []struct {
		name    string
		user    userInfo
		expTime int
	}{
		{
			name:    "Valid data",
			user:    userInfo{Username: "Vlad", Password: "qwerty123"},
			expTime: 10,
		},

		{
			name:    "Name contains smiles",
			user:    userInfo{Username: "😁", Password: "321ytrewq"},
			expTime: 10,
		},

		{
			name:    "Empty user data",
			user:    userInfo{Username: "", Password: ""},
			expTime: 0,
		},

		{
			name:    "Negative expTime",
			user:    userInfo{Username: "Vlad", Password: "qwerty123"},
			expTime: -123,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := createToken(tc.user, tc.expTime)
			if err != nil {
				t.Errorf("got error: %v", err)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	cases := []struct {
		name    string
		token   string
		isValid bool
	}{
		{
			name: "Valid data with secret key 'bananchiki'",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJVc2VybmFtZSI6IlZsYWQiLCJleHAiOjE3MTc0OTIzOTEsInJvbGUiOjB9." +
				"8b_J5iP09TJhhDeZKTkNOpdHTDYoJ6thlWzB3mWYkNw",
			isValid: false,
		},

		{
			name:    "Name contains smiles",
			token:   "it shouldn't be correct",
			isValid: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := verifyToken(tc.token)
			if valid != tc.isValid {
				t.Errorf("valid expectation is to be %t, but got %t", tc.isValid, valid)
			}
			if valid && err != nil {
				t.Errorf("got error: %v", err)
			}
		})
	}
}

func TestCreateAndVerify(t *testing.T) {
	usr := userInfo{Username: "Vlad", Password: "qwerty123"}
	token, err := createToken(usr, 10)
	if err != nil {
		t.Errorf("creating token: %v", err)
	}

	valid, err := verifyToken(token)
	if !valid || err != nil {
		t.Errorf("token should be valid with no errors, valid: %t, err: %v", valid, err)
	}
}
