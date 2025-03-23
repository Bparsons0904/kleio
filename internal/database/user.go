package database

import (
	"database/sql"
	"log/slog"
)

type User struct {
	Username string
	Token    string
}

func (s *service) GetUser() (User, error) {
	var user User
	err := s.db.QueryRow("SELECT username, token FROM auth").Scan(&user.Username, &user.Token)
	if err != nil {
		if err != sql.ErrNoRows {
			slog.Error("Database query error", "error", err)
		}
		return User{}, err
	}

	return user, nil
}
