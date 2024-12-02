package entity

import (
	"time"

	"github.com/rs/xid"
)

type Post struct {
	ID        xid.ID    `bson:"_id"`
	Text      string    `validate:"required,min=1,max=280"`
	AuthorID  xid.ID    `bson:"author_id" validate:"required"`
	IsDeleted bool      `bson:"deleted"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (p Post) CreatedAt() time.Time {
	return p.ID.Time()
}

type PostIDProjection struct {
	ID       xid.ID `bson:"_id"`
	AuthorID xid.ID `bson:"author_id"`
}
