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
		Follow(user *models.User, toFollowUUID string) (*models.User, error)
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
	// First check if the user exists on the DB
	// I've implemented a really simple and BAD user management.
	// This is just for testing
	userBox, err := u.UserRepository.New(user)
	if err != nil {
		return nil, err
	}

	newUser := userBox.Unbox()

	return newUser, nil
}

func (u *UService) Follow(user *models.User, toFollowUUID string) (*models.User, error) {
	knownFollowing := []string{}

	for _, following := range user.Following {
		knownFollowing = append(knownFollowing, following.UUID)
	}

	ubox, err := u.UserRepository.Update(user, map[string]interface{}{
		"following": append(knownFollowing, toFollowUUID),
	})
	if err != nil {
		return nil, err
	}

	return ubox.Unbox(), nil
}

func (u *UService) Feed(user *models.User) ([]*models.Post, error) {
	// TODO: Add pagination.
	// Simulate a feed when no followers exists  by just getting all of them
	followingUUIDS := []string{}
	for _, u := range user.Following {
		followingUUIDS = append(followingUUIDS, u.UUID)
	}

	var foundFeed []repository.PostBox

	if len(followingUUIDS) > 0 {
		feed, err := u.PostRepository.FindBy(repository.FindPostWithUserUUIDInList(followingUUIDS))
		if err != nil {
			return nil, err
		}
		foundFeed = feed
	} else {
		feed, err := u.PostRepository.FindBy(repository.FindPostWithUserNotInUUIDLists([]string{user.UUID}))
		if err != nil {
			return nil, err
		}
		foundFeed = feed
	}

	var output []*models.Post
	for _, post := range foundFeed {
		unboxed := post.Unbox()
		output = append(output, unboxed)
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
