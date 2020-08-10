package main

import (
	models "../models"
	"fmt"
	_ "github.com/satori/go.uuid"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"net/http"
	"net/url"
	"time"
)


func getMain(w http.ResponseWriter,r *http.Request,params url.Values) {
	t,err := template.ParseFiles("../templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	username, authed := authenticated(r)
	sortBy := r.FormValue("sortBy")

	response := struct {
		Posts []models.Post
		Authed bool
	}{
		Posts: posts.Body,
		Authed: authed,
	}

	switch sortBy {
		case "created":
			if authed {
				user,err := users.UserByName(username)
				if err != nil {
					fmt.Println(err.Error())
					break
				}

				p,err := models.SortedPosts(sortBy,user)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				response.Posts = p.Body
			}
		default:
	}



	t.Execute(w,response)
}

//func handleMain(w http.ResponseWriter,r *http.Request,params url.Values) {
//	t,err := template.ParseFiles("../templates/index.html")
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	fmt.Println("post")
//	sortBy := r.FormValue("sortBy")
//	userName , status := authenticated(r)
//
//	if status != http.StatusOK{
//		fmt.Fprintf(w,"ERROR: %v",status)
//		return
//	}
//
//	user,err := users.UserByName(userName)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//
//	p,err := models.SortedPosts(sortBy,user)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//
//	response := struct {
//		Posts []models.Post
//		Authed bool
//	}{
//		Posts: p.Body,
//		Authed: authenticated(r),
//	}
//
//	w.WriteHeader(http.StatusOK)
//	t.Execute(w,response)
//}


func handlePostPage(w http.ResponseWriter, r *http.Request,params url.Values) {
	fmt.Fprintf(w,"%v",params.Get("id"))
}

func writePost(w http.ResponseWriter, r *http.Request,params url.Values){
	t,err := template.ParseFiles("../templates/write.html")

	_ , ok := authenticated(r)
	if !ok {
		http.Redirect(w,r,"/authentication",http.StatusUnauthorized)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "write",nil)
}

func savepostHandler(w http.ResponseWriter, r *http.Request,params url.Values){
	var post models.Post
	var err error

	post.Id = GenerateId()
	post.Description = r.FormValue("description")
	t := time.Now()
	post.PostDate = t.Format(time.RFC1123)
	userid , ok := authenticated(r)
	if !ok {
		http.Redirect(w,r,"/authentication",http.StatusUnauthorized)
		return
	}
	post.UserId = userid
	post.Category = "choumi"
	post.Theme = r.FormValue("theme")

	err = NewPost(post)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	http.Redirect(w,r,"/", 302)
}

func handleAuth(w http.ResponseWriter, r *http.Request,params url.Values) {
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

func getAuth(w http.ResponseWriter, r *http.Request,params url.Values){
	t,err := template.ParseFiles("../templates/authentication.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w,"%v",http.StatusInternalServerError)
		return
	}
	t.Execute(w,nil)
}

func handleRegistration(w http.ResponseWriter, r *http.Request,params url.Values) {
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
}

func getRegistration(w http.ResponseWriter, r *http.Request,params url.Values) {
	t,err := template.ParseFiles("../templates/registration.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w,"500 Internal server error")
		return
	}
	t.Execute(w,nil)
}

