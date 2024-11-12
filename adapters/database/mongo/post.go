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
	Post struct {
		UUID     *primitive.ObjectID `bson:"_id,omitempty"`
		UserID   primitive.ObjectID  `bson:"user_id"`
		Content  string              `bson:"content"`
		Type     int                 `bson:"type"`
		ParentID *primitive.ObjectID `bson:"parent_id"`

		User     *User                `bson:"user"`
		Likes    []primitive.ObjectID `bson:"likes"`
		Comments []primitive.ObjectID `bson:"comments"`

		CreatedAt time.Time  `bson:"created_at"`
		UpdatedAt time.Time  `bson:"updated_at"`
		DeletedAt *time.Time `bson:"deleted_at"`
	}

	PostWithUser struct {
		Post Post `bson:"post"`
		User User `bson:"user"`
	}

	MongoPostRepository struct {
		baseRepository *MongoRepository
		postCollection *mongo.Collection
	}
)

func FromPostModel(post *models.Post) *Post {
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

	return &Post{
		UUID:    getObjectID(post.UUID),
		UserID:  *getObjectID(post.User.UUID),
		Content: post.Content(),
		Type:    int(post.Type),
	}
}

func (p *Post) Unbox() *models.Post {
	return &models.Post{
		UUID: p.UUID.Hex(),
	}
}

func (p *PostWithUser) Unbox() *models.Post {
	post := p.Post.Unbox()
	user := p.User.Unbox()
	post.User = user
	return post
}

func NewMongoPostRepository(mongoRepository *MongoRepository) *MongoPostRepository {
	return &MongoPostRepository{
		mongoRepository,
		mongoRepository.mongoClient.Database(mongoRepository.DBName).Collection("post"),
	}
}

func (m *MongoPostRepository) New(post *models.Post) (repository.PostBox, error) {
	model := FromPostModel(post)
	inserted, err := m.postCollection.InsertOne(context.TODO(), model)
	if err != nil {
		return nil, err
	}

	insertedID := inserted.InsertedID.(primitive.ObjectID)
	model.UUID = &insertedID
	return model, nil
}

func (m *MongoPostRepository) Post(uuid string) (repository.PostBox, error) {

	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{"_id": objectID}},
		},
		{
			{"$lookup", bson.M{
				"from":         "users",
				"localfield":   "user_id",
				"foreignField": "_id",
				"as":           "user",
			}},
		},
		{
			{"$unwind", bson.M{"path": "$user"}},
		},
	}

	cursor, err := m.postCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	results := []*PostWithUser{}
	for cursor.Next(context.TODO()) {
		var result PostWithUser
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results[0], nil
}

func (m *MongoPostRepository) FindBy(options ...repository.FindPostWithOption) ([]repository.PostBox, error) {
	// TODO: Implement this. Not required right now.

	queryOptions := &repository.FindPostByOptions{}

	for _, option := range options {
		option(queryOptions)
	}

	var pipeline mongo.Pipeline
	var output []repository.PostBox

	query := bson.D{}

	if queryOptions.UUID != "" {
		query = append(query, bson.E{Key: "_id", Value: queryOptions.UUID})
	}

	if queryOptions.ParentUUID != "" {
		objectID, err := primitive.ObjectIDFromHex(queryOptions.ParentUUID)
		if err != nil {
			return nil, err
		}

		filter := bson.E{Key: "parent_id", Value: objectID}
		query = append(query, filter)
	}

	if queryOptions.PostType != 0 {
		query = append(query, bson.E{Key: "type", Value: int(queryOptions.PostType)})
	}

	if len(queryOptions.UserUUIDS) > 0 {

		objectIDS, err := m.toObjectIDs(queryOptions.UserUUIDS)
		if err != nil {
			return nil, err
		}

		pipeline = mongo.Pipeline{
			{{"$match", bson.D{{"user_id", bson.D{{"$in", objectIDS}}}}}},
			{{"$lookup", bson.D{
				{"from", "users"},
				{"localField", "user_id"},
				{"foreignField", "_id"},
				{"as", "user"},
			}}},
			{
				{"$unwind", bson.M{"path": "$user"}},
			},
		}

		cursor, err := m.postCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			return nil, err
		}

		for cursor.Next(context.TODO()) {
			var result PostWithUser
			if err := cursor.Decode(&result); err != nil {
				return nil, err
			}

			output = append(output, &result)
		}

		if err := cursor.Err(); err != nil {
			return nil, err
		}

		return output, nil
	}

	cursor, err := m.postCollection.Find(context.TODO(), query)

	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var result Post
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		output = append(output, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func (m *MongoPostRepository) Delete(post *models.Post) error {
	// TODO: Implement this, Not required right now
	return nil
}

func (m *MongoPostRepository) toObjectIDs(uuids []string) ([]primitive.ObjectID, error) {
	output := []primitive.ObjectID{}
	for _, uuid := range uuids {
		objectID, err := primitive.ObjectIDFromHex(uuid)
		if err != nil {
			return nil, err
		}

		output = append(output, objectID)
	}

	return output, nil

}
