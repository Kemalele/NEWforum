package services

import (
	"errors"
	"fmt"

	models "../models"
)

func NewComment(comment models.Comment) error {
	err := validComment(comment)
	if err != nil {
		return err
	}

	err = models.AddComment(comment, models.Db)
	if err != nil {
		return err
	}

	return nil
}

func validComment(c models.Comment) error {
	if len(c.Description) < 3 {
		fmt.Println(c.Description)
		return errors.New("title must be at least 3 symbol")
	}
	return nil
}
