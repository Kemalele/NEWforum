package models

import "fmt"

type Comment struct {
	Id string
	Description string
	PostDate string
	User User
	PostId string
}

func CommentsByPostId(postId string) ([]Comment,error) {
	var comments []Comment
	query := fmt.Sprintf("SELECT * FROM comment WHERE PostId LIKE '%s'", postId)
	rows,err := Db.Query(query)
	if err != nil {
		return nil,err
	}

	for rows.Next() {
		comment := Comment{}
		err := rows.Scan(&comment.Id,&comment.Description,&comment.PostDate,&comment.User.Id,&comment.PostId)
		if err != nil {
			return nil,err
		}

		comment.User, err = UserById(comment.User.Id)
		if err != nil {
			return nil, err
		}

		comments = append(comments,comment)
	}

	return comments,nil
}

func AddComment(comment Comment,sql SQLDB) error{
	query := "INSERT INTO COMMENT (Id,Description,Post_date,UserId,PostId) values ($1,$2,$3,$4,$5)"
	_, err := sql.Exec(query,comment.Id,comment.Description,comment.PostDate,comment.User.Id,comment.PostId)

	if err != nil {
		return err
	}
	return nil
}

func DeleteComment(commentId string, sql SQLDB) error {

	result, err := sql.Exec("DELETE FROM comment WHERE Id = $1", commentId)
	if err != nil {
		return err
	}
	return nil
}