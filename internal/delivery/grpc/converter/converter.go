package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rs/xid"

	gen "github.com/Karzoug/meower-post-service/internal/delivery/grpc/gen/post/v1"
	"github.com/Karzoug/meower-post-service/internal/post/entity"
)

func ToProtoPost(post entity.Post) *gen.Post {
	return &gen.Post{
		Id:          post.ID.String(),
		Text:        post.Text,
		AuthorId:    post.AuthorID.String(),
		UpdatedTime: timestamppb.New(post.UpdatedAt),
		Deleted:     post.IsDeleted,
	}
}

func ToProtoPosts(posts []entity.Post) []*gen.Post {
	pps := make([]*gen.Post, len(posts))
	for i := range posts {
		pps[i] = ToProtoPost(posts[i])
	}

	return pps
}

func ToProtoPostIDProjections(posts []entity.PostIDProjection) []*gen.PostIdProjection {
	pps := make([]*gen.PostIdProjection, len(posts))
	for i := range posts {
		pps[i] = &gen.PostIdProjection{
			Id:       posts[i].ID.String(),
			AuthorId: posts[i].AuthorID.String(),
		}
	}

	return pps
}

func FromProtoPost(post *gen.Post) (entity.Post, error) {
	id, err := xid.FromString(post.Id)
	if err != nil {
		id = xid.NilID()
	}
	authorID, err := xid.FromString(post.AuthorId)
	if err != nil {
		return entity.Post{}, err
	}
	return entity.Post{
		ID:       id,
		Text:     post.Text,
		AuthorID: authorID,
	}, nil
}
