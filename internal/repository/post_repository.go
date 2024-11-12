package repository

import "github.com/wmb1207/ublogging/internal/models"

type (
	FindPostByOptions struct {
		UUID         string
		User         *models.User
		UserUUID     string
		UserUUIDS    []string
		UserEmail    string
		UserUsername string
		ParentUUID   string
		PostType     models.PostType
	}

	FindPostWithOption func(*FindPostByOptions)

	PostBox interface {
		Unbox() *models.Post
	}

	PostRepository interface {
		Post(uuid string) (PostBox, error)
		New(post *models.Post) (PostBox, error)
		FindBy(options ...FindPostWithOption) ([]PostBox, error)
		Delete(post *models.Post) error
	}
)

func FindPostWithUserUUID(uuid string) FindPostWithOption {
	return func(f *FindPostByOptions) {
		f.UUID = uuid
	}
}

func FindPostWithUserUUIDInList(uuids []string) FindPostWithOption {
	return func(f *FindPostByOptions) {
		f.UserUUIDS = uuids
	}
}

func FindPostWithParentUUID(uuid string) FindPostWithOption {
	return func(f *FindPostByOptions) {
		f.ParentUUID = uuid
	}
}

func FindPostWithType(postType models.PostType) FindPostWithOption {
	return func(f *FindPostByOptions) {
		f.PostType = postType
	}
}