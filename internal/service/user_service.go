package service

import (
	"github.com/wmb1207/ublogging/internal/models"
	"github.com/wmb1207/ublogging/internal/repository"
)

type (
	UserService interface {
		New(user *models.User) (*models.User, error)
		Feed(user *models.User) ([]*models.Post, error)
		User(uuid string) (*models.User, error)
	}

	UService struct {
		UserRepository repository.UserRepository
		PostRepository repository.PostRepository
	}
)

func NewUserService(userRepo repository.UserRepository, postRepo repository.PostRepository) UserService {
	return &UService{
		userRepo,
		postRepo,
	}
}

func (u *UService) New(user *models.User) (*models.User, error) {
	userBox, err := u.UserRepository.New(user)
	if err != nil {
		return nil, err
	}

	newUser := userBox.Unbox()

	return newUser, nil
}

func (u *UService) Feed(user *models.User) ([]*models.Post, error) {
	// TODO: Add pagination.
	followingUUIDS := []string{}
	for _, u := range user.Following {
		followingUUIDS = append(followingUUIDS, u.UUID)
	}

	feed, err := u.PostRepository.FindBy(repository.FindPostWithUserUUIDInList(followingUUIDS))
	if err != nil {
		return nil, err
	}

	var output []*models.Post
	for _, post := range feed {
		output = append(output, post.Unbox())
	}

	return output, nil
}

func (u *UService) User(uuid string) (*models.User, error) {
	ubox, err := u.UserRepository.User(uuid)
	if err != nil {
		return nil, err
	}

	user := ubox.Unbox()

	feed, err := u.Feed(user)

	if err != nil {
		return nil, err
	}

	user.Feed = feed
	return user, nil
}
