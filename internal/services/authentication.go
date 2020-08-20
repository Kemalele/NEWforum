package services

import (
	models "../models"
	"errors"
	"net/http"
)

func CorrectUser(username, password string) error {
	user, err := models.UserByName(username)
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
		} else {
			return errors.New("wrong password")

		}

	} else {
		return errors.New("wrong username")
	}

}

func Authenticated(r *http.Request, cache map[string]string) (string, bool) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return "", false
	}

	sessionToken := c.Value

	// nickname
	response := cache[sessionToken]
	if response == "" {
		return "", false
	}

	return response, true
}
