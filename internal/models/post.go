package models

import "time"

type (
	PostType int

	ContentTooLongError struct {
		Length  int
		Message string
	}

	Like struct {
		User      User   `json:"user"`
		CreatedAt string `json:created_at"`
	}

	Post struct {
		UUID        string   `json:"uuid"`
		User        *User    `json:"user"`
		ParentUUID  *string  `json:"parent_uuid,omitempty"`
		Type        PostType `json:"type"`
		ContentData string   `json:"content"`
		Likes       []Like   `json:"likes"`
		Repost      []Post   `json:"_,omitempty"`
		CreatedAt   string   `json:"created_at"`
		Comments    []Post   `json:"comments,"`
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
		UUID:        "",
		User:        user,
		Type:        TypePost,
		ContentData: content,
	}
}

func (c ContentTooLongError) Error() string {
	return c.Message

}

func (p *Post) Content() string {
	return p.ContentData
}

func (p *Post) SetContent(content string) error {
	if len(content) > MaxPostContentLength {
		return ContentTooLongError{
			len(content),
			"Post Content exceded maximum length",
		}
	}
	p.ContentData = content
	return nil
}

func (p *Post) RePost(user *User) {
	repost := Post{
		UUID: "",
		User: user,
		Type: TypeRepost,
	}

	p.Repost = append(p.Repost, repost)
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

	p.Likes = append(p.Likes, like)
}

func (p *Post) Comment(content string, user *User) *Post {
	comment := Post{
		UUID:        "",
		User:        user,
		Type:        TypeComment,
		ContentData: content,
		ParentUUID:  &p.UUID,
	}

	p.Comments = append(p.Comments, comment)
	return &comment
}
