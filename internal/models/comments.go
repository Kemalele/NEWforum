package models

import "fmt"

type CommentDTO struct {
	Comment  Comment
	Likes    int
	Dislikes int
}

type Comment struct {
	Id          string
	Description string
	PostDate    string
	User        User
	Post        Post
}

func CommentsByPostId(postId string) ([]CommentDTO, error) {
	var comments []CommentDTO
	query := fmt.Sprintf("SELECT * FROM comment WHERE PostId LIKE '%s'", postId)
	rows, err := Db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		comment := Comment{}
		err := rows.Scan(&comment.Id, &comment.Description, &comment.PostDate, &comment.User.Id, &comment.Post.Id)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		comment.User, err = UserById(comment.User.Id)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		comment.Post, err = PostById(comment.Post.Id)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		likes, err := LikedCommentCount(comment.Id, "like")
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		dislikes, err := LikedCommentCount(comment.Id, "dislike")
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		comments = append(comments, CommentDTO{Comment: comment, Likes: likes, Dislikes: dislikes})
	}

	return comments, nil
}

func CommentById(commentId string) (Comment, error) {
	comment := Comment{}
	query := fmt.Sprintf("SELECT * FROM comment WHERE Id LIKE '%s'", commentId)
	rows, err := Db.Query(query)
	if err != nil {
		return Comment{}, err
	}

	for rows.Next() {
		err := rows.Scan(&comment.Id, &comment.Description, &comment.PostDate, &comment.User.Id, &comment.Post.Id)
		if err != nil {
			return Comment{}, err
		}

		comment.User, err = UserById(comment.User.Id)
		if err != nil {
			return Comment{}, err
		}

		comment.Post, err = PostById(comment.Post.Id)
		if err != nil {
			return Comment{}, err
		}

	}

	return comment, nil
}

func AddComment(comment Comment, sql SQLDB) error {
	query := "INSERT INTO COMMENT (Id,Description,Post_date,UserId,PostId) values ($1,$2,$3,$4,$5)"
	_, err := sql.Exec(query, comment.Id, comment.Description, comment.PostDate, comment.User.Id, comment.Post.Id)

	if err != nil {
		return err
	}
	return nil
}

func DeleteComment(commentId string, sql SQLDB) error {
	_, err := sql.Exec("DELETE FROM comment WHERE Id = $1", commentId)
	if err != nil {
		return err
	}
	return nil
}
