package entity

import (
	"time"

	"github.com/rs/xid"

	"github.com/Karzoug/meower-post-service/pkg/validator"
)

type Post struct {
	ID        xid.ID
	Text      string `validate:"required,min=1,max=280"`
	AuthorID  xid.ID
	IsDeleted bool
	UpdatedAt time.Time
}

func (p Post) CreatedAt() time.Time {
	return p.ID.Time()
}

func NewPost(authorID xid.ID, text string) (Post, error) {
	p := Post{
		Text:      text,
		AuthorID:  authorID,
		UpdatedAt: time.Now(),
	}

	return p, validator.Struct(p)
}

type PostIDProjection struct {
	ID       xid.ID
	AuthorID xid.ID
}
