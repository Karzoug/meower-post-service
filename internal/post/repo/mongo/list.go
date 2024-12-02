package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rs/xid"

	"github.com/Karzoug/meower-post-service/internal/post/entity"
)

// List returns a list of posts by the author ID.
func (pr postRepo) List(ctx context.Context, authorID, fromID xid.ID, limit int) ([]entity.Post, xid.ID, error) {
	filter := bson.D{
		{
			Key:   "author_id",
			Value: authorID,
		},
	}

	if !fromID.IsNil() {
		filter = append(filter, bson.E{
			Key:   "_id",
			Value: bson.D{{Key: "$gte", Value: fromID}},
		})
	}

	cursor, err := pr.client.
		Database(dbName).
		Collection(collectionName).
		Find(ctx, filter, options.Find().
			SetLimit(int64(limit+1)).
			SetSort(bson.D{{Key: "_id", Value: 1}}),
		)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []entity.Post{}, xid.NilID(), nil
		}
		return nil, xid.NilID(), err
	}
	defer cursor.Close(context.Background())

	var lastID xid.ID
	posts := make([]entity.Post, 0, limit)
	for i := 0; cursor.Next(ctx); i++ {
		var res entity.Post
		if err := cursor.Decode(&res); err != nil {
			return nil, xid.NilID(), err
		}

		if i == limit {
			lastID = res.ID
			break
		}

		posts = append(posts, res)
	}
	if err := cursor.Err(); err != nil {
		return nil, xid.NilID(), err
	}

	if len(posts) < limit {
		return posts, xid.NilID(), nil
	}

	return posts, lastID, nil
}

// ListIDProjection returns a list of post id projections by the author IDs.
func (pr postRepo) ListIDProjections(ctx context.Context, authorIDs []xid.ID, fromID xid.ID, limit int) ([]entity.PostIDProjection, xid.ID, error) {
	filter := bson.D{
		{
			Key: "author_id",
			Value: bson.M{
				"$in": authorIDs,
			},
		},
	}

	if !fromID.IsNil() {
		filter = append(filter, bson.E{
			Key:   "_id",
			Value: bson.D{{Key: "$gte", Value: fromID}},
		})
	}

	cursor, err := pr.client.
		Database(dbName).
		Collection(collectionName).
		Find(ctx, filter, options.Find().
			SetLimit(int64(limit+1)).
			SetSort(bson.D{{Key: "_id", Value: 1}}).
			SetProjection(bson.D{{Key: "author_id", Value: 1}}),
		)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []entity.PostIDProjection{}, xid.NilID(), nil
		}
		return nil, xid.NilID(), err
	}
	defer cursor.Close(context.Background())

	var lastID xid.ID
	posts := make([]entity.PostIDProjection, 0, limit)
	for i := 0; cursor.Next(ctx); i++ {
		var res entity.PostIDProjection
		if err := cursor.Decode(&res); err != nil {
			return nil, xid.NilID(), err
		}

		if i == limit {
			lastID = res.ID
			break
		}

		posts = append(posts, res)
	}
	if err := cursor.Err(); err != nil {
		return nil, xid.NilID(), err
	}

	if len(posts) < limit {
		return posts, xid.NilID(), nil
	}

	return posts, lastID, nil
}
