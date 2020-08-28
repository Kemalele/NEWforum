package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	models "../models"
	services "../services"
)

func HandlePostPage(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../internal/templates/post.html")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post, err := models.PostById(params.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	username, authed := services.Authenticated(r, Cache)
	user, err := models.UserByName(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	comments, err := models.CommentsByPostId(post.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	response := struct {
		Post     models.Post
		Authed   bool
		Comments []models.CommentDTO
		User     models.User
	}{
		post,
		authed,
		comments,
		user,
	}
	t.Execute(w, response)
}

func WritePost(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../internal/templates/write.html")

	_, ok := services.Authenticated(r, Cache)
	if !ok {
		http.Redirect(w, r, "/authentication", http.StatusUnauthorized)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "write", nil)
}

func SavePostHandler(w http.ResponseWriter, r *http.Request, params url.Values) {
	var post models.Post
	var err error

	post.Id = services.GenerateId()
	post.Description = r.FormValue("description")
	t := time.Now()
	post.PostDate = t.Format(time.RFC1123)
	username, ok := services.Authenticated(r, Cache)
	if !ok {
		http.Redirect(w, r, "/authentication", http.StatusUnauthorized)
		return
	}

	user, err := models.UserByName(username)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post.User.Id = user.Id
	fmt.Println(user.Id)
	post.Category.Id, err = models.ValidateCategory(r.FormValue("category"))
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	post.Title = r.FormValue("theme")

	err = services.NewPost(post)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	http.Redirect(w, r, "/", 302)
}
