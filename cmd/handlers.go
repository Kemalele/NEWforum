package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	models "../models"
	_ "github.com/satori/go.uuid"
	uuid "github.com/satori/go.uuid"
)

func getMain(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username, authed := authenticated(r)
	sortBy := r.FormValue("sortBy")

	response := struct {
		Posts  []models.Post
		Authed bool
	}{
		Posts:  nil,
		Authed: authed,
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
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			break
		}
		response.Posts = posts
	}

	t.Execute(w, response)
}

func handlePostPage(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../templates/post.html")
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

	username, authed := authenticated(r)
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
		Comments []models.Comment
		User     models.User
	}{
		post,
		authed,
		comments,
		user,
	}
	t.Execute(w, response)
}

func saveCommentHandler(w http.ResponseWriter, r *http.Request, params url.Values) {
	postId := params.Get("id")
	post, err := models.PostById(postId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	username, _ := authenticated(r)

	user, err := models.UserByName(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	newId, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	comment := models.Comment{
		Id:          newId.String(),
		Description: r.FormValue("text"),
		PostDate:    time.Now().String(),
		User:        user,
		Post:        post,
	}

	models.AddComment(comment, models.Db)
	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
	return
}

func deleteCommentHandler(w http.ResponseWriter, r *http.Request, params url.Values) {
	postId := params.Get("id")
	commentId := params.Get("commentId")

	err := models.DeleteComment(commentId, models.Db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)

}

func writePost(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../templates/write.html")

	_, ok := authenticated(r)
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

func handleModeration(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../templates/moderation.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "moderation", nil)
}

func handleModerationSave(w http.ResponseWriter, r *http.Request, params url.Values) {
	var category models.Category
	var err error

	category.Id = GenerateId()
	category.Name = r.FormValue("category")

	err = models.AddCategory(category, models.Db)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	http.Redirect(w, r, "/", 302)
}

func savepostHandler(w http.ResponseWriter, r *http.Request, params url.Values) {
	var post models.Post
	var err error

	post.Id = GenerateId()
	post.Description = r.FormValue("description")
	t := time.Now()
	post.PostDate = t.Format(time.RFC1123)
	userid, ok := authenticated(r)
	if !ok {
		http.Redirect(w, r, "/authentication", http.StatusUnauthorized)
		return
	}
	post.User.Id = userid

	post.Category.Id = models.ValidateCategory(r.FormValue("category"))

	post.Title = r.FormValue("theme")

	err = NewPost(post)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	http.Redirect(w, r, "/", 302)
}

func handleAuth(w http.ResponseWriter, r *http.Request, params url.Values) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	err := correctUser(username, password)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	sessionToken, _ := uuid.NewV4()
	cache[sessionToken.String()] = username
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.String(),
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func getAuth(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../templates/authentication.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func handleRegistration(w http.ResponseWriter, r *http.Request, params url.Values) {
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

	err = register(user)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func getRegistration(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../templates/registration.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 Internal server error")
		return
	}
	t.Execute(w, nil)
}
