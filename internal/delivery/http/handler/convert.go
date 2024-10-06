package handler

import (
	"net/http"

	"github.com/Karzoug/meower-post-service/internal/delivery/http/gen"
	"github.com/Karzoug/meower-post-service/pkg/ucerr"
	"github.com/Karzoug/meower-post-service/pkg/ucerr/codes"
)

type errorJSONRespose struct {
	Body       gen.ErrorResponse
	StatusCode int
}

func toErrorJSONRespose(err ucerr.Error) errorJSONRespose {
	return errorJSONRespose{
		StatusCode: toHTTPError(err),
		Body: gen.ErrorResponse{
			Error: err.Error(),
		},
	}
}

func toHTTPError(err ucerr.Error) int {
	switch err.Code() {
	case codes.NotFound:
		return http.StatusNotFound
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
