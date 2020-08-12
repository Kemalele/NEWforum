package models

type Comment struct {
	Id string
	Description string
	PostDate string
	UserId string
	PostId string
}

func AllComments() ([]Comment,error) {
	var comments []Comment
	rows,err := Db.Query("SELECT * FROM Comment")
	if err != nil {
		return nil,err
	}
	for rows.Next() {
		comment := Comment{}
		err := rows.Scan(&comment.Id,&comment.Description,&comment.PostDate,&comment.UserId,&comment.PostId)
		if err != nil {
			return nil,err
		}
		comments = append(comments,comment)
	}
	return comments,nil
}

func AddComment(comment Comment,sql SQLDB) error{
	_,err := sql.Exec("INSERT INTO COMMENT (Id,Description,Post_date,UserId,PostId) values ($1,$2,$3,$4,$5)",comment.Id,comment.Description,comment.PostDate,comment.UserId,comment.PostId)
	if err != nil {
		return err
	}
	return nil
}