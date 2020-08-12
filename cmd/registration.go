package main

import (
	models "../models"
	"errors"
)

func register(usr models.User)  error {

	err := validUser(usr)
	if err != nil {
		return err
	}

	err = uniqueUser(usr)
	if err != nil {
		return err
	}

	usr.Password,err = encrypt(usr.Password)
	if err != nil {
		return err
	}

	models.AddUser(usr,models.Db)
	return nil
}

func validUser(usr models.User) error{
	if len(usr.Password) < 6 {
		return errors.New("password must be at least 6 symbols")
	}

	if len(usr.Username) < 3 {
		return errors.New("username must be at least 3 symbols")
	}

	return nil
}

func uniqueUser(newUsr models.User) error{
	user,err := models.UserByName(newUsr.Username)
	if err != nil {
		return err
	}

	if user != (models.User{}) {
		return errors.New("username already taken")
	}

	return nil
}