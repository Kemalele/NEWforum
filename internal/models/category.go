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
		return Category{}, err
	}

	for rows.Next() {
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return Category{}, err
		}
	}
	return category, nil
}

func ValidateCategory(name string) string {
	category := Category{}
	query := fmt.Sprintf("SELECT * FROM Category WHERE Id LIKE '%s'", name)
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

	fmt.Println("-----------------------------", category)
	if name == category.Name {
		return category.Id
	} else {
		return "wrong name of category"
	}

}
