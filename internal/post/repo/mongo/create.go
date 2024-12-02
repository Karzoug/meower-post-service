package mongo

import (
	"context"

	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Karzoug/meower-post-service/internal/post/entity"
	"github.com/Karzoug/meower-post-service/internal/post/repo"
)

func (pr postRepo) Create(ctx context.Context, post *entity.Post) error {
	post.ID = xid.New()
	post.UpdatedAt = post.ID.Time()

	_, err := pr.client.
		Database(dbName).
		Collection(collectionName).
		InsertOne(ctx, post)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return repo.ErrRecordAlreadyExists
		}
		return err
	}

	return nil
}
