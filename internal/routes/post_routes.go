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
	//categories, err := models.CategoriesByName(params.Get("Name"))
	t.ExecuteTemplate(w, "write", nil)
}

func SavePostHandler(w http.ResponseWriter, r *http.Request, params url.Values) {
	var post models.Post
	var err error

	var postcategories models.PostsCategories

	post.Id = services.GenerateId()
	r.ParseForm()
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
	post.Title = r.FormValue("theme")
	post.User.Id = user.Id

	r.ParseMultipartForm(0)
	arr := []string{}
	if r.FormValue("standard") == "standard" {
		postcategories.Id = services.GenerateId()
		postcategories.Category.Id, err = models.ValidateCategory(r.FormValue("standard"))
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		postcategories.Post.Id = post.Id

		arr = append(arr, "standard")
	}
	if r.FormValue("shadow") == "shadow" {
		postcategories.Id = services.GenerateId()
		postcategories.Category.Id, err = models.ValidateCategory(r.FormValue("shadow"))
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		postcategories.Post.Id = post.Id
		arr = append(arr, "shadow")
	}
	if r.FormValue("thinkertoy") == "thinkertoy" {
		postcategories.Id = services.GenerateId()
		postcategories.Category.Id, err = models.ValidateCategory(r.FormValue("shadow"))
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		postcategories.Post.Id = post.Id
		arr = append(arr, "thinkertoy")
	} else {
		fmt.Println("Wrong name of category")
	}
	fmt.Println(postcategories)
	fmt.Println(arr)
	// postcategories.Category.Id, ok = models.ValidateCategory(r.FormValue("standard"))
	// if !ok {
	// 	fmt.Fprintf(w, err.Error())
	// 	return
	// }
	// postcategories.Category.Id, ok = models.ValidateCategory(r.FormValue("shadow"))
	// if !ok {
	// 	fmt.Fprintf(w, err.Error())
	// 	return
	// }
	// postcategories.Category.Id, ok = models.ValidateCategory(r.FormValue("thinkertoy"))
	// if !ok {
	// 	fmt.Fprintf(w, err.Error())
	// 	return
	// }
	err = services.NewPost(post)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	http.Redirect(w, r, "/", 302)
}
