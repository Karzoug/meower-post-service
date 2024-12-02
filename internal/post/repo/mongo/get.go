package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/rs/xid"

	"github.com/Karzoug/meower-post-service/internal/post/entity"
	"github.com/Karzoug/meower-post-service/internal/post/repo"
)

func (pr postRepo) GetOne(ctx context.Context, id xid.ID) (entity.Post, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	var post entity.Post
	if err := pr.client.
		Database(dbName).
		Collection(collectionName).
		FindOne(ctx, filter).
		Decode(&post); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Post{}, repo.ErrRecordNotFound
		}
		return entity.Post{}, err
	}

	return post, nil
}

func (pr postRepo) GetMany(ctx context.Context, ids []xid.ID) (posts []entity.Post, err error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	cur, err := pr.client.
		Database(dbName).
		Collection(collectionName).
		Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	posts = make([]entity.Post, 0, len(ids))
	for cur.Next(ctx) {
		var post entity.Post
		if err := cur.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if cur.Err() != nil {
		return nil, cur.Err()
	}

	return posts, nil
}
