package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPost(t *testing.T) {
	user := User{UUID: "user-123"}
	post := NewPost("Hello, World!", &user)

	assert.Equal(t, post.User.UUID, user.UUID)
	assert.Equal(t, post.content, "Hello, World!")

}

func TestSetContentTooLong(t *testing.T) {
	user := User{UUID: "user-123"}
	post := NewPost("Short content", &user)

	err := post.SetContent("This content is way too long and exceeds the maximum allowed length of 128 characters, which is something we really want to avoid in our application.")

	var contentTooLongError ContentTooLongError
	assert.Error(t, err)
	assert.True(t, errors.As(err, &contentTooLongError))
}

func TestLike(t *testing.T) {
	user1 := User{UUID: "user-123"}
	user2 := User{UUID: "user-456"}
	post := NewPost("Content", &user1)

	post.Like(user2)

	assert.Equal(t, len(post.likes), 1)
	assert.Equal(t, post.likes[0].User.UUID, user2.UUID)
}

func TestComment(t *testing.T) {
	user := User{UUID: "user-123"}
	post := NewPost("Content", &user)

	post.Comment("This is a comment", &user)

	assert.Equal(t, len(post.comments), 1)
	assert.Equal(t, post.comments[0].content, "This is a comment")

}
