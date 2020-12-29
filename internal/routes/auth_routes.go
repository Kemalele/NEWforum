package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	services "../services"
	uuid "github.com/satori/go.uuid"
)

func GetAuth(w http.ResponseWriter, r *http.Request, params url.Values) {
	if r.URL.Path != "/authentication" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Page not found")
		return
	}
	t, err := template.ParseFiles("../internal/templates/authentication.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func HandleAuth(w http.ResponseWriter, r *http.Request, params url.Values) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	err := services.CorrectUser(username, password)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	if Cache.UserExists(username) {
		Cache.DeleteUser(username)
	}

	sessionToken, _ := uuid.NewV4()
	Cache.Add(username, sessionToken.String())
	fmt.Println("AUTH - ", Cache)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.String(),
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
