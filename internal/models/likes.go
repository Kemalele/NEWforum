package models

import "fmt"

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

func LikedPostCount(postId string) (int, error) {
	var likes int
	query := fmt.Sprintf("SELECT COUNT(*) FROM likedPosts WHERE PostId LIKE '%s'", postId)
	err := Db.QueryRow(query).Scan(likes)
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

func DeleteLikedPost(userId, postId string, sql SQLDB) error {
	query := fmt.Sprintf("DELETE FROM likedPosts WHERE UserId LIKE '%s' AND PostId LIKE '%s'", userId, postId)
	_, err := sql.Exec(query)
	if err != nil {
		return err
	}

	return nil
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
