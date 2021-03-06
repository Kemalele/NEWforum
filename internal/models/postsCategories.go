package models

import "fmt"

type PostsCategories struct {
	Id       string
	Category Category
	Post     Post
}

func CategoriesByPostId(postId string) ([]Category, error) {
	var categories []Category
	query := fmt.Sprintf("SELECT CategoryId FROM PostsCategories WHERE PostId LIKE '%s'", postId)
	rows, err := Db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category Category

		err := rows.Scan(&category.Id)
		if err != nil {
			return nil, err
		}

		category, err = CategoryById(category.Id)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func AddCategoryToPost(postCategory PostsCategories, sql SQLDB) error {
	query := "INSERT INTO PostsCategories (Id, PostId, CategoryId) values ($1,$2,$3)"
	_, err := sql.Exec(query, postCategory.Id, postCategory.Post.Id, postCategory.Category.Id)
	if err != nil {
		fmt.Println("postsCategories : 42")
		return err
	}
	return nil
}
