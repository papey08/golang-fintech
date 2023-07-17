package httpgin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homework10/internal/model/errs"
	"homework10/internal/model/users"
	"net/http"
	"testing"
)

type userData struct {
	UserResp userResponse `json:"data"`
}

func (s *httpServerTestSuite) getUserByID(url string) (userData, int, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(s.baseURL+"/api/v1"+url), nil)
	if err != nil {
		return userData{}, 0, fmt.Errorf("unable to create request: %w", err)
	}
	var response userData
	code, err := s.getResponse(req, &response)
	if err != nil {
		return userData{}, 0, err
	}
	return response, code, nil
}

type getUserByIDMocks struct {
	id   int64
	user users.User
	err  error
}

type getUserByIDTest struct {
	givenURL           string
	expectedUser       userResponse
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestGetUserByID() {

	mocks := []getUserByIDMocks{
		{
			id: 2,
			user: users.User{
				ID:       2,
				Nickname: "papey08",
				Email:    "email02@mail.com",
			},
			err: nil,
		},
		{
			id:   50,
			user: users.User{},
			err:  errs.UserNotExist,
		},
	}

	tests := []getUserByIDTest{

		// test of successful getting user
		{
			givenURL: "/users/2",
			expectedUser: userResponse{
				ID:       2,
				Nickname: "papey08",
				Email:    "email02@mail.com",
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting code 404
		{
			givenURL:           "/users/50",
			expectedUser:       userResponse{},
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting code 400
		{
			givenURL:           "/users/abc",
			expectedUser:       userResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("GetUserByID", mock.Anything, m.id).Return(m.user, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.getUserByID(test.givenURL)
		assert.Equal(s.T(), test.expectedUser, resp.UserResp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}

	s.app.AssertCalled(s.T(), "GetUserByID", mock.Anything, mock.AnythingOfType("int64"))
}

func (s *httpServerTestSuite) createUser(body map[string]any) (userData, int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return userData{}, 0, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.baseURL+"/api/v1/users", bytes.NewReader(data))
	if err != nil {
		return userData{}, 0, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userData
	code, err := s.getResponse(req, &response)
	if err != nil {
		return userData{}, 0, err
	}

	return response, code, nil
}

type createUserMocks struct {
	nickname string
	email    string
	user     users.User
	err      error
}

type createUserTest struct {
	givenBody          map[string]any
	expectedUser       userResponse
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestCreateUser() {
	mocks := []createUserMocks{
		{
			nickname: "papey08",
			email:    "email02@mail.com",
			user: users.User{
				ID:       2,
				Nickname: "papey08",
				Email:    "email02@mail.com",
			},
			err: nil,
		},
	}

	tests := []createUserTest{
		// test of successful creating user
		{
			givenBody: map[string]any{
				"nickname": "papey08",
				"email":    "email02@mail.com",
			},
			expectedUser: userResponse{
				ID:       2,
				Nickname: "papey08",
				Email:    "email02@mail.com",
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting code 400
		{
			givenBody: map[string]any{
				"nickname": 0,
				"email":    0,
			},
			expectedUser:       userResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("CreateUser", mock.Anything, m.nickname, m.email).Return(m.user, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.createUser(test.givenBody)
		assert.Equal(s.T(), test.expectedUser, resp.UserResp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}

	s.app.AssertCalled(s.T(), "CreateUser", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"))
}

func (s *httpServerTestSuite) updateUser(url string, body map[string]any) (userData, int, error) {

	data, err := json.Marshal(body)
	if err != nil {
		return userData{}, 0, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, s.baseURL+"/api/v1"+url, bytes.NewReader(data))
	if err != nil {
		return userData{}, 0, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userData
	code, err := s.getResponse(req, &response)
	if err != nil {
		return userData{}, 0, err
	}

	return response, code, nil
}

type updateUserMocks struct {
	id       int64
	nickname string
	email    string
	user     users.User
	err      error
}

type updateUserTest struct {
	givenURL           string
	givenBody          map[string]any
	expectedUser       userResponse
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestUpdateUser() {
	mocks := []updateUserMocks{
		{
			id:       2,
			nickname: "papey08",
			email:    "new.email@mail.com",
			user: users.User{
				ID:       2,
				Nickname: "papey08",
				Email:    "new.email@mail.com",
			},
			err: nil,
		},
		{
			id:       50,
			nickname: "papey08",
			email:    "new.email@mail.com",
			user:     users.User{},
			err:      errs.UserNotExist,
		},
	}
	tests := []updateUserTest{
		// test of successful updating user
		{
			givenURL: "/users/2",
			givenBody: map[string]any{
				"nickname": "papey08",
				"email":    "new.email@mail.com",
			},
			expectedUser: userResponse{
				ID:       2,
				Nickname: "papey08",
				Email:    "new.email@mail.com",
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting code 404
		{
			givenURL: "/users/50",
			givenBody: map[string]any{
				"nickname": "papey08",
				"email":    "new.email@mail.com",
			},
			expectedUser:       userResponse{},
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting code 400
		{
			givenURL: "/users/abc",
			givenBody: map[string]any{
				"nickname": "papey08",
				"email":    "new.email@mail.com",
			},
			expectedUser:       userResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},

		// test of getting code 400
		{
			givenURL: "/users/2",
			givenBody: map[string]any{
				"nickname": 0,
				"email":    0,
			},
			expectedUser:       userResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("UpdateUser", mock.Anything, m.id, m.nickname, m.email).Return(m.user, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.updateUser(test.givenURL, test.givenBody)
		assert.Equal(s.T(), test.expectedUser, resp.UserResp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}

	s.app.AssertCalled(s.T(), "UpdateUser", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
}

func (s *httpServerTestSuite) deleteUser(url string) (int, error) {
	req, err := http.NewRequest(http.MethodDelete, s.baseURL+"/api/v1"+url, nil)
	if err != nil {
		return 0, fmt.Errorf("unable to create request: %w", err)
	}

	var response userData
	code, err := s.getResponse(req, &response)
	if err != nil {
		return 0, err
	}

	return code, nil
}

type deleteUserMocks struct {
	id  int64
	err error
}

type deleteUserTest struct {
	givenURL           string
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestDeleteUser() {
	mocks := []deleteUserMocks{
		{
			id:  2,
			err: nil,
		},
		{
			id:  50,
			err: errs.UserNotExist,
		},
	}
	tests := []deleteUserTest{
		// test of successful deleting user
		{
			givenURL:           "/users/2",
			expectedStatusCode: http.StatusOK,
		},

		// test of getting code 404
		{
			givenURL:           "/users/50",
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting code 400
		{
			givenURL:           "/users/abc",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("DeleteUser", mock.Anything, m.id).Return(m.err).Once()
	}

	for _, test := range tests {
		code, err := s.deleteUser(test.givenURL)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}
	s.app.AssertCalled(s.T(), "DeleteUser", mock.Anything, mock.AnythingOfType("int64"))
}

func BenchmarkGetUserByID(b *testing.B) {
	s := new(httpServerTestSuite)
	httpServerTestSuiteInit(s)
	for i := 0; i < b.N; i++ {
		s.app.On("GetUserByID", mock.Anything, int64(i)).Return(users.User{
			ID:       int64(i),
			Nickname: "papey08",
			Email:    "email@mail.com",
		}, nil).Once()
	}
	for i := 0; i < b.N; i++ {
		_, _, _ = s.getUserByID(fmt.Sprintf("/users/%d", i))
	}
}

func BenchmarkCreateUser(b *testing.B) {
	s := new(httpServerTestSuite)
	httpServerTestSuiteInit(s)
	for i := 0; i < b.N; i++ {
		s.app.On("CreateUser", mock.Anything, "papey08", "email@mail.com").Return(users.User{
			ID:       int64(i),
			Nickname: "papey08",
			Email:    "email@mail.com",
		}, nil).Once()
	}
	for i := 0; i < b.N; i++ {
		_, _, _ = s.createUser(map[string]any{
			"nickname": "papey08",
			"email":    "email@mail.com",
		})
	}
}

func BenchmarkDeleteUser(b *testing.B) {
	s := new(httpServerTestSuite)
	httpServerTestSuiteInit(s)
	for i := 0; i < b.N; i++ {
		s.app.On("DeleteUser", mock.Anything, int64(i)).Return(nil).Once()
	}
	for i := 0; i < b.N; i++ {
		_, _ = s.deleteUser(fmt.Sprintf("/users/%d", i))
	}
}
