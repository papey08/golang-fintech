package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homework10/internal/model/ads"
	"homework10/internal/model/errs"
	"homework10/internal/model/filter"
	"homework10/internal/ports/grpc/pb"
	"time"
)

func adResponseToAd(ad *pb.AdResponse) ads.Ad {
	if ad == nil {
		return ads.Ad{}
	} else {
		return ads.Ad{
			ID:        ad.Id,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorId,
			Published: ad.Published,
			CreationDate: ads.Date{
				Day:   int(ad.CreationDate.Day),
				Month: time.Month(ad.CreationDate.Month),
				Year:  int(ad.CreationDate.Year),
			},
			UpdateDate: ads.Date{
				Day:   int(ad.UpdateDate.Day),
				Month: time.Month(ad.UpdateDate.Month),
				Year:  int(ad.UpdateDate.Year),
			},
		}
	}
}

func listAdsResponseToAdsList(resp *pb.ListAdResponse) []ads.Ad {
	if resp == nil {
		return make([]ads.Ad, 0)
	}
	list := make([]ads.Ad, 0, len(resp.List))
	for _, a := range resp.List {
		list = append(list, adResponseToAd(a))
	}
	return list
}

type createAdMocks struct {
	ctx      context.Context
	title    string
	text     string
	authorID int64
	ad       ads.Ad
	err      error
}

type createAdTest struct {
	givenContext  context.Context
	givenReq      *pb.CreateAdRequest
	expectedAd    ads.Ad
	expectedError error
}

