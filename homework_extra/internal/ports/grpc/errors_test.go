package grpc

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework_extra/internal/model/errs"
	"testing"
)

type errorToGRPCErrorTest struct {
	name              string
	givenError        error
	expectedErrorCode codes.Code
}

func TestErrorToGRPCError(t *testing.T) {
	tests := []errorToGRPCErrorTest{
		{
			name:              "test of mapping errs.AdNotExist",
			givenError:        errs.AdNotExist,
			expectedErrorCode: codes.NotFound,
		},
		{
			name:              "test of mapping errs.UserNotExist",
			givenError:        errs.UserNotExist,
			expectedErrorCode: codes.NotFound,
		},
		{
			name:              "test of mapping errs.WrongUserError",
			givenError:        errs.WrongUserError,
			expectedErrorCode: codes.PermissionDenied,
		},
		{
			name:              "test of mapping errs.ValidationError",
			givenError:        errs.ValidationError,
			expectedErrorCode: codes.FailedPrecondition,
		},
		{
			name:              "test of mapping errs.AdRepositoryError",
			givenError:        errs.AdRepositoryError,
			expectedErrorCode: codes.ResourceExhausted,
		},
		{
			name:              "test of mapping errs.UserRepositoryError",
			givenError:        errs.UserRepositoryError,
			expectedErrorCode: codes.ResourceExhausted,
		},
		{
			name:              "test of mapping unexpected error",
			givenError:        errors.New("unexpected error"),
			expectedErrorCode: codes.Unknown,
		},
		{
			name:              "test of mapping nil error",
			givenError:        nil,
			expectedErrorCode: codes.OK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ErrorToGRPCError(test.givenError)
			assert.Equal(t, test.expectedErrorCode, status.Code(err))
		})
	}
}
