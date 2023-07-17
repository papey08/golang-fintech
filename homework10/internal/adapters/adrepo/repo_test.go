package adrepo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework10/internal/model/ads"
	"homework10/internal/model/errs"
	"homework10/internal/model/filter"
	"sync"
	"testing"
)

type adRepoTestSuite struct {
	suite.Suite
	adRepo AdRepository
}

func (s *adRepoTestSuite) SetupSuite() {
	// init adRepo
	s.adRepo = AdRepository{
		data:   make(map[int64]ads.Ad),
		freeID: 0,
		mu:     sync.RWMutex{},
	}

	ctx := context.Background()

	// fill repo with fake ads
	_, _ = s.adRepo.AddAd(ctx, ads.Ad{
		Title:        "ad00",
		Text:         "default ad",
		AuthorID:     0,
		Published:    false,
		CreationDate: ads.Date{Day: 30, Month: 4, Year: 2023},
		UpdateDate:   ads.Date{Day: 30, Month: 4, Year: 2023},
	})
	_, _ = s.adRepo.AddAd(ctx, ads.Ad{
		Title:        "ad01",
		Text:         "published ad",
		AuthorID:     1,
		Published:    true,
		CreationDate: ads.Date{Day: 30, Month: 4, Year: 2023},
		UpdateDate:   ads.Date{Day: 30, Month: 4, Year: 2023},
	})
	_, _ = s.adRepo.AddAd(ctx, ads.Ad{
		Title:        "ad02",
		Text:         "ad with other creation date",
		AuthorID:     2,
		Published:    false,
		CreationDate: ads.Date{Day: 29, Month: 4, Year: 2023},
		UpdateDate:   ads.Date{Day: 29, Month: 4, Year: 2023},
	})
	_, _ = s.adRepo.AddAd(ctx, ads.Ad{
		Title:        "diff_ad03",
		Text:         "ad with other title",
		AuthorID:     3,
		Published:    false,
		CreationDate: ads.Date{Day: 30, Month: 4, Year: 2023},
		UpdateDate:   ads.Date{Day: 30, Month: 4, Year: 2023},
	})
	_, _ = s.adRepo.AddAd(ctx, ads.Ad{
		Title:        "diff_ad04",
		Text:         "ad with other title",
		AuthorID:     4,
		Published:    false,
		CreationDate: ads.Date{Day: 30, Month: 4, Year: 2023},
		UpdateDate:   ads.Date{Day: 30, Month: 4, Year: 2023},
	})
}

func (s *adRepoTestSuite) TearDownSuite() {

}

type getAdByIDTest struct {
	givenContext  context.Context
	givenID       int64
	expectedAd    ads.Ad
	expectedError error
}

