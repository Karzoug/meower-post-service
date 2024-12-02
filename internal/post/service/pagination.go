package service

import "github.com/rs/xid"

type ListPostsPagination struct {
	Token xid.ID
	Size  int `validate:"omitempty,min=1,max=100"`
}

type ListPostIDProjectionsPagination struct {
	Token xid.ID
	Size  int `validate:"omitempty,min=1,max=1000"`
}
