package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/xid"
	"google.golang.org/grpc/codes"

	"github.com/Karzoug/meower-common-go/auth"
	"github.com/Karzoug/meower-common-go/ucerr"

	"github.com/Karzoug/meower-post-service/internal/post/entity"
	repoerr "github.com/Karzoug/meower-post-service/internal/post/repo"
	"github.com/Karzoug/meower-post-service/pkg/validator"
)

type PostService struct {
	repo repository
}

// NewPostService creates new PostService.
func NewPostService(repo repository) PostService {
	return PostService{
		repo: repo,
	}
}

// CreatePost creates a new post.
func (ps PostService) CreatePost(ctx context.Context, post entity.Post) (entity.Post, error) {
	if err := validator.Struct(post); err != nil {
		return entity.Post{}, ucerr.NewError(
			err,
			"invalid post format: "+err.Error(),
			codes.InvalidArgument,
		)
	}

	if auth.UserIDFromContext(ctx).Compare(post.AuthorID) != 0 {
		return entity.Post{}, ucerr.NewError(
			nil,
			"cannot create post for another user",
			codes.PermissionDenied,
		)
	}

	if err := ps.repo.Create(ctx, &post); err != nil {
		return entity.Post{}, ucerr.NewInternalError(fmt.Errorf("repo error: %w", err))
	}

	return post, nil
}

// GetPost finds post by ID.
func (ps PostService) GetPost(ctx context.Context, id xid.ID) (entity.Post, error) {
	post, err := ps.repo.GetOne(ctx, id)
	if err != nil {
		if errors.Is(err, repoerr.ErrRecordNotFound) {
			return entity.Post{}, ucerr.NewError(
				nil,
				"not found post",
				codes.NotFound,
			)
		}
		return entity.Post{}, ucerr.NewInternalError(fmt.Errorf("repo error: %w", err))
	}

	if post.IsDeleted && auth.UserIDFromContext(ctx).Compare(post.AuthorID) != 0 {
		hideDataOfDeletedPost(&post)
	}

	return post, nil
}

func (ps PostService) DeletePost(ctx context.Context, id xid.ID) error {
	update := func(post *entity.Post) error {
		if auth.UserIDFromContext(ctx).Compare(post.AuthorID) != 0 {
			return ucerr.NewError(
				nil,
				"cannot delete post for another user",
				codes.PermissionDenied,
			)
		}

		if post.IsDeleted {
			return repoerr.ErrNoAffected
		}
		post.IsDeleted = true

		return nil
	}

	if err := ps.repo.Update(ctx, id, update); err != nil {
		switch {
		case errors.Is(err, repoerr.ErrNoAffected):
			return nil
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return ucerr.NewError(
				nil,
				"not found post",
				codes.NotFound,
			)
		case errors.Is(err, repoerr.ErrAborted):
			return ucerr.NewError(
				nil,
				"operation aborted: post not affected, try again",
				codes.Aborted,
			)
		default:
			return ucerr.NewInternalError(fmt.Errorf("repo error: %w", err))
		}
	}

	return nil
}

// BatchGetPosts finds posts by IDs.
func (ps PostService) BatchGetPosts(ctx context.Context, ids []xid.ID) ([]entity.Post, error) {
	posts, err := ps.repo.GetMany(ctx, ids)
	if err != nil {
		return nil, ucerr.NewInternalError(fmt.Errorf("repo error: %w", err))
	}

	for i := range posts {
		if posts[i].IsDeleted {
			hideDataOfDeletedPost(&posts[i])
		}
	}

	if len(posts) != len(ids) {
		return nil, ucerr.NewError(
			nil,
			"not found post(-s)",
			codes.NotFound,
		)
	}

	return posts, nil
}

// ListPosts returns a list of posts by the author ID with pagination.
// Deleted posts are omitted by default.
func (ps PostService) ListPosts(ctx context.Context, authorID xid.ID, pgn ListPostsPagination) (posts []entity.Post, nextID xid.ID, err error) {
	if pgn.Size < 0 {
		return nil, xid.NilID(), ucerr.NewError(
			nil,
			"invalid pagination parameter: negative size",
			codes.InvalidArgument,
		)
	}
	if pgn.Size == 0 {
		pgn.Size = 100
	} else if pgn.Size > 100 {
		pgn.Size = 100
	}

	posts, nextID, err = ps.repo.List(ctx, authorID, pgn.Token, pgn.Size)
	if err != nil {
		return nil, xid.NilID(), ucerr.NewInternalError(fmt.Errorf("repo error: %w", err))
	}

	return posts, nextID, nil
}

// ListPostIDProjections returns a list of post id projections by the author IDs with pagination.
// Deleted posts are omitted by default.
func (ps PostService) ListPostIDProjections(ctx context.Context, authorIDs []xid.ID, pgn ListPostIDProjectionsPagination) (projections []entity.PostIDProjection, nextID xid.ID, err error) {
	if pgn.Size < 0 {
		return nil, xid.NilID(), ucerr.NewError(
			nil,
			"invalid pagination parameter: negative size",
			codes.InvalidArgument,
		)
	}
	if pgn.Size == 0 {
		pgn.Size = 100
	} else if pgn.Size > 1000 {
		pgn.Size = 1000
	}

	projections, nextID, err = ps.repo.ListIDProjections(ctx, authorIDs, pgn.Token, pgn.Size)
	if err != nil {
		return nil, xid.NilID(), ucerr.NewInternalError(fmt.Errorf("repo error: %w", err))
	}

	return projections, nextID, nil
}

// hideDataOfDeletedPost prevents leak of deleted post data.
func hideDataOfDeletedPost(post *entity.Post) {
	post.Text = "(deleted by user)"
}
