package main

import (
	models "../models"
	"fmt"
	"log"
	"net/http"
	router "../pkg/router"
	"os"
)

var cache map[string]string

func main() {
	err := models.Init("forum.db")
	if err != nil {
		log.Fatal(err)
	}
	cache = make(map[string]string)

	//css/
	//css := http.FileServer(http.Dir("css"))
	//http.Handle("/css/", http.StripPrefix("/css/", css))

	r := router.New(getMain)
	r.Handle("GET","/",getMain)
	r.Handle("GET","/write",writePost)
	r.Handle("GET","/registration",getRegistration)
	r.Handle("GET","/authentication",getAuth)
	r.Handle("GET","/post/:id",handlePostPage)

	//r.Handle("POST","/",handleMain)
	r.Handle("POST","/savePost",savepostHandler)
	r.Handle("POST","/registration",handleRegistration)
	r.Handle("POST","/authentication",handleAuth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3030"
	}

	fmt.Printf("app is running on %s\n",port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
