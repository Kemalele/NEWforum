package models

type Liked struct {
	Id string
	Value string
	PostId string
	UserId string
}

func AllLikes() ([]Liked,error) {
	var likes []Liked
	rows,err := Db.Query("SELECT * FROM Liked")
	if err != nil {
		return nil,err
	}

	for rows.Next() {
		liked := Liked{}
		err := rows.Scan(&liked.Id,&liked.Value,&liked.PostId,&liked.UserId)
		if err != nil {
			return nil,err
		}
		likes = append(likes,liked)
	}
	return likes,nil
}

func AddLike(liked Liked,sql SQLDB) error{
	_,err := sql.Exec("INSERT INTO LIKED (Id,Value,PostId,UserId) values ($1,$2,$3,$4)",liked.Id,liked.Value,liked.PostId,liked.UserId)
	if err != nil {
		return err
	}
	return nil
}
