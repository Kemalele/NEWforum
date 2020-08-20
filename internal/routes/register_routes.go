package routes

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"time"
	"html/template"
	models "../models"
	services "../services"
)

func HandleRegistration(w http.ResponseWriter, r *http.Request, params url.Values) {
	var user models.User
	var err error
	id, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 Internal server error")
		return
	}

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Email = r.FormValue("email")
	user.Id = id.String()
	user.RegistrationDate = time.Now().String()

	err = services.Register(user)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func GetRegistration(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../internal/templates/registration.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 Internal server error")
		return
	}
	t.Execute(w, nil)
}
