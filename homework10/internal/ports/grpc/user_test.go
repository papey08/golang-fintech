package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homework10/internal/model/errs"
	"homework10/internal/model/users"
	"homework10/internal/ports/grpc/pb"
	"testing"
)

func userResponseToUser(u *pb.UserResponse) users.User {
	if u == nil {
		return users.User{}
	} else {
		return users.User{
			ID:       u.Id,
			Nickname: u.Nickname,
			Email:    u.Email,
		}
	}
}

type getUserByIDMocks struct {
	ctx context.Context
	id  int64
	usr users.User
	err error
}

type getUserByIDTest struct {
	givenContext  context.Context
	givenReq      *pb.GetUserByIDRequest
	expectedUser  users.User
	expectedError error
}

func (s *grpcServerTestSuite) TestGetUserByID() {
	ctx := context.Background()

	tests := []getUserByIDTest{
		// test of getting nil UserResponse and error
		{
			givenContext:  ctx,
			givenReq:      &pb.GetUserByIDRequest{Id: 50},
			expectedUser:  users.User{},
			expectedError: ErrorToGRPCError(errs.UserNotExist),
		},

		// test of getting UserResponse and nil error
		{
			givenContext: ctx,
			givenReq:     &pb.GetUserByIDRequest{Id: 2},
			expectedUser: users.User{
				ID:       2,
				Nickname: "papey08",
				Email:    "email02@mail.com",
			},
			expectedError: nil,
		},
	}
	mocks := []getUserByIDMocks{
		{
			ctx: ctx,
			id:  50,
			usr: users.User{},
			err: errs.UserNotExist,
		},
		{
			ctx: ctx,
			id:  2,
			usr: users.User{
				ID:       2,
				Nickname: "papey08",
				Email:    "email02@mail.com",
			},
			err: nil,
		},
	}

	for _, m := range mocks {
		s.app.On("GetUserByID", m.ctx, m.id).Return(m.usr, m.err).Once()
	}

	for _, test := range tests {
		usrResp, err := s.server.GetUserByID(test.givenContext, test.givenReq)
		usr := userResponseToUser(usrResp)
		assert.Equal(s.T(), test.expectedUser, usr)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "GetUserByID", mock.Anything, mock.AnythingOfType("int64"))
}

type createUserMocks struct {
	ctx      context.Context
	nickname string
	email    string
	usr      users.User
	err      error
}

type createUserTest struct {
	givenContext  context.Context
	givenReq      *pb.CreateUserRequest
	expectedUser  users.User
	expectedError error
}

func (s *grpcServerTestSuite) TestCreateUser() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []createUserTest{
		// test of getting nil UserResponse and error
		{
			givenContext: doneCtx,
			givenReq: &pb.CreateUserRequest{
				Nickname: "papey08",
				Email:    "email02@mail.com",
			},
			expectedUser:  users.User{},
			expectedError: ErrorToGRPCError(errs.UserRepositoryError),
		},

		// test of getting UserResponse and nil error
		{
			givenContext: ctx,
			givenReq: &pb.CreateUserRequest{
				Nickname: "papey08",
				Email:    "email02@mail.com",
			},
			expectedUser: users.User{
				ID:       2,
				Nickname: "papey08",
				Email:    "email02@mail.com",
			},
			expectedError: nil,
		},
	}
	mocks := []createUserMocks{
		{
			ctx:      doneCtx,
			nickname: "papey08",
			email:    "email02@mail.com",
			usr:      users.User{},
			err:      errs.UserRepositoryError,
		},
		{
			ctx:      ctx,
			nickname: "papey08",
			email:    "email02@mail.com",
			usr: users.User{
				ID:       2,
				Nickname: "papey08",
				Email:    "email02@mail.com",
			},
			err: nil,
		},
	}

	for _, m := range mocks {
		s.app.On("CreateUser", m.ctx, m.nickname, m.email).Return(m.usr, m.err).Once()
	}

	for _, test := range tests {
		usrResp, err := s.server.CreateUser(test.givenContext, test.givenReq)
		usr := userResponseToUser(usrResp)
		assert.Equal(s.T(), test.expectedUser, usr)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "CreateUser", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"))
}

type updateUserMocks struct {
	ctx      context.Context
	id       int64
	nickname string
	email    string
	usr      users.User
	err      error
}

type updateUserTest struct {
	givenContext  context.Context
	givenReq      *pb.UpdateUserRequest
	expectedUser  users.User
	expectedError error
}

