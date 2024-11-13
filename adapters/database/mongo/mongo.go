package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MissingDBURIError  struct{}
	MissingDBNameError struct{}

	NotFoundError struct {
		Detail string
	}

	MongoRepositoryConfig func(*MongoRepository)

	MongoRepository struct {
		DBURI       string
		DBName      string
		mongoClient *mongo.Client
	}
)

func (m MissingDBNameError) Error() string {
	return "Missing DB Name during MongoRepository initialization"
}

func (m MissingDBURIError) Error() string {
	return "Missing dbURI Name during MongoRepository initialization"
}

func (n NotFoundError) Error() string {
	return n.Detail
}

func WithDBURI(uri string) MongoRepositoryConfig {
	return func(mr *MongoRepository) {
		mr.DBURI = uri
	}
}

func WithDBName(dbName string) MongoRepositoryConfig {
	return func(mr *MongoRepository) {
		mr.DBName = dbName
	}
}

func NewMongoRepository(buildOptions ...MongoRepositoryConfig) (*MongoRepository, error) {
	mr := &MongoRepository{}

	for _, option := range buildOptions {
		option(mr)
	}

	if mr.DBName == "" {
		return nil, MissingDBNameError{}
	}

	if mr.DBURI == "" {
		return nil, MissingDBURIError{}
	}

	clientOptions := options.Client().ApplyURI(mr.DBURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	mr.mongoClient = client

	return mr, nil
}

func toObjectIDs(uuids []string) ([]primitive.ObjectID, error) {
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
