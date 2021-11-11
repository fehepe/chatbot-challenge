package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/fehepe/chatbot-challenge/internal/models"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("../../web/*.html"))
}

// loginHandler serves form for users to login with
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "login.html", nil)
}

// loginAuthHandler authenticates user login
func LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("username:", username, "password:", password)
	// retrieve password from db to compare (hash) with user supplied password's hash

	params := map[string]string{
		"username": username,
		"password": password,
	}

	response, err := Post(params, "/api/login")
	if err != nil {
		fmt.Println("error trying to  Hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "Conection Error.")
		return
	}
	var resp models.Response

	if err = json.Unmarshal(response, &resp); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	if resp.Message == "success" {
		//fmt.Fprint(w, "You have successfully logged in :)")
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")
}

// registerHandler serves form for registring new users
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

//registerAuthHandler creates new user in database
func RegisterAuthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	name := r.FormValue("name")
	password := r.FormValue("password")
	if password == "" || name == "" || username == "" || len(password) < 4 || len(name) < 4 || len(username) < 4 {
		tpl.ExecuteTemplate(w, "register.html", "please check username, name and password must be longer than 4")
		return
	}

	params := map[string]string{
		"name":     name,
		"username": username,
		"password": password,
	}

	response, err := Post(params, "/api/register")
	if err != nil {
		fmt.Println("error trying to register the user")
		tpl.ExecuteTemplate(w, "register.html", "Conection Error.")
		return
	}
	var resp models.User

	if err = json.Unmarshal(response, &resp); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	if resp.Id != 0 {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	fmt.Println("error inserting new user")
	tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
}
