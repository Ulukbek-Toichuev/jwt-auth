package db

import (
	"database/sql"
	"fmt"
	entity "jwt-auth/internal/entity"
	"log"
)

const (
	insert_query = `INSERT INTO users (username, email, password, role, created_date, deleted_date) VALUES (?, ?, ?, ?, ?, NULL);`
	select_query = `SELECT user_id, username, email, password, role, created_date, deleted_date FROM users`
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db}
}

func (as *UserStore) CreateUser(user entity.UserEntity) (int, error) {
	res, err := as.db.Exec(insert_query, user.Username, user.Email, user.Password, user.Role, user.CreatedDate)
	if err != nil {
		log.Printf("%v", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("%v", err)
		return 0, err
	}
	return int(id), nil
}

func (as *UserStore) GetUserByEmail(email string) (entity.UserEntity, error) {
	select_by_email_query := fmt.Sprintf("%s %s", select_query, "WHERE email = ?;")
	row := as.db.QueryRow(select_by_email_query, email)
	var user entity.UserEntity
	err := row.Scan(&user.UserId, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedDate, &user.DeletedDate)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (as *UserStore) GetUsers() ([]entity.UserEntity, error) {
	result := make([]entity.UserEntity, 0)
	rows, err := as.db.Query(select_query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user entity.UserEntity
		err := rows.Scan(&user.UserId, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedDate, &user.DeletedDate)
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}

	return result, nil
}

func (as *UserStore) DeleteUser(id int) {

}
