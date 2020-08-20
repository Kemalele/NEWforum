package routes

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"time"
	services "../services"
	models "../models"
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
		PostId:      postId,
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
