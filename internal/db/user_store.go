package db

import (
	"database/sql"
	"fmt"
	entity "jwt-auth/internal/entity"
	"log"
	"time"
)

const (
	insert_query      = `INSERT INTO users (username, email, password, role, created_date, deleted_date) VALUES (?, ?, ?, ?, ?, NULL);`
	select_query      = `SELECT user_id, username, email, password, role, created_date, deleted_date FROM users`
	update_role_query = `UPDATE users set role = ? where email = ?`
	delete_user_query = `UPDATE users set deleted_date = ? where email = ?`
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
	select_by_email_query := fmt.Sprintf("%s %s", select_query, "WHERE email = ? AND deleted_date IS NULL;")
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
	rows, err := as.db.Query(fmt.Sprintf("%s %s", select_query, "WHERE deleted_date IS NULL;"))
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

func (as *UserStore) ChangeUsersRole(role, email string) (int, error) {
	res, err := as.db.Exec(update_role_query, role, email)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (as *UserStore) DeleteUser(email string) (int, error) {
	res, err := as.db.Exec(delete_user_query, time.Now(), email)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
