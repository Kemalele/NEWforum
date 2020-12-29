package models

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SQLDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var Db *sql.DB

func Init(dataSource string) error {
	var err error

	// if db is not exists create db
	if _, err = os.Stat(dataSource); os.IsNotExist(err) {
		if err = CreateDB(dataSource); err != nil {
			return err
		}
	}

	Db, err = sql.Open("sqlite3", dataSource)
	if err != nil {
		return err
	}

	return nil
}

func CreateDB(dataSource string) error {
	os.Create(dataSource)

	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE likedPosts (Id TEXT PRIMARY KEY NOT NULL, Value TEXT NOT NULL, PostId REFERENCES Post(Id), UserId REFERENCES User(Id))")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE likedComments (Id TEXT PRIMARY KEY NOT NULL,Value TEXT NOT NULL,CommentId REFERENCES Comment(Id),UserId REFERENCES User(Id))")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE User(Id TEXT PRIMARY KEY NOT NULL,Username TEXT NOT NULL,Password TEXT NOT NULL,Email TEXT NOT NULL,RegistrationDate TEXT NOT NULL)")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE PostsCategories (Id TEXT PRIMARY KEY NOT NULL,PostId REFERENCES Post(Id),CategoryId REFERENCES Category(Id))")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE Post (Id TEXT PRIMARY KEY NOT NULL, Description TEXT NOT NULL, Post_date TEXT NOT NULL, UserId , Title TEXT NOT NULL)")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE Comment(Id TEXT PRIMARY KEY NOT NULL,Description TEXT NOT NULL,Post_date TEXT NOT NULL,UserId REFERENCES User(Id),PostId REFERENCES Post(Id))")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE Category(Id TEXT PRIMARY KEY NOT NULL,Name TEXT NOT NULL)")
	if err != nil {
		return err
	}

	categories := []Category{
		{Id: "1", Name: "standard"},
		{Id: "2", Name: "shadow"},
		{Id: "3", Name: "thinkertoy"},
	}

	for _, category := range categories {
		AddCategory(category, db)
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}
