package user_repo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework10/internal/model/errs"
	"homework10/internal/model/users"
	"sync"
	"testing"
)

type userRepoTestSuite struct {
	suite.Suite
	userRepo UserRepository
}

func (s *userRepoTestSuite) SetupSuite() {
	// init userRepo
	s.userRepo = UserRepository{
		data:   make(map[int64]users.User),
		freeID: 0,
		mu:     sync.RWMutex{},
	}

	ctx := context.Background()

	// fill repo with fake users
	_, _ = s.userRepo.AddUser(ctx, users.User{
		Nickname: "user00",
		Email:    "email00@mail.com",
	})
	_, _ = s.userRepo.AddUser(ctx, users.User{
		Nickname: "user01",
		Email:    "email01@mail.com",
	})
	_, _ = s.userRepo.AddUser(ctx, users.User{
		Nickname: "user02",
		Email:    "email02@mail.com",
	})
	_, _ = s.userRepo.AddUser(ctx, users.User{
		Nickname: "user03",
		Email:    "email03@mail.com",
	})
}

func (s *userRepoTestSuite) TearDownSuite() {

}

type getUserByIDTest struct {
	givenContext  context.Context
	givenID       int64
	expectedUser  users.User
	expectedError error
}

func (s *userRepoTestSuite) TestGetUserByID() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []getUserByIDTest{
		// test of getting user with done context
		{
			givenContext:  doneCtx,
			givenID:       1,
			expectedUser:  users.User{},
			expectedError: errs.UserRepositoryError,
		},

		// test of getting non-existing user
		{
			givenContext:  ctx,
			givenID:       50,
			expectedUser:  users.User{},
			expectedError: errs.UserNotExist,
		},

		// test of successful getting user
		{
			givenContext: ctx,
			givenID:      2,
			expectedUser: users.User{
				ID:       2,
				Nickname: "user02",
				Email:    "email02@mail.com",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		usr, err := s.userRepo.GetUserByID(test.givenContext, test.givenID)
		assert.Equal(s.T(), test.expectedUser, usr)
		assert.Equal(s.T(), test.expectedError, err)
	}
}

type addUserTest struct {
	givenContext  context.Context
	givenUser     users.User
	expectedUser  users.User
	expectedError error
}

func (s *userRepoTestSuite) TestAddUser() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []addUserTest{
		// test of adding user with done context
		{
			givenContext: doneCtx,
			givenUser: users.User{
				ID:       0,
				Nickname: "notAddedUser",
				Email:    "notAddedUser@mail.com",
			},
			expectedUser:  users.User{},
			expectedError: errs.UserRepositoryError,
		},

		// test of successful adding of user
		{
			givenContext: ctx,
			givenUser: users.User{
				ID:       0,
				Nickname: "user04",
				Email:    "email04@mail.com",
			},
			expectedUser: users.User{
				ID:       4,
				Nickname: "user04",
				Email:    "email04@mail.com",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		usr, err := s.userRepo.AddUser(test.givenContext, test.givenUser)
		assert.Equal(s.T(), test.expectedUser, usr)
		assert.Equal(s.T(), test.expectedError, err)
	}
}

type updateUserFieldsTest struct {
	givenContext  context.Context
	givenID       int64
	givenUser     users.User
	expectedUser  users.User
	expectedError error
}

func (s *userRepoTestSuite) TestUpdateUserFields() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []updateUserFieldsTest{
		// test of updating user with done context
		{
			givenContext: doneCtx,
			givenID:      2,
			givenUser: users.User{
				Nickname: "updatedUser02",
				Email:    "email02@mail.com",
			},
			expectedUser:  users.User{},
			expectedError: errs.UserRepositoryError,
		},

		// test of updating non-existing user
		{
			givenContext: ctx,
			givenID:      50,
			givenUser: users.User{
				Nickname: "nonExistingUser",
				Email:    "nonExistingUser@mail.com",
			},
			expectedUser:  users.User{},
			expectedError: errs.UserNotExist,
		},

		// test of successful updating user
		{
			givenContext: ctx,
			givenID:      2,
			givenUser: users.User{
				Nickname: "updatedUser02",
				Email:    "email02@mail.com",
			},
			expectedUser: users.User{
				ID:       2,
				Nickname: "updatedUser02",
				Email:    "email02@mail.com",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		usr, err := s.userRepo.UpdateUserFields(test.givenContext, test.givenID, test.givenUser)
		assert.Equal(s.T(), test.expectedUser, usr)
		assert.Equal(s.T(), test.expectedError, err)
	}
}

type deleteUserTest struct {
	givenContext  context.Context
	givenID       int64
	expectedError error
}

func (s *userRepoTestSuite) TestDeleteUser() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []deleteUserTest{
		// test of deleting user with done context
		{
			givenContext:  doneCtx,
			givenID:       2,
			expectedError: errs.UserRepositoryError,
		},

		// test of deleting non-existing user
		{
			givenContext:  ctx,
			givenID:       50,
			expectedError: errs.UserNotExist,
		},

		// test of successful deleting user
		{
			givenContext:  ctx,
			givenID:       1,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		err := s.userRepo.DeleteUser(test.givenContext, test.givenID)
		assert.Equal(s.T(), test.expectedError, err)

		if err == nil {
			_, getErr := s.userRepo.GetUserByID(test.givenContext, test.givenID)
			assert.Equal(s.T(), errs.UserNotExist, getErr)
		}
	}
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(userRepoTestSuite))
}
