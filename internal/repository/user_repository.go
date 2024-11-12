package repository

import "github.com/wmb1207/ublogging/internal/models"

type (
	FindUserByOptions struct {
		UUID         string
		Email        string
		Username     string
		CreationTime string
	}

	UserBox interface {
		Unbox() *models.User
	}

	FindUserWithOption func(*FindUserByOptions)

	UserRepository interface {
		User(uuid string) (UserBox, error)
		New(user *models.User) (UserBox, error)
		FindBy(options ...FindUserWithOption) ([]UserBox, error)
		Delete(user *models.User) error
	}
)
