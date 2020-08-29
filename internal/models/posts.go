package models

import (
	"errors"
	"fmt"
)

type PostDTO struct {
	Post     Post
	Likes    int
	Dislikes int
}

type Post struct {
	Id          string
	Description string
	PostDate    string
	User        User
	Category    Category
	Title       string
}

func AllPosts() ([]PostDTO, error) {
	rows, err := Db.Query("SELECT * FROM Post")
	var postsLikes []PostDTO
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.Id, &post.Description, &post.PostDate, &post.User.Id, &post.Category.Id, &post.Title)
		if err != nil {
			return nil, err
		}
		post.Category, err = CategoryById(post.Category.Id)
		if err != nil {
			return nil, err
		}
		post.User, err = UserById(post.User.Id)
		if err != nil {
			return nil, err
		}

		likes, err := LikedPostCount(post.Id, "like")
		if err != nil {
			return nil, err
		}

		dislikes, err := LikedPostCount(post.Id, "dislike")
		if err != nil {
			return nil, err
		}

		postsLikes = append(postsLikes, PostDTO{Post: post, Likes: likes, Dislikes: dislikes})
	}
	return postsLikes, nil
}

func AddPost(post Post, sql SQLDB) error {
	_, err := sql.Exec("INSERT INTO POST (Id,Description,Post_date,UserId,CategoryId,Title) values ($1,$2,$3,$4,$5,$6)", post.Id, post.Description, post.PostDate, post.User.Id, post.Category.Id, post.Title)
	if err != nil {
		return err
	}
	return nil
}

func PostById(id string) (Post, error) {
	post := Post{}
	query := fmt.Sprintf("SELECT * FROM Post WHERE Id LIKE '%s'", id)
	rows, err := Db.Query(query)
	if err != nil {
		return Post{}, err
	}

	for rows.Next() {
		err := rows.Scan(&post.Id, &post.Description, &post.PostDate, &post.User.Id, &post.Category.Id, &post.Title)
		if err != nil {
			return Post{}, err
		}

		post.User, err = UserById(post.User.Id)
		if err != nil {
			return Post{}, err
		}

		post.Category, err = CategoryById(post.Category.Id)
		if err != nil {
			return Post{}, err
		}
	}

	return post, nil
}

func SortedPosts(sortBy string, user User) ([]PostDTO, error) {
	var query string
	var postsLikes []PostDTO

	if sortBy == "created" {
		query = fmt.Sprintf("SELECT * FROM POST WHERE UserId LIKE '%s';", user.Id)
	} else if sortBy == "liked" {
		query = fmt.Sprintf("SELECT p.Id, p.Description, p.Post_date, p.UserId, p.CategoryId, p.Title FROM Post p LEFT JOIN likedPosts l ON p.Id = l.PostId WHERE l.UserId LIKE '%s' AND l.Value LIKE 'like';", user.Id)
	} else {
		return postsLikes, errors.New("no such parameter to sort")
	}

	rows, err := Db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.Id, &post.Description, &post.PostDate, &post.User.Id, &post.Category.Id, &post.Title)
		if err != nil {
			return nil, err
		}

		post.Category, err = CategoryById(post.Category.Id)
		if err != nil {
			return nil, err
		}

		post.User, err = UserById(post.User.Id)
		if err != nil {
			return nil, err
		}

		likes, err := LikedPostCount(post.Id, "like")
		if err != nil {
			return nil, err
		}

		dislikes, err := LikedPostCount(post.Id, "dislike")
		if err != nil {
			return nil, err
		}

		postsLikes = append(postsLikes, PostDTO{Post: post, Likes: likes, Dislikes: dislikes})
	}

	return postsLikes, nil
}
