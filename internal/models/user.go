package models

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mahdihagh80/forms/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	db *sqlx.DB
}

func NewUserModel(db *sqlx.DB) UserModel {
	return UserModel{
		db: db,
	}
}

func (um UserModel) Create(ctx context.Context, user services.UserData) (int, error) {
	query := "INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)"
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("error while hashing user password : %w", err)
	}
	
	res, err := um.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("error while inserting new user to db : %w", err)
	}
	
	userId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error while inserting new user to db and get his userId: %w", err)
	}
	return int(userId), nil
}

func (um UserModel) SignIn(ctx context.Context, email, password string) (int, error) {
	query := "SELECT id, password FROM users WHERE email=?"
	var userId int
	var hashedPassword string

	err := um.db.QueryRowContext(ctx, query, email).Scan(&userId, &hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("error while fetching user data from db : %w", err)
	}

	err = checkPassword(hashedPassword, password)
	if err != nil {
		return 0, fmt.Errorf("password mismatched : %w", err)
	}
	return userId, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
