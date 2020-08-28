package models

import (
	"fmt"
)

type Category struct {
	Id   string
	Name string
}

func AddCategory(category Category, sql SQLDB) error {
	query := "INSERT INTO CATEGORY (Id,Name) values ($1,$2)"
	_, err := sql.Exec(query, category.Id, category.Name)
	fmt.Println("1")
	rows, err := Db.Query("SELECT * FROM Category")
	if err != nil {
		return err
	}
	var id string
	var name string

	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		fmt.Println(id)
		fmt.Println(name)
	}

	rows.Close()
	if err != nil {
		return err
	}
	return nil
}

func CategoryById(Id string) (Category, error) {
	category := Category{}
	query := fmt.Sprintf("SELECT * FROM Category WHERE Id LIKE '%s'", Id)
	rows, err := Db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return Category{}, err
	}

	for rows.Next() {
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			fmt.Println(err.Error())
			return Category{}, err
		}
	}

	return category, nil
}

func ValidateCategory(name string) (string, error) {
	category := Category{}
	query := fmt.Sprintf("SELECT * FROM Category WHERE Name LIKE '%s'", name)
	rows, err := Db.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("-----------------------------")
	if name != category.Name {
		return "", err
	}
	return category.Id, nil
}

func AllCategories() error {
	rows, err := Db.Query("SELECT * FROM Category")
	if err != nil {
		return err
	}
	for rows.Next() {
		category := Category{}
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return err
		}
		fmt.Println(category)

	}
	return nil
}