func (s *adRepoTestSuite) TestGetAdByID() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []getAdByIDTest{
		// test of getting ad with done context
		{
			givenContext:  doneCtx,
			givenID:       1,
			expectedAd:    ads.Ad{},
			expectedError: errs.AdRepositoryError,
		},

		// test of getting non-existing ad
		{
			givenContext:  ctx,
			givenID:       50,
			expectedAd:    ads.Ad{},
			expectedError: errs.AdNotExist,
		},

		// test of successful getting ad
		{
			givenContext: ctx,
			givenID:      2,
			expectedAd: ads.Ad{
				ID:           2,
				Title:        "ad02",
				Text:         "ad with other creation date",
				AuthorID:     2,
				Published:    false,
				CreationDate: ads.Date{Day: 29, Month: 4, Year: 2023},
				UpdateDate:   ads.Date{Day: 29, Month: 4, Year: 2023},
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		ad, err := s.adRepo.GetAdByID(test.givenContext, test.givenID)
		assert.Equal(s.T(), test.expectedAd, ad)
		assert.Equal(s.T(), test.expectedError, err)
	}
}

type addAdTest struct {
	givenContext  context.Context
	givenAd       ads.Ad
	expectedAd    ads.Ad
	expectedError error
}

func (s *adRepoTestSuite) TestAddAd() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []addAdTest{
		// test of adding ad with done context
		{
			givenContext: doneCtx,
			givenAd: ads.Ad{
				Title:        "notAddedAd",
				Text:         "text of not added ad",
				AuthorID:     0,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedAd:    ads.Ad{},
			expectedError: errs.AdRepositoryError,
		},

		// test of successful adding ad
		{
			givenContext: ctx,
			givenAd: ads.Ad{
				Title:        "hello",
				Text:         "world",
				AuthorID:     2,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedAd: ads.Ad{
				ID:           5,
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

	for _, test := range tests {
		ad, err := s.adRepo.AddAd(test.givenContext, test.givenAd)
		assert.Equal(s.T(), test.expectedAd, ad)
		assert.Equal(s.T(), test.expectedError, err)
		if err != nil {
			_ = s.adRepo.DeleteAd(test.givenContext, ad.ID)
		}
	}
}

type updateAdFieldsTest struct {
	givenContext  context.Context
	givenID       int64
	givenAd       ads.Ad
	expectedAd    ads.Ad
	expectedError error
}

func (s *adRepoTestSuite) TestUpdateAdFields() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []updateAdFieldsTest{
		// test of updating ad with done context
		{
			givenContext: doneCtx,
			givenID:      2,
			givenAd: ads.Ad{
				Title:        "notUpdatedAd02",
				Text:         "ad with other creation date",
				AuthorID:     2,
				Published:    false,
				CreationDate: ads.Date{Day: 29, Month: 4, Year: 2023},
				UpdateDate:   ads.CurrentDate(),
			},
			expectedAd:    ads.Ad{},
			expectedError: errs.AdRepositoryError,
		},

		// test of updating non-existing ad
		{
			givenContext: ctx,
			givenID:      50,
			givenAd: ads.Ad{
				Title:        "nonExistingAd",
				Text:         "text of non-existing ad",
				AuthorID:     0,
				Published:    false,
				CreationDate: ads.Date{Day: 30, Month: 4, Year: 2023},
				UpdateDate:   ads.CurrentDate(),
			},
			expectedAd:    ads.Ad{},
			expectedError: errs.AdNotExist,
		},

		// test of successful updating ad
		{
			givenContext: ctx,
			givenID:      4,
			givenAd: ads.Ad{
				Title:        "diff_ad04",
				Text:         "updated ad with other title",
				AuthorID:     4,
				Published:    false,
				CreationDate: ads.Date{Day: 30, Month: 4, Year: 2023},
				UpdateDate:   ads.CurrentDate(),
			},
			expectedAd: ads.Ad{
				ID:           4,
				Title:        "diff_ad04",
				Text:         "updated ad with other title",
				AuthorID:     4,
				Published:    false,
				CreationDate: ads.Date{Day: 30, Month: 4, Year: 2023},
				UpdateDate:   ads.CurrentDate(),
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		ad, err := s.adRepo.UpdateAdFields(test.givenContext, test.givenID, test.givenAd)
		assert.Equal(s.T(), test.expectedAd, ad)
		assert.Equal(s.T(), test.expectedError, err)
	}
}

type getAdsListTest struct {
	givenContext    context.Context
	givenFilter     filter.Filter
	expectedListLen int
	expectedError   error
}

func (s *adRepoTestSuite) TestGetAdsList() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []getAdsListTest{
		// test of getting ads list with done context
		{
			givenContext:    doneCtx,
			givenFilter:     filter.Filter{},
			expectedListLen: 0,
			expectedError:   errs.AdRepositoryError,
		},

		// test of getting ads list with default filter
		{
			givenContext:    ctx,
			givenFilter:     filter.Filter{},
			expectedListLen: 1,
			expectedError:   nil,
		},

		// test of getting ads list with ads created by author 1
		{
			givenContext: ctx,
			givenFilter: filter.Filter{
				AuthorBy: true,
				AuthorID: 1,
			},
			expectedListLen: 1,
			expectedError:   nil,
		},

		// test of getting ads list with ads created on 29.04.2023
		{
			givenContext: ctx,
			givenFilter: filter.Filter{
				DateBy: true,
				Date:   ads.Date{Day: 29, Month: 04, Year: 2023},
			},
		},

		// test of getting ads list with both author and date filters
		{
			givenContext: ctx,
			givenFilter: filter.Filter{
				PublishedBy: true,
				AuthorBy:    true,
				AuthorID:    1,
				DateBy:      true,
				Date:        ads.Date{Day: 29, Month: 4, Year: 2023},
			},
			expectedListLen: 0,
			expectedError:   nil,
		},
	}

	for _, test := range tests {
		list, err := s.adRepo.GetAdsList(test.givenContext, test.givenFilter)
		assert.Len(s.T(), list, test.expectedListLen)
		assert.Equal(s.T(), test.expectedError, err)
	}
}

type searchAdsTest struct {
	givenContext    context.Context
	givenPattern    string
	expectedListLen int
	expectedError   error
}

func (s *adRepoTestSuite) TestSearchAds() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []searchAdsTest{
		// test of searching ads with done context
		{
			givenContext:    doneCtx,
			givenPattern:    "ad",
			expectedListLen: 0,
			expectedError:   errs.AdRepositoryError,
		},

		// test of searching non-existing ad
		{
			givenContext:    ctx,
			givenPattern:    "non-existing ad",
			expectedListLen: 0,
			expectedError:   nil,
		},

		// test of searching ad by full match
		{
			givenContext:    ctx,
			givenPattern:    "diff_ad03",
			expectedListLen: 1,
			expectedError:   nil,
		},

		// test of searching ad by prefix match
		{
			givenContext:    ctx,
			givenPattern:    "diff_ad",
			expectedListLen: 2,
			expectedError:   nil,
		},
	}

	for _, test := range tests {
		list, err := s.adRepo.SearchAds(test.givenContext, test.givenPattern)
		assert.Len(s.T(), list, test.expectedListLen)
		assert.Equal(s.T(), test.expectedError, err)
	}
}

type deleteAdTest struct {
	givenContext  context.Context
	givenID       int64
	expectedError error
}

func (s *adRepoTestSuite) TestDeleteAd() {
	ctx := context.Background()
	doneCtx, cancel := context.WithCancel(ctx)
	cancel()

	tests := []deleteAdTest{
		// test of deleting ad with done context
		{
			givenContext:  doneCtx,
			givenID:       1,
			expectedError: errs.AdRepositoryError,
		},

		// test of deleting non-existing ad
		{
			givenContext:  ctx,
			givenID:       50,
			expectedError: errs.AdNotExist,
		},

		// test of successful deleting ad
		{
			givenContext:  ctx,
			givenID:       0,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		err := s.adRepo.DeleteAd(test.givenContext, test.givenID)
		assert.Equal(s.T(), test.expectedError, err)
		if err == nil {
			_, getErr := s.adRepo.GetAdByID(test.givenContext, test.givenID)
			assert.Equal(s.T(), errs.AdNotExist, getErr)
		}
	}
}

func TestAdRepoTestSuite(t *testing.T) {
	suite.Run(t, new(adRepoTestSuite))
}
