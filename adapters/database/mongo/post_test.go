package mongodb

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/wmb1207/ublogging/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreatePost(t *testing.T) {
	baseRepo, err := NewMongoRepository(WithDBName(TestDBName), WithDBURI(TestDBURI))
	if err != nil {
		t.Errorf(err.Error())
	}

	repo := NewMongoUserRepository(baseRepo)
	ubox, err := repo.New(&models.User{
		Email:    "testUser@email.com",
		Username: "testuser123",
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	unboxed := ubox.Unbox()
	fmt.Println(unboxed.UUID)

	if ubox.Unbox().UUID == "" {
		t.Errorf("error unboxing user -> missing UUID")
	}

	postRepo := NewMongoPostRepository(baseRepo)
	userObjectID, err := primitive.ObjectIDFromHex(unboxed.UUID)
	if err != nil {
		t.Errorf(err.Error())
	}
	post := &models.Post{
		User: &models.User{UUID: userObjectID.Hex(), Username: "", Email: ""},
	}
	post.SetContent("TEST CONTENT 1")

	pbox, err := postRepo.New(post)

	if err != nil {
		t.Errorf(err.Error())
	}

	if pbox.Unbox().UUID == "" {
		t.Errorf("error unboxing post -> missing UUID")
	}
}

func TestListPosts(t *testing.T) {
	baseRepo, err := NewMongoRepository(WithDBName(TestDBName), WithDBURI(TestDBURI))
	if err != nil {
		t.Errorf(err.Error())
	}

	repo := NewMongoUserRepository(baseRepo)
	ubox, err := repo.New(&models.User{
		Email:    "testUser@email.com",
		Username: "testuser123",
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	unboxed := ubox.Unbox()
	fmt.Println(unboxed.UUID)

	if ubox.Unbox().UUID == "" {
		t.Errorf("error unboxing user -> missing UUID")
	}

	postRepo := NewMongoPostRepository(baseRepo)
	userObjectID, err := primitive.ObjectIDFromHex(unboxed.UUID)
	if err != nil {
		t.Errorf(err.Error())
	}

	posts := make([]models.Post, 10)

	for i := 0; i < 10; i++ {
		post := &models.Post{
			User: &models.User{UUID: userObjectID.Hex(), Username: "", Email: ""},
		}
		post.SetContent("TEST_CONTENT " + strconv.Itoa(i+1))
		posts[i] = *post

	}

	for _, post := range posts {
		pbox, err := postRepo.New(&post)

		if err != nil {
			t.Errorf(err.Error())
		}

		if pbox.Unbox().UUID == "" {
			t.Errorf("error unboxing post -> missing UUID")
		}

	}

}
