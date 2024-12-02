package service

import "github.com/rs/xid"

type ListPostsPagination struct {
	Token xid.ID
	Size  int
}

type ListPostIDProjectionsPagination struct {
	Token xid.ID
	Size  int
}
