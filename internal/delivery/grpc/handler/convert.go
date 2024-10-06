package handler

import (
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Karzoug/meower-post-service/pkg/ucerr"
	"github.com/Karzoug/meower-post-service/pkg/ucerr/codes"
)

func toGRPCError(err ucerr.Error) error {
	switch err.Code() {
	case codes.NotFound:
		return status.Error(grpcCodes.NotFound, err.Error())
	case codes.InvalidArgument:
		return status.Error(grpcCodes.InvalidArgument, err.Error())
	case codes.AlreadyExists:
		return status.Error(grpcCodes.AlreadyExists, err.Error())
	case codes.Unauthenticated:
		return status.Error(grpcCodes.Unauthenticated, err.Error())
	default:
		return status.Error(grpcCodes.Internal, err.Error())
	}
}
