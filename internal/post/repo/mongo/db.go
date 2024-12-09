package mongo

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	dbName         = "meower"
	collectionName = "posts"
)

type postRepo struct {
	client *mongo.Client
}

func NewPostRepo(mongoClient *mongo.Client) postRepo {
	return postRepo{client: mongoClient}
}
