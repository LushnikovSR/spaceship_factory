package inventory

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrPartNotFound = status.Error(codes.NotFound, "part is not found")
