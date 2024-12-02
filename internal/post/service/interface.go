package service

import (
	"context"

	"github.com/rs/xid"

	"github.com/Karzoug/meower-post-service/internal/post/entity"
)

type repository interface {
	// Create creates a new post.
	Create(ctx context.Context, post *entity.Post) error
	// GetOne finds post by ID.
	GetOne(ctx context.Context, id xid.ID) (entity.Post, error)
	// GetMany finds posts by IDs, omitting any not found posts.
	GetMany(ctx context.Context, ids []xid.ID) ([]entity.Post, error)
	// List returns a list of posts by the author ID with pagination.
	// Deleted posts are omitted by default.
	List(ctx context.Context, authorID, fromID xid.ID, limit int,
	) (posts []entity.Post, nextID xid.ID, err error)
	// ListIDProjections returns a list of post id projections by the author IDs with pagination.
	// Deleted posts are omitted by default.
	ListIDProjections(ctx context.Context, authorIDs []xid.ID, fromID xid.ID, limit int,
	) (projections []entity.PostIDProjection, nextID xid.ID, err error)
}
