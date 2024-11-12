package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wmb1207/ublogging/internal/models"
)

func TestNewPostService(t *testing.T) {
	mockPostRepo := new(MockPostRepository)
	service := NewPostService(mockPostRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &PService{}, service)
}

func TestPService_New(t *testing.T) {
	mockPostRepo := new(MockPostRepository)
	service := &PService{PostRepository: mockPostRepo}

	user := &models.User{UUID: "1", Username: "user2", Email: "test@com"}
	content := "content 1"
	p1 := &models.Post{UUID: "1", User: user}
	p1.SetContent(content)

	newPost := &MockPostBox{data: *p1}
	mockPostRepo.On("New", mock.Anything).Return(newPost, nil)

	post, err := service.New("content 1", p1.User)

	assert.NoError(t, err)
	assert.Equal(t, content, post.Content())
	mockPostRepo.AssertExpectations(t)
}

func TestPService_Post(t *testing.T) {
	mockPostRepo := new(MockPostRepository)
	service := &PService{PostRepository: mockPostRepo}
	postUUID := "post1"
	expectedPost := &models.Post{UUID: postUUID}
	mockPostRepo.On("Post", postUUID).Return(&MockPostBox{data: *expectedPost}, nil)

	post, err := service.Post(postUUID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost.UUID, post.UUID)
	mockPostRepo.AssertExpectations(t)
}

func TestPService_Comment(t *testing.T) {
	mockPostRepo := new(MockPostRepository)
	service := &PService{PostRepository: mockPostRepo}

	comment := "Nice post!"

	p1 := &models.Post{UUID: "1", User: &models.User{UUID: "1", Username: "user2", Email: "test@com"}}
	p1.SetContent("content 1")
	p2 := &models.Post{UUID: "2", User: &models.User{UUID: "2", Username: "user2", Email: "test@com"}}
	p2.SetContent("content 2")

	newComment := p1.Comment(comment, p1.User)

	mockPostRepo.On("New", newComment).Return(&MockPostBox{data: *newComment}, nil)

	updatedPost, err := service.Comment(comment, p1, p1.User)

	assert.NoError(t, err)
	assert.Equal(t, p1, updatedPost)
	mockPostRepo.AssertExpectations(t)
}
