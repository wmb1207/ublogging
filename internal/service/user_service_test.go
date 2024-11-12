package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wmb1207/ublogging/internal/models"
	"github.com/wmb1207/ublogging/internal/repository"
)

func TestNewUserService(t *testing.T) {
	userRepo := new(MockUserRepository)
	postRepo := new(MockPostRepository)
	service := NewUserService(userRepo, postRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &UService{}, service)
}

func TestUService_New(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockPostRepo := new(MockPostRepository)
	service := &UService{UserRepository: mockUserRepo, PostRepository: mockPostRepo}

	user := &models.User{UUID: "user1"}

	mockUserRepo.On("New", user).Return(&MockUserBox{data: *user}, nil)

	createdUser, err := service.New(user)

	assert.NoError(t, err)
	assert.Equal(t, user.UUID, createdUser.UUID)
	mockUserRepo.AssertExpectations(t)
}

func TestUService_Feed(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockPostRepo := new(MockPostRepository)
	service := &UService{UserRepository: mockUserRepo, PostRepository: mockPostRepo}

	user := &models.User{Following: []*models.User{{UUID: "user2"}, {UUID: "user3"}}}

	p1 := &models.Post{UUID: "1", User: &models.User{UUID: "1", Username: "user2", Email: "test@com"}}
	p1.SetContent("content 1")
	p2 := &models.Post{UUID: "2", User: &models.User{UUID: "2", Username: "user2", Email: "test@com"}}
	p2.SetContent("content 2")

	posts := []repository.PostBox{
		&MockPostBox{data: *p1}, &MockPostBox{data: *p2},
	}

	mockPostRepo.On("FindBy", mock.Anything).Return(posts, nil)

	feed, err := service.Feed(user)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(feed))
	assert.Equal(t, "content 1", feed[0].Content())
	assert.Equal(t, "content 2", feed[1].Content())
	mockPostRepo.AssertExpectations(t)
}

func TestUService_User(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockPostRepo := new(MockPostRepository)
	service := &UService{UserRepository: mockUserRepo, PostRepository: mockPostRepo}

	user := &models.User{UUID: "user1", Following: []*models.User{{UUID: "user2"}}}
	mockUserRepo.On("User", "user1").Return(&MockUserBox{data: *user}, nil)
	p1 := &models.Post{UUID: "1", User: &models.User{UUID: "1", Username: "user2", Email: "test@com"}}
	p1.SetContent("content 1")
	posts := []repository.PostBox{
		&MockPostBox{data: *p1},
	}

	mockPostRepo.On("FindBy", mock.Anything).Return(posts, nil)

	resultUser, err := service.User("user1")

	assert.NoError(t, err)
	assert.Equal(t, user.UUID, resultUser.UUID)
	assert.Equal(t, 1, len(resultUser.Feed))
	assert.Equal(t, "content 1", resultUser.Feed[0].Content())
	mockUserRepo.AssertExpectations(t)
	mockPostRepo.AssertExpectations(t)
}

func TestUService_New_Error(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockPostRepo := new(MockPostRepository)
	service := &UService{UserRepository: mockUserRepo, PostRepository: mockPostRepo}

	user := &models.User{UUID: "user1"}

	mockUserRepo.On("New", user).Return(&MockUserBox{}, errors.New("error creating user"))

	createdUser, err := service.New(user)

	assert.Error(t, err)
	assert.Nil(t, createdUser)
	mockUserRepo.AssertExpectations(t)
}

func TestUService_User_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockPostRepo := new(MockPostRepository)
	service := &UService{UserRepository: mockUserRepo, PostRepository: mockPostRepo}

	mockUserRepo.On("User", "user1").Return(&MockUserBox{}, errors.New("user not found"))

	user, err := service.User("user1")

	assert.Error(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertExpectations(t)
}
