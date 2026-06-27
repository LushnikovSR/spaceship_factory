package payment

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrMissingArgument = status.Error(codes.InvalidArgument, "argument must be specified")
