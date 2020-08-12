package models

import (
	"errors"
	"fmt"
)
type Post struct {
	Id          string
	Description string
	PostDate    string
	UserId      string
	Category 	string
	Theme 		string
}

func AllPosts() ([]Post,error) {
	rows,err := Db.Query("SELECT * FROM Post")
	var posts []Post
	if err != nil {
		return nil,err
	}
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.Id,&post.Description,&post.PostDate,&post.UserId,&post.Category,&post.Theme)
		if err != nil {
			return nil, err
		}
		posts = append(posts,post)
	}
	return posts,nil
}

func AddPost(post Post,sql SQLDB) error {
	_,err := sql.Exec("INSERT INTO POST (Id,Description,Post_date,UserId,Category,Theme) values ($1,$2,$3,$4,$5,$6)",post.Id,post.Description,post.PostDate,post.UserId,post.Category,post.Theme)
	if err != nil { 
		return err
	}
	return nil
}

//func GetPostById(Id string)(Post,error) {
//	for _,post := range p.Body {
//		if post.Id == Id {
//			return post,nil
//		}
//	}
//	return Post{},errors.New("no such user")
//}

func SortedPosts(sortBy string,user User) ([]Post,error) {
	var query string
	var posts []Post

	if sortBy == "created" {
		query = fmt.Sprintf("SELECT * FROM POST ORDER BY CASE userid WHEN '%s' THEN 1 ELSE 2 END;", user.Username)
	}else {
		return posts,errors.New("no such parameter to sort")
	}

	rows,err := Db.Query(query)
	if err != nil {
		return nil,err
	}

	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.Id,&post.Description,&post.PostDate,&post.UserId,&post.Category,&post.Theme)
		if err != nil {
			return nil,err
		}
		posts = append(posts,post)
	}

	return posts,nil
}