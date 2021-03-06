package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	models "../internal/models"
	routes "../internal/routes"
	services "../internal/services"
	router "../pkg/router"
)

func main() {
	err := models.Init("forum.db")
	if err != nil {
		log.Fatal(err)
	}
	routes.Cache = *services.NewCache()
	r := router.New(routes.GetMain)
	r.Handle("GET", "/", routes.GetMain)
	r.Handle("GET", "/write", routes.WritePost)
	r.Handle("GET", "/registration", routes.GetRegistration)
	r.Handle("GET", "/authentication", routes.GetAuth)
	r.Handle("GET", "/post/:id", routes.HandlePostPage)
	r.Handle("GET", "/logout", routes.HandleLogout)

	r.Handle("POST", "/post/:id/_method=POST", routes.SaveCommentHandler)
	r.Handle("POST", "/savePost", routes.SavePostHandler)
	r.Handle("POST", "/registration", routes.HandleRegistration)
	r.Handle("POST", "/authentication", routes.HandleAuth)
	r.Handle("POST", "/rate", routes.Rate)

	r.Handle("POST", "/post/:id/_method=DELETE", routes.DeleteCommentHandler)

	r.Handle("GET", "/moderation", routes.HandleModeration)
	r.Handle("POST", "/moderationSave", routes.HandleModerationSave)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("app is running on %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
