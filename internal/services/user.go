package services

import (
	"context"
	"fmt"
)

type UserService interface {
	Create(ctx context.Context, user UserData) (userId int, err error)
	SignIn(ctx context.Context, email, password string) (userId int, err error)
}

type UserStore interface {
	Create(ctx context.Context, user UserData) (userId int, err error)
	SignIn(ctx context.Context, email, password string) (userId int, err error)
}

type UserData struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type User struct {
	store UserStore
}

func NewUserService(store UserStore) User {
	return User{
		store: store,
	}
}

func (u User) Create(ctx context.Context, user UserData) (int, error) {
	userId, err := u.store.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("error while creating user : %w", err)
	}
	return userId, nil
}

func (u User) SignIn(ctx context.Context, email, password string) (int, error) {
	userId, err := u.store.SignIn(ctx, email, password)
	if err != nil {
		return 0, fmt.Errorf("error while signing : %w", err)
	}
	return userId, nil
}
