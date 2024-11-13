package repository

import "github.com/wmb1207/ublogging/internal/models"

type (
	FindUserByOptions struct {
		UUID         string
		UUIDS        []string
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
		Update(user *models.User, toUpdate map[string]interface{}) (UserBox, error)
	}
)

func WithUUID(uuid string) FindUserWithOption {
	return func(o *FindUserByOptions) {
		o.UUID = uuid
	}
}

func WithUUIDS(uuids []string) FindUserWithOption {
	return func(o *FindUserByOptions) {
		o.UUIDS = uuids
	}
}

func WithEmail(email string) FindUserWithOption {
	return func(o *FindUserByOptions) {
		o.Email = email
	}
}

func WithUsername(username string) FindUserWithOption {
	return func(o *FindUserByOptions) {
		o.Username = username
	}
}

func WithCreationTime(creationTime string) FindUserWithOption {
	return func(o *FindUserByOptions) {
		o.CreationTime = creationTime
	}
}
