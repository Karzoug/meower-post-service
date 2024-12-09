package post

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/rs/xid"

	"github.com/Karzoug/meower-common-go/auth"
	"github.com/Karzoug/meower-post-service/internal/delivery/grpc/converter"
	gen "github.com/Karzoug/meower-post-service/internal/delivery/grpc/gen/post/v1"
	"github.com/Karzoug/meower-post-service/internal/post/service"
)

func RegisterService(ps service.PostService) func(grpcServer *grpc.Server) {
	hdl := handlers{
		postService: ps,
	}
	return func(grpcServer *grpc.Server) {
		gen.RegisterPostServiceServer(grpcServer, hdl)
	}
}

type handlers struct {
	gen.UnimplementedPostServiceServer
	postService service.PostService
}

func (h handlers) CreatePost(ctx context.Context, req *gen.CreatePostRequest) (*gen.Post, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	p, err := converter.FromProtoPost(req.Post)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid post format")
	}

	post, err := h.postService.CreatePost(ctx, p)
	if err != nil {
		return nil, err
	}

	return converter.ToProtoPost(post), nil
}

func (h handlers) GetPost(ctx context.Context, req *gen.GetPostRequest) (*gen.Post, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	id, err := xid.FromString(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id: "+req.Id)
	}

	post, err := h.postService.GetPost(ctx, auth.UserIDFromContext(ctx), id)
	if err != nil {
		return nil, err
	}

	return converter.ToProtoPost(post), nil
}

func (h handlers) DeletePost(ctx context.Context, req *gen.DeletePostRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	id, err := xid.FromString(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id: "+req.Id)
	}

	if err := h.postService.DeletePost(ctx, auth.UserIDFromContext(ctx), id); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h handlers) BatchGetPosts(ctx context.Context, req *gen.BatchGetPostsRequest) (*gen.BatchGetPostsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ids := make([]xid.ID, len(req.Ids))
	var err error
	for i, id := range req.Ids {
		ids[i], err = xid.FromString(id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid id: "+id)
		}
	}

	posts, err := h.postService.BatchGetPosts(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &gen.BatchGetPostsResponse{
		Posts: converter.ToProtoPosts(posts),
	}, nil
}

func (h handlers) ListPosts(ctx context.Context, req *gen.ListPostsRequest) (*gen.ListPostsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	id, err := xid.FromString(req.Parent)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id: "+req.Parent)
	}

	token := xid.NilID()
	if len(req.PageToken) != 0 {
		tkn, err := xid.FromString(req.PageToken)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid page token: "+req.PageToken)
		}
		token = tkn
	}

	posts, nextID, err := h.postService.ListPosts(ctx,
		id,
		service.ListPostsPagination{
			Token: token,
			Size:  int(req.PageSize),
		},
	)
	if err != nil {
		return nil, err
	}

	var nextToken string
	if nextID.IsNil() {
		nextToken = ""
	} else {
		nextToken = nextID.String()
	}

	return &gen.ListPostsResponse{
		Posts:         converter.ToProtoPosts(posts),
		NextPageToken: nextToken,
	}, nil
}

func (h handlers) ListPostIdProjections(ctx context.Context, req *gen.ListPostIdProjectionsRequest) (*gen.ListPostIdProjectionsResponse, error) { //nolint:stylecheck,revive // method names generating by codegen
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ids := make([]xid.ID, len(req.Parents))
	var err error
	for i, id := range req.Parents {
		ids[i], err = xid.FromString(id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid id: "+id)
		}
	}

	token := xid.NilID()
	if len(req.PageToken) != 0 {
		tkn, err := xid.FromString(req.PageToken)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid page token: "+req.PageToken)
		}
		token = tkn
	}

	projections, nextID, err := h.postService.ListPostIDProjections(ctx,
		ids,
		service.ListPostIDProjectionsPagination{
			Token: token,
			Size:  int(req.PageSize),
		},
	)
	if err != nil {
		return nil, err
	}

	var nextToken string
	if nextID.IsNil() {
		nextToken = ""
	} else {
		nextToken = nextID.String()
	}

	return &gen.ListPostIdProjectionsResponse{
		PostIdProjections: converter.ToProtoPostIDProjections(projections),
		NextPageToken:     nextToken,
	}, nil
}
