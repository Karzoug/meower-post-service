package service

import (
	"context"

	"github.com/rs/xid"

	"github.com/Karzoug/meower-post-service/internal/post/entity"
)

type PostService struct {
}

func NewPostService(repo repository) PostService {
	return PostService{}
}

func (ps PostService) CreatePost(ctx context.Context, post entity.Post) (entity.Post, error) {
	panic("implement me")
}

func (ps PostService) GetPost(ctx context.Context, id xid.ID) (entity.Post, error) {
	panic("implement me")
}

func (ps PostService) BatchGetPosts(ctx context.Context, ids []xid.ID) ([]entity.Post, error) {
	panic("implement me")
}

func (ps PostService) ListPosts(ctx context.Context, authorID xid.ID, pgn ListPostsPagination) (posts []entity.Post, nextID xid.ID, err error) {
	panic("implement me")
}

func (ps PostService) ListPostIDProjections(ctx context.Context, authorIDs []xid.ID, pgn ListPostIDProjectionsPagination) (projections []entity.PostIDProjection, nextID xid.ID, err error) {
	panic("implement me")
}
