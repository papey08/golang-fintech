package app

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	admocks "homework10/internal/app/adrepo_mocks"
	usermocks "homework10/internal/app/user_repo_mocks"
	"homework10/internal/model/ads"
	"homework10/internal/model/errs"
	"homework10/internal/model/filter"
	"homework10/internal/model/users"
	"testing"
)

type myAppTestSuite struct {
	suite.Suite
	userRepo *usermocks.UserRepository
	adRepo   *admocks.AdRepository
	service  App
}

func (s *myAppTestSuite) SetupSuite() {
	s.userRepo = new(usermocks.UserRepository)
	s.adRepo = new(admocks.AdRepository)
	s.service = NewApp(s.adRepo, s.userRepo)
}

func (s *myAppTestSuite) TearDownSuite() {
	// userRepo и adRepo не требуют TearDown, потому что основаны на map
}

type createAdMocks struct {
	userID           int64
	getUserByIDError error

	adToCreate ads.Ad
	createdAd  ads.Ad
	addAdError error
}

type createAdTest struct {
	givenTitle    string
	givenText     string
	givenUserID   int64
	expectedAd    ads.Ad
	expectedError error
}

func (s *myAppTestSuite) TestCreateAd() {
	mocks := []createAdMocks{
		{
			userID:           50,
			getUserByIDError: errs.UserNotExist,
		},
		{
			userID:           5,
			getUserByIDError: nil,
		},
		{
			userID:           5,
			getUserByIDError: nil,

			adToCreate: ads.Ad{
				ID:           0,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			createdAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			addAdError: nil,
		},
	}
	tests := []createAdTest{
		// test of creating ad by non-existing user
		{
			givenTitle:    "hello",
			givenText:     "world",
			givenUserID:   50,
			expectedAd:    ads.Ad{},
			expectedError: errs.UserNotExist,
		},

		// test of creating invalid ad
		{
			givenTitle:    "invalid ad with an empty text",
			givenText:     "",
			givenUserID:   5,
			expectedAd:    ads.Ad{},
			expectedError: errs.ValidationError,
		},

		// test of correct creating of ad
		{
			givenTitle:  "hello",
			givenText:   "world",
			givenUserID: 5,
			expectedAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedError: nil,
		},
	}

	for _, m := range mocks {
		s.userRepo.On("GetUserByID", mock.Anything, m.userID).Return(users.User{}, m.getUserByIDError).Once()
		s.adRepo.On("AddAd", mock.Anything, m.adToCreate).Return(m.createdAd, m.addAdError).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		ad, err := s.service.CreateAd(ctx, test.givenTitle, test.givenText, test.givenUserID)
		assert.Equal(s.T(), test.expectedAd, ad)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.userRepo.AssertCalled(s.T(), "GetUserByID", mock.Anything, mock.AnythingOfType("int64"))
	s.adRepo.AssertCalled(s.T(), "AddAd", mock.Anything, mock.AnythingOfType("ads.Ad"))
}

type changeAdStatusMocks struct {
	userID           int64
	getUserByIDError error

	adID           int64
	gotAd          ads.Ad
	getAdByIDError error

	idToUpdate  int64
	newAd       ads.Ad
	changedAd   ads.Ad
	changeError error
}

type changeAdStatusTest struct {
	givenAdID      int64
	givenUserID    int64
	givenPublished bool
	expectedAd     ads.Ad
	expectedError  error
}

func (s *myAppTestSuite) TestChangeAdStatus() {
	mocks := []changeAdStatusMocks{
		{
			userID:           50,
			getUserByIDError: errs.UserNotExist,
		},
		{
			userID:           5,
			getUserByIDError: nil,
			adID:             500,
			gotAd:            ads.Ad{},
			getAdByIDError:   errs.AdNotExist,
		},
		{
			userID:           10,
			getUserByIDError: nil,
			adID:             5,
			gotAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			getAdByIDError: nil,
		},
		{
			userID:           5,
			getUserByIDError: nil,
			adID:             5,
			gotAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			getAdByIDError: nil,
			idToUpdate:     5,
			newAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    true,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			changedAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    true,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			changeError: nil,
		},
	}
	tests := []changeAdStatusTest{
		// test of changing status by non-existing user
		{
			givenAdID:      50,
			givenUserID:    50,
			givenPublished: true,
			expectedAd:     ads.Ad{},
			expectedError:  errs.UserNotExist,
		},

		// test of changing status of non-existing ad
		{
			givenAdID:      500,
			givenUserID:    5,
			givenPublished: true,
			expectedAd:     ads.Ad{},
			expectedError:  errs.AdNotExist,
		},

		// test of changing status by a wrong user
		{

			givenAdID:      5,
			givenUserID:    10,
			givenPublished: true,
			expectedAd:     ads.Ad{},
			expectedError:  errs.WrongUserError,
		},

		// test of correct changing of ad status
		{
			givenAdID:      5,
			givenUserID:    5,
			givenPublished: true,
			expectedAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    true,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedError: nil,
		},
	}

	for _, m := range mocks {
		s.userRepo.On("GetUserByID", mock.Anything, m.userID).Return(users.User{}, m.getUserByIDError).Once()
		s.adRepo.On("GetAdByID", mock.Anything, m.adID).Return(m.gotAd, m.getAdByIDError).Once()
		s.adRepo.On("UpdateAdFields", mock.Anything, m.idToUpdate, m.newAd).Return(m.changedAd, m.changeError).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		changedAd, err := s.service.ChangeAdStatus(ctx, test.givenAdID, test.givenUserID, test.givenPublished)
		assert.Equal(s.T(), test.expectedAd, changedAd)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.userRepo.AssertCalled(s.T(), "GetUserByID", mock.Anything, mock.AnythingOfType("int64"))
	s.adRepo.AssertCalled(s.T(), "GetAdByID", mock.Anything, mock.AnythingOfType("int64"))
	s.adRepo.AssertCalled(s.T(), "UpdateAdFields", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("ads.Ad"))
}

type updateAdMocks struct {
	userID           int64
	getUserByIDError error

	adID           int64
	gotAd          ads.Ad
	getAdByIDError error

	idToUpdate  int64
	newAd       ads.Ad
	changedAd   ads.Ad
	changeError error
}

type updateAdTest struct {
	givenAdID     int64
	givenUserID   int64
	givenTitle    string
	givenText     string
	expectedAd    ads.Ad
	expectedError error
}

func (s *myAppTestSuite) TestUpdateAd() {
	mocks := []updateAdMocks{
		{
			userID:           50,
			getUserByIDError: errs.UserNotExist,
		},
		{
			userID:           5,
			getUserByIDError: nil,
			adID:             500,
			gotAd:            ads.Ad{},
			getAdByIDError:   errs.AdNotExist,
		},
		{
			userID:           10,
			getUserByIDError: nil,
			adID:             5,
			gotAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			getAdByIDError: nil,
		},
		{
			userID:           5,
			getUserByIDError: nil,
			adID:             5,
			gotAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			getAdByIDError: nil,
		},
		{
			userID:           5,
			getUserByIDError: nil,
			adID:             5,
			gotAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			getAdByIDError: nil,
			idToUpdate:     5,
			newAd: ads.Ad{
				ID:           5,
				Title:        "привет",
				Text:         "мир",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			changedAd: ads.Ad{
				ID:           5,
				Title:        "привет",
				Text:         "мир",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			changeError: nil,
		},
	}
	tests := []updateAdTest{
		// test of updating ad by non-existing user
		{

			givenAdID:     50,
			givenUserID:   50,
			givenTitle:    "привет",
			givenText:     "мир",
			expectedAd:    ads.Ad{},
			expectedError: errs.UserNotExist,
		},

		// test of updating of non-existing ad
		{

			givenAdID:     500,
			givenUserID:   5,
			givenTitle:    "привет",
			givenText:     "мир",
			expectedAd:    ads.Ad{},
			expectedError: errs.AdNotExist,
		},

		// test of updating ad by a wrong user
		{

			givenAdID:     5,
			givenUserID:   10,
			givenTitle:    "привет",
			givenText:     "мир",
			expectedAd:    ads.Ad{},
			expectedError: errs.WrongUserError,
		},

		// test of updating ad fields to invalid fields
		{

			givenAdID:     5,
			givenUserID:   5,
			givenTitle:    "invalid ad wit an empty text",
			givenText:     "",
			expectedAd:    ads.Ad{},
			expectedError: errs.ValidationError,
		},

		// test of correct updating ad
		{

			givenAdID:   5,
			givenUserID: 5,
			givenTitle:  "привет",
			givenText:   "мир",
			expectedAd: ads.Ad{
				ID:           5,
				Title:        "привет",
				Text:         "мир",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedError: nil,
		},
	}

	for _, m := range mocks {
		s.userRepo.On("GetUserByID", mock.Anything, m.userID).Return(users.User{}, m.getUserByIDError).Once()
		s.adRepo.On("GetAdByID", mock.Anything, m.adID).Return(m.gotAd, m.getAdByIDError).Once()
		s.adRepo.On("UpdateAdFields", mock.Anything, m.idToUpdate, m.newAd).Return(m.changedAd, m.changeError).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		changedAd, err := s.service.UpdateAd(ctx, test.givenAdID, test.givenUserID, test.givenTitle, test.givenText)
		assert.Equal(s.T(), test.expectedAd, changedAd)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.userRepo.AssertCalled(s.T(), "GetUserByID", mock.Anything, mock.AnythingOfType("int64"))
	s.adRepo.AssertCalled(s.T(), "GetAdByID", mock.Anything, mock.AnythingOfType("int64"))
	s.adRepo.AssertCalled(s.T(), "UpdateAdFields", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("ads.Ad"))
}

type getAdByIDTest struct {
	givenID       int64
	expectedAd    ads.Ad
	expectedError error
}

func (s *myAppTestSuite) TestGetAdByID() {
	tests := []getAdByIDTest{
		{
			givenID: 5,
			expectedAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedError: nil,
		},
		{
			givenID:       500,
			expectedAd:    ads.Ad{},
			expectedError: errs.AdNotExist,
		},
	}

	for _, test := range tests {
		s.adRepo.On("GetAdByID", mock.Anything, test.givenID).Return(test.expectedAd, test.expectedError).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		gotAd, err := s.service.GetAdByID(ctx, test.givenID)
		assert.Equal(s.T(), test.expectedAd, gotAd)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.adRepo.AssertCalled(s.T(), "GetAdByID", mock.Anything, mock.AnythingOfType("int64"))
}

type getAdsList struct {
	givenFilter     filter.Filter
	expectedAdsList []ads.Ad
	expectedError   error
}

func (s *myAppTestSuite) TestGetAdsList() {
	tests := []getAdsList{
		{
			givenFilter:     filter.Filter{},
			expectedAdsList: []ads.Ad{},
			expectedError:   nil,
		},
		{
			givenFilter: filter.Filter{
				PublishedBy: true,
				AuthorBy:    false,
				AuthorID:    0,
				DateBy:      false,
				Date:        ads.Date{},
			},
			expectedAdsList: []ads.Ad{
				{
					ID:           5,
					Title:        "hello",
					Text:         "world",
					AuthorID:     5,
					Published:    false,
					CreationDate: ads.CurrentDate(),
					UpdateDate:   ads.CurrentDate(),
				},
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		s.adRepo.On("GetAdsList", mock.Anything, test.givenFilter).Return(test.expectedAdsList, test.expectedError).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		list, err := s.service.GetAdsList(ctx, test.givenFilter)
		assert.Equal(s.T(), test.expectedAdsList, list)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.adRepo.AssertCalled(s.T(), "GetAdsList", mock.Anything, mock.AnythingOfType("filter.Filter"))
}

type searchAds struct {
	givenPattern  string
	expectedAds   []ads.Ad
	expectedError error
}

func (s *myAppTestSuite) TestSearchAds() {
	tests := []searchAds{
		{
			givenPattern: "hello",
			expectedAds: []ads.Ad{
				{
					ID:           5,
					Title:        "hello",
					Text:         "world",
					AuthorID:     5,
					Published:    false,
					CreationDate: ads.CurrentDate(),
					UpdateDate:   ads.CurrentDate(),
				},
			},
			expectedError: nil,
		},
		{
			givenPattern:  "привет",
			expectedAds:   []ads.Ad{},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		s.adRepo.On("SearchAds", mock.Anything, test.givenPattern).Return(test.expectedAds, test.expectedError).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		list, err := s.service.SearchAds(ctx, test.givenPattern)
		assert.Equal(s.T(), test.expectedAds, list)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.adRepo.AssertCalled(s.T(), "SearchAds", mock.Anything, mock.AnythingOfType("string"))
}

type deleteAdMocks struct {
	userID           int64
	gotUser          users.User
	getUserByIDError error

	adID           int64
	gotAd          ads.Ad
	getAdByIDError error

	idToDelete  int64
	deleteError error
}

type deleteAdTest struct {
	givenUserID   int64
	givenAdID     int64
	expectedError error
}

func (s *myAppTestSuite) TestDeleteAd() {
	mocks := []deleteAdMocks{
		{
			userID:           50,
			gotUser:          users.User{},
			getUserByIDError: errs.UserNotExist,
		},
		{
			userID: 5,
			gotUser: users.User{
				ID:       5,
				Nickname: "papey08",
				Email:    "email05@mail.com",
			},
			getUserByIDError: nil,

			adID:           500,
			gotAd:          ads.Ad{},
			getAdByIDError: errs.AdNotExist,
		},
		{
			userID: 10,
			gotUser: users.User{
				ID:       10,
				Nickname: "wrong_user",
				Email:    "wrong@mail.com",
			},
			getUserByIDError: nil,

			adID: 5,
			gotAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			getAdByIDError: nil,
		},
		{
			userID: 5,
			gotUser: users.User{
				ID:       5,
				Nickname: "papey08",
				Email:    "email05@mail.com",
			},
			getUserByIDError: nil,

			adID: 5,
			gotAd: ads.Ad{
				ID:           5,
				Title:        "hello",
				Text:         "world",
				AuthorID:     5,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			getAdByIDError: nil,

			idToDelete:  5,
			deleteError: nil,
		},
	}
	tests := []deleteAdTest{
		// test of deleting ad by non-existing user
		{

			givenUserID:   50,
			givenAdID:     5,
			expectedError: errs.UserNotExist,
		},

		//test of deleting non-existing ad
		{

			givenUserID:   5,
			givenAdID:     500,
			expectedError: errs.AdNotExist,
		},

		// test of deleting ad by wrong user
		{

			givenUserID:   10,
			givenAdID:     5,
			expectedError: errs.WrongUserError,
		},

		// test of correct ad delete
		{

			givenUserID:   5,
			givenAdID:     5,
			expectedError: nil,
		},
	}

	for _, m := range mocks {
		s.userRepo.On("GetUserByID", mock.Anything, m.userID).Return(m.gotUser, m.getUserByIDError).Once()
		s.adRepo.On("GetAdByID", mock.Anything, m.adID).Return(m.gotAd, m.getAdByIDError).Once()
		s.adRepo.On("DeleteAd", mock.Anything, m.idToDelete).Return(m.deleteError).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		err := s.service.DeleteAd(ctx, test.givenUserID, test.givenAdID)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.userRepo.AssertCalled(s.T(), "GetUserByID", mock.Anything, mock.AnythingOfType("int64"))
	s.adRepo.AssertCalled(s.T(), "GetAdByID", mock.Anything, mock.AnythingOfType("int64"))
	s.adRepo.AssertCalled(s.T(), "DeleteAd", mock.Anything, mock.AnythingOfType("int64"))
}

type getUserByIDTest struct {
	givenID       int64
	expectedUser  users.User
	expectedError error
}

func (s *myAppTestSuite) TestGetUserByID() {
	tests := []getUserByIDTest{
		{
			givenID: 0,
			expectedUser: users.User{
				ID:       0,
				Nickname: "papey08",
				Email:    "email00@mail.com",
			},
			expectedError: nil,
		},
		{
			givenID: 1,
			expectedUser: users.User{
				ID:       1,
				Nickname: "one_more_user",
				Email:    "email01@mail.com",
			},
			expectedError: nil,
		},
		{
			givenID:       2,
			expectedUser:  users.User{},
			expectedError: errs.UserNotExist,
		},
	}

	for _, test := range tests {
		s.userRepo.On("GetUserByID", mock.Anything, test.givenID).Return(test.expectedUser, test.expectedError).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		gotUser, err := s.service.GetUserByID(ctx, test.givenID)
		assert.Equal(s.T(), test.expectedUser, gotUser)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.userRepo.AssertCalled(s.T(), "GetUserByID", mock.Anything, mock.AnythingOfType("int64"))
}

type createUserTest struct {
	givenUser    users.User
	expectedUser users.User
}

func (s *myAppTestSuite) TestCreateUser() {
	tests := []createUserTest{
		{
			givenUser: users.User{
				ID:       0,
				Nickname: "papey08",
				Email:    "email00@mail.com",
			},
			expectedUser: users.User{
				ID:       0,
				Nickname: "papey08",
				Email:    "email00@mail.com",
			},
		},
		{
			givenUser: users.User{
				ID:       0,
				Nickname: "one_more_user",
				Email:    "email01@mail.com",
			},
			expectedUser: users.User{
				ID:       1,
				Nickname: "one_more_user",
				Email:    "email01@mail.com",
			},
		},
	}

	for _, test := range tests {
		s.userRepo.On("AddUser", mock.Anything, test.givenUser).Return(test.expectedUser, nil).Once()
	}

	ctx := context.Background()

	createdUser, err := s.service.CreateUser(ctx, "papey08", "email00@mail.com")
	assert.Equal(s.T(), tests[0].expectedUser, createdUser)
	assert.Nil(s.T(), err)

	createdUser, err = s.service.CreateUser(ctx, "one_more_user", "email01@mail.com")
	assert.Equal(s.T(), tests[1].expectedUser, createdUser)
	assert.Nil(s.T(), err)

	s.userRepo.AssertCalled(s.T(), "AddUser", mock.Anything, mock.AnythingOfType("users.User"))
}

type updateUserTest struct {
	givenID       int64
	givenUser     users.User
	expectedUser  users.User
	expectedError error
}

func (s *myAppTestSuite) TestUpdateUser() {
	tests := []updateUserTest{
		{
			givenID: 5,
			givenUser: users.User{
				ID:       0,
				Nickname: "papey08",
				Email:    "changed_email@mail.com",
			},
			expectedUser: users.User{
				ID:       5,
				Nickname: "papey08",
				Email:    "changed_email@mail.com",
			},
			expectedError: nil,
		},
		{
			givenID: 50,
			givenUser: users.User{
				ID:       0,
				Nickname: "changed_nickname08",
				Email:    "changed_email@mail.com",
			},
			expectedUser:  users.User{},
			expectedError: errs.UserNotExist,
		},
	}

	for _, test := range tests {
		s.userRepo.On("UpdateUserFields", mock.Anything, test.givenID, test.givenUser).Return(test.expectedUser, test.expectedError).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		gotUser, err := s.service.UpdateUser(ctx, test.givenID, test.givenUser.Nickname, test.givenUser.Email)
		assert.Equal(s.T(), test.expectedUser, gotUser)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.userRepo.AssertCalled(s.T(), "UpdateUserFields", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("users.User"))
}

type deleteUserTest struct {
	givenID       int64
	expectedError error
}

func (s *myAppTestSuite) TestDeleteUser() {
	tests := []deleteUserTest{
		{
			givenID:       5,
			expectedError: nil,
		},
		{
			givenID:       50,
			expectedError: errs.UserNotExist,
		},
	}

	for _, test := range tests {
		s.userRepo.On("DeleteUser", mock.Anything, test.givenID).Return(test.expectedError).Once()
	}

	ctx := context.Background()

	for _, test := range tests {
		err := s.service.DeleteUser(ctx, test.givenID)
		assert.Equal(s.T(), test.expectedError, err)
	}
	s.userRepo.AssertCalled(s.T(), "DeleteUser", mock.Anything, mock.AnythingOfType("int64"))
}

func TestMyAppTestSuite(t *testing.T) {
	suite.Run(t, new(myAppTestSuite))
}
