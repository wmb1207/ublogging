package service

import (
	"github.com/stretchr/testify/mock"
	"github.com/wmb1207/ublogging/internal/models"
	"github.com/wmb1207/ublogging/internal/repository"
)

type (
	MockUserBox struct {
		data models.User
	}

	MockPostBox struct {
		data models.Post
	}

	MockUserRepository struct {
		mock.Mock
	}
	MockPostRepository struct {
		mock.Mock
	}
)

func (m *MockUserRepository) New(user *models.User) (repository.UserBox, error) {
	args := m.Called(user)
	return args.Get(0).(repository.UserBox), args.Error(1)
}

func (m *MockUserRepository) User(uuid string) (repository.UserBox, error) {
	args := m.Called(uuid)
	return args.Get(0).(repository.UserBox), args.Error(1)
}

func (m *MockUserRepository) FindBy(options ...repository.FindUserWithOption) ([]repository.UserBox, error) {
	args := m.Called(options)
	return args.Get(0).([]repository.UserBox), args.Error(1)
}

func (m *MockUserRepository) Delete(user *models.User) error {
	args := m.Called(user)
	return args.Error(1)
}

func (m *MockPostRepository) FindBy(options ...repository.FindPostWithOption) ([]repository.PostBox, error) {
	args := m.Called(options)
	return args.Get(0).([]repository.PostBox), args.Error(1)
}

func (m *MockPostRepository) Post(uuid string) (repository.PostBox, error) {
	args := m.Called(uuid)
	return args.Get(0).(repository.PostBox), args.Error(1)
}

func (m *MockPostRepository) New(post *models.Post) (repository.PostBox, error) {
	args := m.Called(post)
	return args.Get(0).(repository.PostBox), args.Error(1)
}

func (m *MockPostRepository) Delete(post *models.Post) error {
	args := m.Called(post)
	return args.Error(1)
}

func (m *MockUserBox) Unbox() *models.User {
	return &m.data
}

func (m *MockPostBox) Unbox() *models.Post {
	return &m.data
}
