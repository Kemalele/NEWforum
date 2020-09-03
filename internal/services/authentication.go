package services

import (
	"errors"
	"net/http"

	models "../models"
)

type Cache struct {
	tokenUser map[string]string
	userToken map[string]string
}

func NewCache() *Cache {
	var c Cache
	c.tokenUser = make(map[string]string)
	c.userToken = make(map[string]string)
	return &c
}

func (c *Cache) UserByToken(token string) string {
	return c.tokenUser[token]
}

func (c *Cache) Add(username, token string) {
	c.tokenUser[token] = username
	c.userToken[username] = token
}

func (c *Cache) DeleteToken(token string) {
	username := c.tokenUser[token]

	delete(c.tokenUser, token)
	delete(c.userToken, username)
}

func (c *Cache) DeleteUser(username string) {
	token := c.userToken[username]

	delete(c.tokenUser, token)
	delete(c.userToken, username)
}

func (c *Cache) TokenExists(token string) bool {
	username := c.tokenUser[token]
	if username == "" {
		return false
	}

	return true
}

func (c *Cache) UserExists(username string) bool {
	token := c.userToken[username]
	if token == "" {
		return false
	}

	return true
}

func (c *Cache) Match(token, user string) bool {
	tokenByUser := c.userToken[user]
	if tokenByUser != token {
		return false
	}

	return true
}

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

func Authenticated(r *http.Request, cache *Cache) (string, bool) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return "", false
	}

	sessionToken := c.Value

	if !cache.TokenExists(sessionToken) {
		cache.DeleteToken(sessionToken)
		return "", false
	} else if !cache.Match(sessionToken, cache.UserByToken(sessionToken)) {
		return "", false
	}

	// nickname
	return cache.UserByToken(sessionToken), true
}