func (s *grpcServerTestSuite) TestUpdateUser() {
	ctx := context.Background()

	tests := []updateUserTest{
		// test of getting nil UserResponse and error
		{
			givenContext: ctx,
			givenReq: &pb.UpdateUserRequest{
				Id:       50,
				Nickname: "new nickname",
				Email:    "new.email@mail.com",
			},
			expectedUser:  users.User{},
			expectedError: ErrorToGRPCError(errs.UserNotExist),
		},

		// test of getting UserResponse and nil error
		{
			givenContext: ctx,
			givenReq: &pb.UpdateUserRequest{
				Id:       2,
				Nickname: "new nickname",
				Email:    "new.email@mail.com",
			},
			expectedUser: users.User{
				ID:       2,
				Nickname: "new nickname",
				Email:    "new.email@mail.com",
			},
			expectedError: nil,
		},
	}
	mocks := []updateUserMocks{
		{
			ctx:      ctx,
			id:       50,
			nickname: "new nickname",
			email:    "new.email@mail.com",
			usr:      users.User{},
			err:      errs.UserNotExist,
		},
		{
			ctx:      ctx,
			id:       2,
			nickname: "new nickname",
			email:    "new.email@mail.com",
			usr: users.User{
				ID:       2,
				Nickname: "new nickname",
				Email:    "new.email@mail.com",
			},
			err: nil,
		},
	}

	for _, m := range mocks {
		s.app.On("UpdateUser", m.ctx, m.id, m.nickname, m.email).Return(m.usr, m.err).Once()
	}

	for _, test := range tests {
		usrResp, err := s.server.UpdateUser(test.givenContext, test.givenReq)
		usr := userResponseToUser(usrResp)
		assert.Equal(s.T(), test.expectedUser, usr)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "UpdateUser", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
}

type deleteUserMocks struct {
	ctx context.Context
	id  int64
	err error
}

type deleteUserTest struct {
	givenContext  context.Context
	givenReq      *pb.DeleteUserRequest
	expectedError error
}

func (s *grpcServerTestSuite) TestDeleteUserTest() {
	ctx := context.Background()

	tests := []deleteUserTest{
		// test of getting error
		{
			givenContext: ctx,
			givenReq: &pb.DeleteUserRequest{
				Id: 50,
			},
			expectedError: ErrorToGRPCError(errs.UserNotExist),
		},

		// test of getting nil error
		{
			givenContext: ctx,
			givenReq: &pb.DeleteUserRequest{
				Id: 2,
			},
			expectedError: nil,
		},
	}
	mocks := []deleteUserMocks{
		{
			ctx: ctx,
			id:  50,
			err: errs.UserNotExist,
		},
		{
			ctx: ctx,
			id:  2,
			err: nil,
		},
	}

	for _, m := range mocks {
		s.app.On("DeleteUser", m.ctx, m.id).Return(m.err).Once()
	}

	for _, test := range tests {
		_, err := s.server.DeleteUser(test.givenContext, test.givenReq)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "DeleteUser", mock.Anything, mock.AnythingOfType("int64"))
}

func BenchmarkGetUserByID(b *testing.B) {
	s := new(grpcServerTestSuite)
	grpcServerTestSuiteInit(s)
	for i := 0; i < b.N; i++ {
		s.app.On("GetUserByID", mock.Anything, int64(i)).Return(users.User{
			ID:       int64(i),
			Nickname: "papey08",
			Email:    "email@mail.com",
		}, nil).Once()
	}
	for i := 0; i < b.N; i++ {
		_, _ = s.client.GetUserByID(context.Background(), &pb.GetUserByIDRequest{Id: int64(i)})
	}
}

func BenchmarkCreateUser(b *testing.B) {
	s := new(grpcServerTestSuite)
	grpcServerTestSuiteInit(s)
	for i := 0; i < b.N; i++ {
		s.app.On("CreateUser", mock.Anything, "papey08", "email@mail.com").Return(users.User{
			ID:       int64(i),
			Nickname: "papey08",
			Email:    "email@mail.com",
		}, nil).Once()
	}
	for i := 0; i < b.N; i++ {
		_, _ = s.client.CreateUser(context.Background(), &pb.CreateUserRequest{
			Nickname: "papey08",
			Email:    "email@mail.com",
		})
	}
}

func BenchmarkDeleteUser(b *testing.B) {
	s := new(grpcServerTestSuite)
	grpcServerTestSuiteInit(s)
	for i := 0; i < b.N; i++ {
		s.app.On("DeleteUser", mock.Anything, int64(i)).Return(nil).Once()
	}
	for i := 0; i < b.N; i++ {
		_, _ = s.client.DeleteUser(context.Background(), &pb.DeleteUserRequest{Id: int64(i)})
	}
}
