package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"html/template"
	services "../services"
	models "../models"
)

var Cache map[string]string


func GetMain(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../internal/templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username, authed := services.Authenticated(r, Cache)
	user, err := models.UserByName(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	sortBy := r.FormValue("sortBy")

	response := struct {
		Posts  []models.Post
		Authed bool
		User models.User
	}{
		Posts:  nil,
		Authed: authed,
		User: user,
	}

	switch sortBy {
	case "created":
		if authed {
			user, err := models.UserByName(username)
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			posts, err := models.SortedPosts(sortBy, user)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			response.Posts = posts
		}
	default:
		posts, err := models.AllPosts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			break
		}
		response.Posts = posts
	}

	t.Execute(w, response)
}
