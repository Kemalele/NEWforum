package models

import (
	"database/sql"
	"fmt"
)

type LikedPost struct {
	Id    string
	Value string
	Post  Post
	User  User
}

type LikedComment struct {
	Id      string
	Value   string
	Comment Comment
	User    User
}

func LikedPostsByPostId(postId string) ([]LikedPost, error) {
	var likes []LikedPost
	query := fmt.Sprintf("SELECT * FROM likedPosts WHERE PostId LIKE '%s'", postId)
	rows, err := Db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		liked := LikedPost{}
		err := rows.Scan(&liked.Id, &liked.Value, &liked.Post.Id, &liked.User.Id)
		if err != nil {
			return nil, err
		}

		liked.User, err = UserById(liked.User.Id)
		if err != nil {
			return nil, err
		}

		liked.Post, err = PostById(liked.Post.Id)
		if err != nil {
			return nil, err
		}

		likes = append(likes, liked)
	}
	return likes, nil
}

func LikedPostCount(postId, action string) (int, error) {
	var likes int
	query := fmt.Sprintf("SELECT COUNT(*) FROM likedPosts WHERE PostId LIKE '%s' AND Value LIKE '%s'", postId, action)
	err := Db.QueryRow(query).Scan(&likes)
	if err != nil {
		return 0, err
	}

	return likes, nil

}

func AddLikedPosts(liked LikedPost, sql SQLDB) error {
	_, err := sql.Exec("INSERT INTO likedPosts (Id,Value,PostId,UserId) values ($1,$2,$3,$4)", liked.Id, liked.Value, liked.Post.Id, liked.User.Id)
	if err != nil {
		return err
	}

	return nil
}
func PostRate(userId, postId string) (LikedPost, error) {
	var rate LikedPost
	query := fmt.Sprintf("SELECT * FROM likedPosts WHERE UserId LIKE '%s' AND PostId LIKE '%s'", userId, postId)
	err := Db.QueryRow(query).Scan(&rate.Id, &rate.Value, &rate.Post.Id, &rate.User.Id)
	if err != nil {
		if err != sql.ErrNoRows {
			return rate, err
		}
	}

	return rate, nil
}

func CommentRate(userId, commentId string) (LikedComment, error) {
	var rate LikedComment
	query := fmt.Sprintf("SELECT * FROM likedComments WHERE UserId LIKE '%s' AND CommentId LIKE '%s'", userId, commentId)
	err := Db.QueryRow(query).Scan(&rate.Id, &rate.Value, &rate.Comment.Id, &rate.User.Id)
	if err != nil {
		if err != sql.ErrNoRows {
			return rate, err
		}
	}

	return rate, nil
}

func DeleteLikedPost(userId, postId string, sql SQLDB) error {
	query := fmt.Sprintf("DELETE FROM likedPosts WHERE UserId LIKE '%s' AND PostId LIKE '%s'", userId, postId)
	_, err := sql.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func DeleteLikedComment(userId, commentId string, sql SQLDB) error {
	query := fmt.Sprintf("DELETE FROM likedComments WHERE UserId LIKE '%s' AND CommentId LIKE '%s'", userId, commentId)
	_, err := sql.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func LikedCommentCount(commentId, action string) (int, error) {
	var likes int
	query := fmt.Sprintf("SELECT COUNT(*) FROM likedComments WHERE CommentId LIKE '%s' AND Value LIKE '%s'", commentId, action)
	err := Db.QueryRow(query).Scan(&likes)
	if err != nil {
		return 0, err
	}

	return likes, nil
}

func LikedCommentsByCommentId(commentId string) ([]LikedComment, error) {
	var likes []LikedComment
	query := fmt.Sprintf("SELECT * FROM likedComments WHERE CommentId LIKE '%s'", commentId)
	rows, err := Db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		liked := LikedComment{}
		err := rows.Scan(&liked.Id, &liked.Value, &liked.Comment.Id, &liked.User.Id)
		if err != nil {
			return nil, err
		}

		liked.User, err = UserById(liked.User.Id)
		if err != nil {
			return nil, err
		}

		liked.Comment, err = CommentById(liked.Comment.Id)
		if err != nil {
			return nil, err
		}

		likes = append(likes, liked)
	}
	return likes, nil
}

func AddLikedComments(liked LikedComment, sql SQLDB) error {
	_, err := sql.Exec("INSERT INTO likedComments (Id,Value,CommentId,UserId) values ($1,$2,$3,$4)", liked.Id, liked.Value, liked.Comment.Id, liked.User.Id)
	if err != nil {
		return err
	}

	return nil
}
