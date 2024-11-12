package mongodb

import (
	"context"
	"time"

	"github.com/wmb1207/ublogging/internal/models"
	"github.com/wmb1207/ublogging/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	User struct {
		UUID      *primitive.ObjectID  `bson:"_id,omitempty"`
		Email     string               `bson:"email"`
		Username  string               `bson:"username"`
		CreatedAt time.Time            `bson:"created_at"`
		UpdatedAt time.Time            `bson:"updated_at"`
		DeletedAt *time.Time           `bson:"deleted_at"`
		Following []primitive.ObjectID `bson:"folowing"`
		Followers []primitive.ObjectID `bson:"folowers"`
	}

	MongoUserRepository struct {
		baseRepository *MongoRepository
		userCollection *mongo.Collection
	}
)

func (u *User) Unbox() *models.User {
	return &models.User{
		UUID:     u.UUID.Hex(),
		Username: u.Username,
		Email:    u.Email,
	}
}

func FromUserModel(user *models.User) *User {
	getObjectID := func(uuid string) *primitive.ObjectID {
		oID, err := primitive.ObjectIDFromHex(uuid)
		var objectID *primitive.ObjectID
		if err != nil {
			objectID = nil
		} else {
			objectID = &oID
		}
		return objectID
	}

	return &User{
		UUID:     getObjectID(user.UUID),
		Email:    user.Email,
		Username: user.Username,
	}
}

func NewMongoUserRepository(mongoRepository *MongoRepository) *MongoUserRepository {
	return &MongoUserRepository{
		mongoRepository,
		mongoRepository.mongoClient.Database(mongoRepository.DBName).Collection("users"),
	}
}

func (m *MongoUserRepository) New(user *models.User) (repository.UserBox, error) {
	model := FromUserModel(user)
	inserted, err := m.userCollection.InsertOne(context.TODO(), model)
	if err != nil {
		return nil, err
	}

	insertedID := inserted.InsertedID.(primitive.ObjectID)
	model.UUID = &insertedID
	return model, nil
}

func (m *MongoUserRepository) User(uuid string) (repository.UserBox, error) {
	var user User

	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	if err := m.userCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return nil, m.handleError(err)
	}

	return &user, nil
}

func (m *MongoUserRepository) FindBy(options ...repository.FindUserWithOption) ([]repository.UserBox, error) {
	queryOptions := &repository.FindUserByOptions{}
	for _, option := range options {
		option(queryOptions)
	}

	filters := bson.M{}

	if queryOptions.UUID != "" {
		objectID, err := primitive.ObjectIDFromHex(queryOptions.UUID)
		if err != nil {
			return nil, err
		}
		filters["_id"] = objectID
	}

	// This should be by regex but well... let's make this simple for now
	if queryOptions.Email != "" {
		filters["email"] = queryOptions.Email
	}

	// This should be by regex but well... let's make this simple for now
	if queryOptions.Username != "" {
		filters["email"] = queryOptions.Email
	}

	output := []repository.UserBox{}

	cursor, err := m.userCollection.Find(context.TODO(), filters)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		output = append(output, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func (m *MongoUserRepository) Delete(user *User) error {
	// Not Implemented for this example
	return nil
}

func (m *MongoUserRepository) handleError(err error) error {
	if err == mongo.ErrNoDocuments {
		return NotFoundError{"No User found"}
	}

	return err
}