func (s *grpcServerTestSuite) TestCreateAd() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []createAdTest{
		// test of getting nil AdResponse and error
		{
			givenContext: doneCtx,
			givenReq: &pb.CreateAdRequest{
				Title:    "hello",
				Text:     "world",
				AuthorId: 2,
			},
			expectedAd:    ads.Ad{},
			expectedError: ErrorToGRPCError(errs.AdRepositoryError),
		},

		// test of getting AdResponse
		{
			givenContext: ctx,
			givenReq: &pb.CreateAdRequest{
				Title:    "hello",
				Text:     "world",
				AuthorId: 2,
			},
			expectedAd: ads.Ad{
				ID:           2,
				Title:        "hello",
				Text:         "world",
				AuthorID:     2,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedError: nil,
		},
	}
	mocks := []createAdMocks{
		{
			ctx:      doneCtx,
			title:    "hello",
			text:     "world",
			authorID: 2,
			ad:       ads.Ad{},
			err:      errs.AdRepositoryError,
		},
		{
			ctx:      ctx,
			title:    "hello",
			text:     "world",
			authorID: 2,
			ad: ads.Ad{
				ID:           2,
				Title:        "hello",
				Text:         "world",
				AuthorID:     2,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			err: nil,
		},
	}

	for _, m := range mocks {
		s.app.On("CreateAd", m.ctx, m.title, m.text, m.authorID).Return(m.ad, m.err).Once()
	}

	for _, test := range tests {
		adResp, err := s.server.CreateAd(test.givenContext, test.givenReq)
		ad := adResponseToAd(adResp)
		assert.Equal(s.T(), test.expectedAd, ad)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "CreateAd", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int64"))
}

type changeAdStatusMocks struct {
	ctx       context.Context
	id        int64
	authorID  int64
	published bool
	ad        ads.Ad
	err       error
}

type changeAdStatusTest struct {
	givenContext  context.Context
	givenReq      *pb.ChangeAdStatusRequest
	expectedAd    ads.Ad
	expectedError error
}

func (s *grpcServerTestSuite) TestChangeAdStatus() {
	ctx := context.Background()

	tests := []changeAdStatusTest{
		// test of getting nil AdResponse and error
		{
			givenContext: ctx,
			givenReq: &pb.ChangeAdStatusRequest{
				Id:        50,
				AuthorId:  2,
				Published: true,
			},
			expectedAd:    ads.Ad{},
			expectedError: ErrorToGRPCError(errs.AdNotExist),
		},

		// test of getting AdResponse
		{
			givenContext: ctx,
			givenReq: &pb.ChangeAdStatusRequest{
				Id:        2,
				AuthorId:  2,
				Published: true,
			},
			expectedAd: ads.Ad{
				ID:           2,
				Title:        "hello",
				Text:         "world",
				AuthorID:     2,
				Published:    true,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedError: nil,
		},
	}
	mocks := []changeAdStatusMocks{
		{
			ctx:       ctx,
			id:        50,
			authorID:  2,
			published: true,
			ad:        ads.Ad{},
			err:       errs.AdNotExist,
		},
		{
			ctx:       ctx,
			id:        2,
			authorID:  2,
			published: true,
			ad: ads.Ad{
				ID:           2,
				Title:        "hello",
				Text:         "world",
				AuthorID:     2,
				Published:    true,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			err: nil,
		},
	}

	for _, m := range mocks {
		s.app.On("ChangeAdStatus", m.ctx, m.id, m.authorID, m.published).Return(m.ad, m.err).Once()
	}

	for _, test := range tests {
		adResp, err := s.server.ChangeAdStatus(test.givenContext, test.givenReq)
		ad := adResponseToAd(adResp)
		assert.Equal(s.T(), test.expectedAd, ad)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "ChangeAdStatus", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("bool"))
}

type updateAdMocks struct {
	ctx      context.Context
	id       int64
	authorID int64
	title    string
	text     string
	ad       ads.Ad
	err      error
}

type updateAdTest struct {
	givenContext  context.Context
	givenReq      *pb.UpdateAdRequest
	expectedAd    ads.Ad
	expectedError error
}

func (s *grpcServerTestSuite) TestUpdateAd() {
	ctx := context.Background()

	tests := []updateAdTest{
		// test of getting nil AdResponse and error
		{
			givenContext: ctx,
			givenReq: &pb.UpdateAdRequest{
				Id:       50,
				Title:    "привет",
				Text:     "мир",
				AuthorId: 2,
			},
			expectedAd:    ads.Ad{},
			expectedError: ErrorToGRPCError(errs.AdNotExist),
		},

		// test of getting AdResponse
		{
			givenContext: ctx,
			givenReq: &pb.UpdateAdRequest{
				Id:       2,
				Title:    "привет",
				Text:     "мир",
				AuthorId: 2,
			},
			expectedAd: ads.Ad{
				ID:           2,
				Title:        "привет",
				Text:         "мир",
				AuthorID:     2,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedError: nil,
		},
	}
	mocks := []updateAdMocks{
		{
			ctx:      ctx,
			id:       50,
			authorID: 2,
			title:    "привет",
			text:     "мир",
			ad:       ads.Ad{},
			err:      errs.AdNotExist,
		},
		{
			ctx:      ctx,
			id:       2,
			authorID: 2,
			title:    "привет",
			text:     "мир",
			ad: ads.Ad{
				ID:           2,
				Title:        "привет",
				Text:         "мир",
				AuthorID:     2,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			err: nil,
		},
	}

	for _, m := range mocks {
		s.app.On("UpdateAd", m.ctx, m.id, m.authorID, m.title, m.text).Return(m.ad, m.err).Once()
	}

	for _, test := range tests {
		adResp, err := s.server.UpdateAd(test.givenContext, test.givenReq)
		ad := adResponseToAd(adResp)
		assert.Equal(s.T(), test.expectedAd, ad)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "UpdateAd", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
}

type getAdByIDMocks struct {
	ctx context.Context
	id  int64
	ad  ads.Ad
	err error
}

type getAdByIDTest struct {
	givenContext  context.Context
	givenReq      *pb.GetAdByIDRequest
	expectedAd    ads.Ad
	expectedError error
}

func (s *grpcServerTestSuite) TestGetAdByID() {
	ctx := context.Background()

	tests := []getAdByIDTest{
		// test of getting nil AdResponse and error
		{
			givenContext: ctx,
			givenReq: &pb.GetAdByIDRequest{
				Id: 50,
			},
			expectedAd:    ads.Ad{},
			expectedError: ErrorToGRPCError(errs.AdNotExist),
		},

		//test of getting AdResponse
		{
			givenContext: ctx,
			givenReq: &pb.GetAdByIDRequest{
				Id: 2,
			},
			expectedAd: ads.Ad{
				ID:           2,
				Title:        "hello",
				Text:         "world",
				AuthorID:     2,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedError: nil,
		},
	}
	mocks := []getAdByIDMocks{
		{
			ctx: ctx,
			id:  50,
			ad:  ads.Ad{},
			err: errs.AdNotExist,
		},
		{
			ctx: ctx,
			id:  2,
			ad: ads.Ad{
				ID:           2,
				Title:        "hello",
				Text:         "world",
				AuthorID:     2,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			err: nil,
		},
	}

	for _, m := range mocks {
		s.app.On("GetAdByID", m.ctx, m.id).Return(m.ad, m.err).Once()
	}

	for _, test := range tests {
		adResp, err := s.server.GetAdByID(test.givenContext, test.givenReq)
		ad := adResponseToAd(adResp)
		assert.Equal(s.T(), test.expectedAd, ad)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "GetAdByID", mock.Anything, mock.AnythingOfType("int64"))
}

type getAdsListMocks struct {
	ctx    context.Context
	f      filter.Filter
	listAd []ads.Ad
	err    error
}

type getAdsListTest struct {
	givenContext  context.Context
	givenReq      *pb.FilterRequest
	expectedList  []ads.Ad
	expectedError error
}

func (s *grpcServerTestSuite) TestGetAdsList() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []getAdsListTest{
		// test of getting nil ListAdResponse and error
		{
			givenContext:  doneCtx,
			givenReq:      &pb.FilterRequest{},
			expectedList:  []ads.Ad{},
			expectedError: ErrorToGRPCError(errs.AdRepositoryError),
		},

		// test of getting ListAdResponse
		{
			givenContext: ctx,
			givenReq:     &pb.FilterRequest{},
			expectedList: []ads.Ad{
				{
					ID:           2,
					Title:        "hello",
					Text:         "world",
					AuthorID:     2,
					Published:    true,
					CreationDate: ads.CurrentDate(),
					UpdateDate:   ads.CurrentDate(),
				},
			},
			expectedError: nil,
		},
	}
	mocks := []getAdsListMocks{
		{
			ctx:    doneCtx,
			f:      filter.Filter{},
			listAd: []ads.Ad{},
			err:    errs.AdRepositoryError,
		},
		{
			ctx: ctx,
			f:   filter.Filter{},
			listAd: []ads.Ad{
				{
					ID:           2,
					Title:        "hello",
					Text:         "world",
					AuthorID:     2,
					Published:    true,
					CreationDate: ads.CurrentDate(),
					UpdateDate:   ads.CurrentDate(),
				},
			},
			err: nil,
		},
	}

	for _, m := range mocks {
		s.app.On("GetAdsList", m.ctx, m.f).Return(m.listAd, m.err).Once()
	}

	for _, test := range tests {
		listAdResp, err := s.server.GetAdsList(test.givenContext, test.givenReq)
		assert.Len(s.T(), listAdsResponseToAdsList(listAdResp), len(test.expectedList))
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "GetAdsList", mock.Anything, mock.AnythingOfType("filter.Filter"))
}

type searchAdsMocks struct {
	ctx     context.Context
	pattern string
	listAd  []ads.Ad
	err     error
}

type searchAdsTest struct {
	givenContext  context.Context
	givenReq      *pb.SearchAdsRequest
	expectedList  []ads.Ad
	expectedError error
}

func (s *grpcServerTestSuite) TestSearchAds() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []searchAdsTest{
		// test of getting nil ListAdResponse and error
		{
			givenContext: doneCtx,
			givenReq: &pb.SearchAdsRequest{
				Pattern: "hello",
			},
			expectedList:  []ads.Ad{},
			expectedError: ErrorToGRPCError(errs.AdRepositoryError),
		},

		// test of getting ListAdResponse
		{
			givenContext: ctx,
			givenReq: &pb.SearchAdsRequest{
				Pattern: "hello",
			},
			expectedList: []ads.Ad{
				{
					ID:           2,
					Title:        "hello",
					Text:         "world",
					AuthorID:     2,
					Published:    true,
					CreationDate: ads.CurrentDate(),
					UpdateDate:   ads.CurrentDate(),
				},
			},
			expectedError: nil,
		},
	}
	mocks := []searchAdsMocks{
		{
			ctx:     doneCtx,
			pattern: "hello",
			listAd:  []ads.Ad{},
			err:     errs.AdRepositoryError,
		},
		{
			ctx:     ctx,
			pattern: "hello",
			listAd: []ads.Ad{
				{
					ID:           2,
					Title:        "hello",
					Text:         "world",
					AuthorID:     2,
					Published:    true,
					CreationDate: ads.CurrentDate(),
					UpdateDate:   ads.CurrentDate(),
				},
			},
			err: nil,
		},
	}

	for _, m := range mocks {
		s.app.On("SearchAds", m.ctx, m.pattern).Return(m.listAd, m.err).Once()
	}

	for _, test := range tests {
		listAdResp, err := s.server.SearchAds(test.givenContext, test.givenReq)
		assert.Len(s.T(), listAdsResponseToAdsList(listAdResp), len(test.expectedList))
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "SearchAds", mock.Anything, mock.AnythingOfType("string"))
}

type deleteAdMocks struct {
	ctx      context.Context
	authorID int64
	id       int64
	err      error
}

type deleteAdTest struct {
	givenContext  context.Context
	givenReq      *pb.DeleteAdRequest
	expectedError error
}

func (s *grpcServerTestSuite) TestDeleteAd() {
	ctx := context.Background()

	tests := []deleteAdTest{
		// test of getting error
		{
			givenContext: ctx,
			givenReq: &pb.DeleteAdRequest{
				Id:       50,
				AuthorId: 2,
			},
			expectedError: ErrorToGRPCError(errs.AdNotExist),
		},

		// test of getting nil error
		{
			givenContext: ctx,
			givenReq: &pb.DeleteAdRequest{
				Id:       2,
				AuthorId: 2,
			},
			expectedError: nil,
		},
	}
	mocks := []deleteAdMocks{
		{
			ctx:      ctx,
			authorID: 2,
			id:       50,
			err:      errs.AdNotExist,
		},
		{
			ctx:      ctx,
			authorID: 2,
			id:       2,
			err:      nil,
		},
	}

	for _, m := range mocks {
		s.app.On("DeleteAd", m.ctx, m.authorID, m.id).Return(m.err).Once()
	}

	for _, test := range tests {
		_, err := s.server.DeleteAd(test.givenContext, test.givenReq)
		assert.Equal(s.T(), test.expectedError, err)
	}

	s.app.AssertCalled(s.T(), "DeleteAd", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64"))
}
