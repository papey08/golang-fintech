package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework9/internal/model/errs"
)

// ErrorToGRPCError returns error with its code from google.golang.org/grpc/codes
func ErrorToGRPCError(err error) error {
	var c codes.Code
	switch err {
	case nil:
		return nil
	case errs.AdNotExist:
		c = codes.NotFound
	case errs.UserNotExist:
		c = codes.NotFound
	case errs.WrongUserError:
		c = codes.PermissionDenied
	case errs.ValidationError:
		c = codes.FailedPrecondition
	case errs.UserAlreadyExists:
		c = codes.AlreadyExists
	case errs.AdRepositoryError:
		c = codes.ResourceExhausted
	case errs.UserRepositoryError:
		c = codes.ResourceExhausted
	default:
		c = codes.Unknown
	}
	return status.Errorf(c, err.Error())
}
