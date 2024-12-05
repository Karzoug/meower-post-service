package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rs/xid"

	"github.com/Karzoug/meower-post-service/internal/post/entity"
	"github.com/Karzoug/meower-post-service/internal/post/repo"
)

func (pr postRepo) Update(ctx context.Context, id xid.ID, updateFn func(post *entity.Post) error) error {
	post, err := pr.GetOne(ctx, id)
	if err != nil {
		return err
	}

	updatedAt := post.UpdatedAt

	if err := updateFn(&post); err != nil {
		return err
	}

	post.UpdatedAt = time.Now()

	filter := bson.M{
		"_id":        id,
		"updated_at": updatedAt,
	}

	res, err := pr.client.
		Database(dbName).
		Collection(collectionName).
		ReplaceOne(ctx, filter, post)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return repo.ErrAborted // post was updated between GetOne and ReplaceOne
	}

	return nil
}
