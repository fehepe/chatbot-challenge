package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/fehepe/chatbot-challenge/internal/models"
	"github.com/fehepe/chatbot-challenge/pkg/hub"
	"github.com/fehepe/chatbot-challenge/pkg/stock"

	"github.com/gorilla/sessions"
)

var (
	key     = []byte("secret")
	tpl     *template.Template
	store   = sessions.NewCookieStore(key)
	hubConn *hub.Hub
)

func init() {
	tpl = template.Must(template.ParseGlob("../../web/*.html"))
	hubConn = hub.NewHub()
	go hubConn.Run()
}

// loginHandler serves form for users to login with
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****LogOutHandler running*****")
	response, err := stock.Post(nil, "/api/logout")
	if err != nil {
		fmt.Println("error doing the login request")
		tpl.ExecuteTemplate(w, "login.html", "Conection Error.")
		return
	}

	var resp models.Response
	if err = json.Unmarshal(response, &resp); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	session, _ := store.Get(r, "session")
	delete(session.Values, "id")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
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

	response, err := stock.Post(params, "/api/login")
	if err != nil {
		fmt.Println("error doing the login request")
		tpl.ExecuteTemplate(w, "login.html", "Conection Error.")
		return
	}

	if strings.Contains(string(response), "message") {
		var resp models.Response
		if err = json.Unmarshal(response, &resp); err != nil {
			log.Fatal("ooopsss! an error occurred, please try again")
		}
		fmt.Println("incorrect password")
		tpl.ExecuteTemplate(w, "login.html", resp.Message)
		return
	}

	var resp models.User
	if err = json.Unmarshal(response, &resp); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	session, _ := store.Get(r, "session")
	session.Values["id"] = resp.Id
	session.Values["name"] = resp.Name
	sessions.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)

}

// registerHandler serves form for registring new users
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****IndexHandler running*****")
	session, _ := store.Get(r, "session")
	_, ok := session.Values["id"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	name, _ := session.Values["name"]

	tpl.ExecuteTemplate(w, "index.html", name)
}

// registerHandler serves form for registring new users
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

func ChatServer(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["id"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	hub.ServeWs(hubConn, w, r)
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

	response, err := stock.Post(params, "/api/register")
	if err != nil {
		fmt.Println("error trying to register the user")
		tpl.ExecuteTemplate(w, "register.html", "Conection Error.")
		return
	}

	var resp models.User
	if err = json.Unmarshal(response, &resp); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}

	if resp.Id == 0 {
		fmt.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}

	session, _ := store.Get(r, "session")
	session.Values["id"] = resp.Id
	session.Values["name"] = resp.Name
	sessions.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
