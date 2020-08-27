package services

import (
	"errors"
	"fmt"

	models "../models"
)

func NewPost(post models.Post) error {
	err := validPost(post)
	if err != nil {
		return err
	}

	err = models.AddPost(post, models.Db)
	if err != nil {
		return err
	}

	return nil
}

func validPost(p models.Post) error {
	if len(p.Title) < 1 {
		fmt.Println(p.Title)
		return errors.New("title must be at least 1 symbol")
	}

	if len(p.Description) < 1 {
		return errors.New("content must be at least 1 symbol")
	}

	return nil
}
