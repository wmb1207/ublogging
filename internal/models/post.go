package models

import "time"

type (
	PostType int

	ContentTooLongError struct {
		Length  int
		Message string
	}

	Like struct {
		User      User
		CreatedAt string
	}

	Post struct {
		UUID string
		User *User

		ParentUUID *string

		Type PostType

		content string
		likes   []Like
		repost  []Post

		createdAt string
		comments  []Post
	}
)

const (
	MaxPostContentLength          = 128
	TypePost             PostType = iota
	TypeComment
	TypeRepost
)

func NewPost(content string, user *User) *Post {
	return &Post{
		UUID:    "",
		User:    user,
		Type:    TypePost,
		content: content,
	}
}

func (c ContentTooLongError) Error() string {
	return c.Message

}

func (p *Post) Content() string {
	return p.content
}

func (p *Post) SetContent(content string) error {
	if len(content) > MaxPostContentLength {
		return ContentTooLongError{
			len(content),
			"Post Content exceded maximum length",
		}
	}
	p.content = content
	return nil
}

func (p *Post) RePost(user *User) {
	repost := Post{
		UUID: "",
		User: user,
		Type: TypeRepost,
	}

	p.repost = append(p.repost, repost)
}

func (p *Post) Like(user User) {

	currentTime := time.Now()
	// Format the timestamp as needed
	formattedTime := currentTime.Format(time.RFC3339)

	// Print the current timestamp

	like := Like{
		User:      user,
		CreatedAt: formattedTime,
	}

	p.likes = append(p.likes, like)
}

func (p *Post) Comment(content string, user *User) *Post {
	comment := Post{
		UUID:       "",
		User:       user,
		Type:       TypeComment,
		content:    content,
		ParentUUID: &p.UUID,
	}

	p.comments = append(p.comments, comment)
	return &comment
}
