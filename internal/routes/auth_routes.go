package routes

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"time"
	"html/template"
	services "../services"
)


func GetAuth(w http.ResponseWriter, r *http.Request, params url.Values) {
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
	sessionToken, _ := uuid.NewV4()
	Cache[sessionToken.String()] = username
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.String(),
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

