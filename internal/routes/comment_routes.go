package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	models "../models"
	services "../services"
	uuid "github.com/satori/go.uuid"
)

func SaveCommentHandler(w http.ResponseWriter, r *http.Request, params url.Values) {
	postId := params.Get("id")
	username, _ := services.Authenticated(r, Cache)

	user, err := models.UserByName(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	post, err := models.PostById(postId)
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

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request, params url.Values) {
	postId := params.Get("id")
	commentId := params.Get("commentId")

	err := models.DeleteComment(commentId, models.Db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

	http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)

}
