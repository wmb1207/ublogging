package mongodb

import (
	"fmt"
	"testing"

	"github.com/wmb1207/ublogging/internal/models"
)

const (
	TestDBName = "dev123"
	TestDBURI  = "mongodb://dev123:dev123@localhost:27017"
)

func TestConnection(t *testing.T) {
	_, err := NewMongoRepository(WithDBName(TestDBName), WithDBURI(TestDBURI))
	if err != nil {
		t.Error(err.Error())
	}
}

func TestCreateUser(t *testing.T) {
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
}

func TestGetUser(t *testing.T) {
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

	found, err := repo.User(unboxed.UUID)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(found)
	if found.Unbox().UUID != unboxed.UUID {
		t.Errorf("Error user with UUID: %s Not found", unboxed.UUID)
	}
}
