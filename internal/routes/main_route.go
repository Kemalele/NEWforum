package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	models "../models"
	services "../services"
	uuid "github.com/satori/go.uuid"
)

var Cache map[string]string

func GetMain(w http.ResponseWriter, r *http.Request, params url.Values) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Page not found")
		return
	}

	t, err := template.ParseFiles("../internal/templates/index.html")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username, authed := services.Authenticated(r, Cache)

	user, err := models.UserByName(username)
	if err != nil {
		fmt.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
	}

	sortBy := r.FormValue("sortBy")
	fmt.Println(sortBy)
	response := struct {
		Posts  []models.PostDTO
		Authed bool
		User   models.User
	}{
		Posts:  nil,
		Authed: authed,
		User:   user,
	}

	switch sortBy {
	case "created", "liked", "standard", "shadow", "thinkertoy":
		if authed {
			user, err := models.UserByName(username)
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				break
			}

			posts, err := models.SortedPosts(sortBy, user)
			if err != nil {
				fmt.Println("main_route : 58")
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			response.Posts = posts
		}

	default:
		posts, err := models.AllPosts()
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.Posts = posts
	}
	err = t.Execute(w, response)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Rate(w http.ResponseWriter, r *http.Request, params url.Values) {
	requestBody := struct {
		Action   string `json:"action"`
		Target   string `json:"target"`
		TargetID string `json:"targetId"`
		UserID   string `json:"userId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if requestBody.Action != "like" && requestBody.Action != "dislike" {

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newId, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	user, err := models.UserById(requestBody.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	switch requestBody.Target {
	case "post":
		models.DeleteLikedPost(requestBody.UserID, requestBody.TargetID, models.Db)

		post, err := models.PostById(requestBody.TargetID)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rate := models.LikedPost{
			Id:    newId.String(),
			Value: requestBody.Action,
			Post:  post,
			User:  user,
		}

		err = models.AddLikedPosts(rate, models.Db)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "comment":
		models.DeleteLikedComment(requestBody.UserID, requestBody.TargetID, models.Db)

		comment, err := models.CommentById(requestBody.TargetID)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rate := models.LikedComment{
			Id:      newId.String(),
			Value:   requestBody.Action,
			Comment: comment,
			User:    user,
		}

		err = models.AddLikedComments(rate, models.Db)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
