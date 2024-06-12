package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type LoginError struct {
	ErrorExists  bool
	ErrorMessage string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var data LoginError

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println("in parsing form:", err)
		}
		login := r.FormValue("login")
		pass := r.FormValue("password")

		if login == "" || pass == "" {
			data.ErrorMessage = "Введите логин и пароль."
			goto executeTemplate
		}

		var jsonStr = []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, login, pass))
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/login", bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Printf("creating request: %v\n", err)
			data.ErrorMessage = "Попробуйте ввести пароль еще раз."
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("sending request: %v\n", err)
			data.ErrorMessage = "Попробуйте ввести пароль еще раз."
			goto executeTemplate
		}
		str, _ := io.ReadAll(res.Body)
		if res.StatusCode != http.StatusOK {
			data.ErrorMessage = "Неправильный логин или пароль."
			goto executeTemplate
		}

		cookie := http.Cookie{
			Name:  "jwtTokenCookie",
			Value: string(str),
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/comics", 301)
	}

executeTemplate:
	tmpl, _ := template.ParseFiles("templates/login.html")
	data.ErrorExists = data.ErrorMessage != ""
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println("executing login template:", err)
	}

}
