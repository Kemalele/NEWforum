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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !services.ValidURL(r.URL.Path, 2) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Page not found")
		return
	}

	post, err := models.PostById(params.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	if post == (models.Post{}) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found")
		return
	}

	username, authed := services.Authenticated(r, &Cache)
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

	if !services.ValidURL(r.URL.Path, 1) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Page not found")
		return
	}

	_, authed := services.Authenticated(r, &Cache)
	if !authed {
		http.Redirect(w, r, "/authentication", http.StatusUnauthorized)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := struct {
		Authed bool
	}{
		Authed: authed,
	}
	t.ExecuteTemplate(w, "write", response)
}

func SavePostHandler(w http.ResponseWriter, r *http.Request, params url.Values) {
	var post models.Post
	var err error
	categories := []string{"standard", "shadow", "thinkertoy"}

	post.Id = services.GenerateId()
	post.Description = r.FormValue("description")
	t := time.Now()
	post.PostDate = t.Format(time.RFC1123)

	username, ok := services.Authenticated(r, &Cache)
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

	post.Title = r.FormValue("theme")
	post.User.Id = user.Id

	for _, name := range categories {
		if r.FormValue(name) == name {
			var postcategories models.PostsCategories
			postcategories.Id = services.GenerateId()
			postcategories.Category.Id, err = models.ValidateCategory(name)
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			postcategories.Post.Id = post.Id
			err := models.AddCategoryToPost(postcategories, models.Db)
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	err = services.NewPost(post)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	http.Redirect(w, r, "/", 302)
}
