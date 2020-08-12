package main

import (
	models "../models"
	"errors"
	"net/http"
)

func correctUser(username,password string) error{
	user,err := models.UserByName(username)
	if err != nil {
		return err
	}

	if user.Username == username {
		decryptedPass, err := decrypt(user.Password)
		if err != nil {
			return err
		}

		if decryptedPass == password {
			return nil
		}else {
			return errors.New("wrong password")

		}

	}else {
		return errors.New("wrong username")
	}

}

func authenticated(r *http.Request) (string,bool){
	c,err := r.Cookie("session_token")
	if err != nil {
		return "", false
	}

	sessionToken := c.Value

	// nickname
	response := cache[sessionToken]
	if response == "" {
		return "", false
	}

	return response,true
}
