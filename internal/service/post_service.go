package service

import (
	"github.com/wmb1207/ublogging/internal/models"
	"github.com/wmb1207/ublogging/internal/repository"
)

type (
	PostService interface {
		New(content string, user *models.User) (*models.Post, error)
		Post(uuid string) (*models.Post, error)
		Comment(comment string, post *models.Post, user *models.User) (*models.Post, error)
		Like(post *models.Post, user *models.User) (*models.Post, error)
		Repost(post *models.Post, user *models.User) (*models.Post, error)
		Comments(post *models.Post) ([]models.Post, error)
	}

	PService struct {
		PostRepository repository.PostRepository
	}
)

func NewPostService(postRepo repository.PostRepository) PostService {
	return &PService{
		PostRepository: postRepo,
	}
}

func (p *PService) New(content string, user *models.User) (*models.Post, error) {
	newPost := models.NewPost(content, user)
	output, err := p.PostRepository.New(newPost)
	if err != nil {
		return nil, err
	}

	return output.Unbox(), nil
}

func (p *PService) Post(uuid string) (*models.Post, error) {
	post, err := p.PostRepository.Post(uuid)
	if err != nil {
		return nil, err
	}

	return post.Unbox(), nil
}

func (p *PService) Comment(comment string, post *models.Post, user *models.User) (*models.Post, error) {
	newComment := post.Comment(comment, user)
	if _, err := p.PostRepository.New(newComment); err != nil {
		return nil, err
	}
	return post, nil
}

func (p *PService) Comments(post *models.Post) ([]models.Post, error) {
	comments, err := p.PostRepository.FindBy(
		repository.FindPostWithParentUUID(post.UUID),
		repository.FindPostWithType(models.TypeComment),
	)
	if err != nil {
		return nil, err
	}

	output := []models.Post{}
	for _, comment := range comments {
		output = append(output, *comment.Unbox())
	}

	return output, nil
}

func (p *PService) Like(post *models.Post, user *models.User) (*models.Post, error) {
	// TODO: Implement this
	return nil, nil
}

func (p *PService) Repost(post *models.Post, user *models.User) (*models.Post, error) {
	// TODO: IMPLEMENT THIS
	return nil, nil
}
