package models

import (
	"fmt"
)

type User struct {
	Id               string
	Username         string
	Password         string
	Email            string
	RegistrationDate string
}

func AllUsers() ([]User, error) {
	rows, err := Db.Query("SELECT * FROM User")
	var users []User

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		usr := User{}
		err := rows.Scan(&usr.Id, &usr.Username, &usr.Password, &usr.Email, &usr.RegistrationDate)
		if err != nil {
			return nil, err
		}
		users = append(users, usr)
	}
	return users, nil
}

func AddUser(usr User, sql SQLDB) error {
	_, err := sql.Exec("INSERT INTO USER (Id,Username,Password,Email,RegistrationDate) values ($1,$2,$3,$4,$5)", usr.Id, usr.Username, usr.Password, usr.Email, usr.RegistrationDate)
	if err != nil {
		return err
	}
	return nil
}

func UserByName(username string) (User, error) {
	usr := User{}
	query := fmt.Sprintf("SELECT * FROM User WHERE username LIKE '%s'", username)
	rows, err := Db.Query(query)
	if err != nil {
		return User{}, err
	}

	for rows.Next() {
		err := rows.Scan(&usr.Id, &usr.Username, &usr.Password, &usr.Email, &usr.RegistrationDate)
		if err != nil {
			return User{}, err
		}
	}

	return usr, nil
}

func UserById(userId string) (User, error) {
	usr := User{}
	query := fmt.Sprintf("SELECT * FROM User WHERE Id LIKE '%s'", userId)
	rows, err := Db.Query(query)
	if err != nil {
		return User{}, err
	}

	for rows.Next() {
		err := rows.Scan(&usr.Id, &usr.Username, &usr.Password, &usr.Email, &usr.RegistrationDate)
		if err != nil {
			return User{}, err
		}
	}

	return usr, nil
}
