package httpgin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homework10/internal/model/ads"
	"homework10/internal/model/errs"
	"homework10/internal/model/filter"
	"net/http"
)

type adData struct {
	AdResp adResponse `json:"data"`
}

type adsData struct {
	AdsResp adsResponse `json:"data"`
}

func (s *httpServerTestSuite) createAd(body map[string]any) (adData, int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return adData{}, 0, fmt.Errorf("unable to marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, s.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adData{}, 0, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	var response adData
	code, err := s.getResponse(req, &response)
	if err != nil {
		return adData{}, 0, err
	}
	return response, code, nil
}

type createAdMocks struct {
	title  string
	text   string
	userID int64
	ad     ads.Ad
	err    error
}

type createAdTest struct {
	givenBody          map[string]any
	expectedAd         adResponse
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestCreateAd() {
	mocks := []createAdMocks{
		{
			title:  "hello",
			text:   "world",
			userID: 2,
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
		{
			title:  "invalid ad with an empty text",
			text:   "",
			userID: 2,
			ad:     ads.Ad{},
			err:    errs.ValidationError,
		},
		{
			title:  "ad with non-existing userID",
			text:   "world",
			userID: 50,
			ad:     ads.Ad{},
			err:    errs.UserNotExist,
		},
	}
	tests := []createAdTest{
		// test of successful creating ad
		{
			givenBody: map[string]any{
				"title":   "hello",
				"text":    "world",
				"user_id": 2,
			},
			expectedAd: adResponse{
				ID:           2,
				Title:        "hello",
				Text:         "world",
				AuthorID:     2,
				Published:    false,
				CreationDate: ads.CurrentDate(),
				UpdateDate:   ads.CurrentDate(),
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting code 400
		{
			givenBody: map[string]any{
				"title":   "invalid ad with an empty text",
				"text":    "",
				"user_id": 2,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},

		// test of getting code 404
		{
			givenBody: map[string]any{
				"title":   "ad with non-existing userID",
				"text":    "world",
				"user_id": 50,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting code 400
		{
			givenBody: map[string]any{
				"title":   "invalid ad request",
				"text":    0,
				"user_id": "hello world",
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("CreateAd", mock.Anything, m.title, m.text, m.userID).Return(m.ad, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.createAd(test.givenBody)
		assert.Equal(s.T(), test.expectedAd, resp.AdResp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}

	s.app.AssertCalled(s.T(), "CreateAd", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int64"))
}

func (s *httpServerTestSuite) changeAdStatus(url string, body map[string]any) (adData, int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return adData{}, 0, fmt.Errorf("unable to marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPut, s.baseURL+"/api/v1"+url, bytes.NewReader(data))
	if err != nil {
		return adData{}, 0, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	var response adData
	code, err := s.getResponse(req, &response)
	if err != nil {
		return adData{}, 0, err
	}
	return response, code, nil
}

type changeAdStatusMocks struct {
	adID      int64
	userID    int64
	published bool
	ad        ads.Ad
	err       error
}

type changeAdStatusTest struct {
	givenURL           string
	givenBody          map[string]any
	expectedAd         adResponse
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestChangeAdStatus() {
	mocks := []changeAdStatusMocks{
		{
			adID:      2,
			userID:    2,
			published: true,
			ad: ads.Ad{
				ID:        2,
				Title:     "hello",
				Text:      "world",
				AuthorID:  2,
				Published: true,
				CreationDate: ads.Date{
					Day:   29,
					Month: 4,
					Year:  2022,
				},
				UpdateDate: ads.CurrentDate(),
			},
			err: nil,
		},
		{
			adID:      50,
			userID:    2,
			published: true,
			ad:        ads.Ad{},
			err:       errs.AdNotExist,
		},
		{
			adID:      2,
			userID:    50,
			published: true,
			ad:        ads.Ad{},
			err:       errs.UserNotExist,
		},
		{
			adID:      2,
			userID:    5,
			published: true,
			ad:        ads.Ad{},
			err:       errs.WrongUserError,
		},
	}

	tests := []changeAdStatusTest{
		// test of correct changing ad status
		{
			givenURL: "/ads/2/status",
			givenBody: map[string]any{
				"published": true,
				"user_id":   2,
			},
			expectedAd: adResponse{
				ID:        2,
				Title:     "hello",
				Text:      "world",
				AuthorID:  2,
				Published: true,
				CreationDate: ads.Date{
					Day:   29,
					Month: 4,
					Year:  2022,
				},
				UpdateDate: ads.CurrentDate(),
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting code 404
		{
			givenURL: "/ads/50/status",
			givenBody: map[string]any{
				"published": true,
				"user_id":   2,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting code 404
		{
			givenURL: "/ads/2/status",
			givenBody: map[string]any{
				"published": true,
				"user_id":   50,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting code 403
		{
			givenURL: "/ads/2/status",
			givenBody: map[string]any{
				"published": true,
				"user_id":   5,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusForbidden,
		},

		// test of getting code 400
		{
			givenURL: "/ads/abc/status",
			givenBody: map[string]any{
				"published": true,
				"user_id":   2,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},

		// test of getting code 400
		{
			givenURL: "/ads/2/status",
			givenBody: map[string]any{
				"published": 2,
				"user_id":   true,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("ChangeAdStatus", mock.Anything, m.adID, m.userID, m.published).Return(m.ad, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.changeAdStatus(test.givenURL, test.givenBody)
		assert.Equal(s.T(), test.expectedAd, resp.AdResp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}

	s.app.AssertCalled(s.T(), "ChangeAdStatus", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("bool"))
}

func (s *httpServerTestSuite) updateAd(url string, body map[string]any) (adData, int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return adData{}, 0, fmt.Errorf("unable to marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPut, s.baseURL+"/api/v1"+url, bytes.NewReader(data))
	if err != nil {
		return adData{}, 0, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	var response adData
	code, err := s.getResponse(req, &response)
	if err != nil {
		return adData{}, 0, err
	}
	return response, code, nil
}

type updateAdMocks struct {
	adID   int64
	userID int64
	title  string
	text   string
	ad     ads.Ad
	err    error
}

type updateAdTest struct {
	givenURL           string
	givenBody          map[string]any
	expectedAd         adResponse
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestUpdateAd() {
	mocks := []updateAdMocks{
		{
			adID:   2,
			userID: 2,
			title:  "привет",
			text:   "мир",
			ad: ads.Ad{
				ID:        2,
				Title:     "привет",
				Text:      "мир",
				AuthorID:  2,
				Published: false,
				CreationDate: ads.Date{
					Day:   29,
					Month: 4,
					Year:  2022,
				},
				UpdateDate: ads.CurrentDate(),
			},
			err: nil,
		},
		{
			adID:   2,
			userID: 5,
			title:  "привет",
			text:   "мир",
			ad:     ads.Ad{},
			err:    errs.WrongUserError,
		},
		{
			adID:   2,
			userID: 2,
			title:  "invalid ad with an empty text",
			text:   "",
			ad:     ads.Ad{},
			err:    errs.ValidationError,
		},
		{
			adID:   50,
			userID: 2,
			title:  "привет",
			text:   "мир",
			ad:     ads.Ad{},
			err:    errs.AdNotExist,
		},
		{
			adID:   2,
			userID: 50,
			title:  "привет",
			text:   "мир",
			ad:     ads.Ad{},
			err:    errs.UserNotExist,
		},
	}
	tests := []updateAdTest{
		// test of successful updating ad
		{
			givenURL: "/ads/2",
			givenBody: map[string]any{
				"title":   "привет",
				"text":    "мир",
				"user_id": 2,
			},
			expectedAd: adResponse{

				ID:        2,
				Title:     "привет",
				Text:      "мир",
				AuthorID:  2,
				Published: false,
				CreationDate: ads.Date{
					Day:   29,
					Month: 4,
					Year:  2022,
				},
				UpdateDate: ads.CurrentDate(),
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting code 403
		{
			givenURL: "/ads/2",
			givenBody: map[string]any{
				"title":   "привет",
				"text":    "мир",
				"user_id": 5,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusForbidden,
		},

		// test of getting code 400
		{
			givenURL: "/ads/2",
			givenBody: map[string]any{
				"title":   "invalid ad with an empty text",
				"text":    "",
				"user_id": 2,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},

		// test of getting code 404
		{
			givenURL: "/ads/50",
			givenBody: map[string]any{
				"title":   "привет",
				"text":    "мир",
				"user_id": 2,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting code 404
		{
			givenURL: "/ads/2",
			givenBody: map[string]any{
				"title":   "привет",
				"text":    "мир",
				"user_id": 50,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting code 400
		{
			givenURL: "/ads/abc",
			givenBody: map[string]any{
				"title":   "привет",
				"text":    "мир",
				"user_id": 2,
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},

		// test of getting code 400
		{
			givenURL: "/ads/2",
			givenBody: map[string]any{
				"title":   "привет",
				"text":    2,
				"user_id": "мир",
			},
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("UpdateAd", mock.Anything, m.adID, m.userID, m.title, m.text).Return(m.ad, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.updateAd(test.givenURL, test.givenBody)
		assert.Equal(s.T(), test.expectedAd, resp.AdResp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}

	s.app.AssertCalled(s.T(), "UpdateAd", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64"), mock.AnythingOfType("string"), mock.AnythingOfType("string"))
}

func (s *httpServerTestSuite) getAdsList(body map[string]any) (adsData, int, error) {
	var data []byte
	if body == nil {
		data = []byte{}
	} else {
		var err error
		data, err = json.Marshal(body)
		if err != nil {
			return adsData{}, 0, fmt.Errorf("unable to marshal: %w", err)
		}
	}
	req, err := http.NewRequest(http.MethodGet, s.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adsData{}, 0, fmt.Errorf("unable to create request: %w", err)
	}
	var response adsData
	code, err := s.getResponse(req, &response)
	if err != nil {
		return adsData{}, 0, err
	}
	return response, code, nil
}

type getAdsListMocks struct {
	f    filter.Filter
	list []ads.Ad
	err  error
}

type getAdsListTest struct {
	givenBody          map[string]any
	expectedList       adsResponse
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestGetAdsList() {
	mocks := []getAdsListMocks{
		{
			f: filter.Filter{},
			list: []ads.Ad{
				{
					ID:        2,
					Title:     "hello",
					Text:      "world",
					AuthorID:  2,
					Published: true,
					CreationDate: ads.Date{
						Day:   29,
						Month: 4,
						Year:  2022,
					},
					UpdateDate: ads.Date{
						Day:   12,
						Month: 10,
						Year:  2022,
					},
				},
				{
					ID:        5,
					Title:     "привет",
					Text:      "мир",
					AuthorID:  5,
					Published: true,
					CreationDate: ads.Date{
						Day:   29,
						Month: 4,
						Year:  2022,
					},
					UpdateDate: ads.Date{
						Day:   12,
						Month: 10,
						Year:  2022,
					},
				},
			},
			err: nil,
		},
		{
			f: filter.Filter{
				PublishedBy: false,
				AuthorBy:    true,
				AuthorID:    2,
				DateBy:      false,
				Date:        ads.Date{},
			},
			list: []ads.Ad{
				{
					ID:        2,
					Title:     "hello",
					Text:      "world",
					AuthorID:  2,
					Published: true,
					CreationDate: ads.Date{
						Day:   29,
						Month: 4,
						Year:  2022,
					},
					UpdateDate: ads.Date{
						Day:   12,
						Month: 10,
						Year:  2022,
					},
				},
			},
			err: nil,
		},
	}
	tests := []getAdsListTest{
		// test of getting ads list with default filter
		{
			givenBody: nil,
			expectedList: []adResponse{
				{
					ID:        2,
					Title:     "hello",
					Text:      "world",
					AuthorID:  2,
					Published: true,
					CreationDate: ads.Date{
						Day:   29,
						Month: 4,
						Year:  2022,
					},
					UpdateDate: ads.Date{
						Day:   12,
						Month: 10,
						Year:  2022,
					},
				},
				{
					ID:        5,
					Title:     "привет",
					Text:      "мир",
					AuthorID:  5,
					Published: true,
					CreationDate: ads.Date{
						Day:   29,
						Month: 4,
						Year:  2022,
					},
					UpdateDate: ads.Date{
						Day:   12,
						Month: 10,
						Year:  2022,
					},
				},
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting ads list with custom filter
		{
			givenBody: map[string]any{
				"published_by": false,
				"author_by":    true,
				"author_id":    2,
				"date_by":      false,
			},
			expectedList: []adResponse{
				{
					ID:        2,
					Title:     "hello",
					Text:      "world",
					AuthorID:  2,
					Published: true,
					CreationDate: ads.Date{
						Day:   29,
						Month: 4,
						Year:  2022,
					},
					UpdateDate: ads.Date{
						Day:   12,
						Month: 10,
						Year:  2022,
					},
				},
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting status code 400
		{
			givenBody: map[string]any{
				"published_by": 2,
				"author_by":    2,
				"author_id":    true,
				"date_by":      2,
			},
			expectedList:       nil,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("GetAdsList", mock.Anything, m.f).Return(m.list, m.err).Once()
	}

	for _, test := range tests {
		list, code, err := s.getAdsList(test.givenBody)
		assert.Equal(s.T(), test.expectedList, list.AdsResp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}

	s.app.AssertCalled(s.T(), "GetAdsList", mock.Anything, mock.AnythingOfType("filter.Filter"))
}

func (s *httpServerTestSuite) searchAds(body map[string]any) (adsData, int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return adsData{}, 0, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, s.baseURL+"/api/v1/ads/search", bytes.NewReader(data))
	if err != nil {
		return adsData{}, 0, fmt.Errorf("unable to create request: %w", err)
	}

	var response adsData
	code, err := s.getResponse(req, &response)
	if err != nil {
		return adsData{}, 0, err
	}

	return response, code, nil
}

type searchAdsMocks struct {
	pattern string
	list    []ads.Ad
	err     error
}

type searchAdsTest struct {
	givenBody          map[string]any
	expectedList       adsResponse
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestSearchAds() {
	mocks := []searchAdsMocks{
		{
			pattern: "hello",
			list: []ads.Ad{
				{
					ID:        2,
					Title:     "hello",
					Text:      "world",
					AuthorID:  2,
					Published: true,
					CreationDate: ads.Date{
						Day:   29,
						Month: 4,
						Year:  2022,
					},
					UpdateDate: ads.Date{
						Day:   12,
						Month: 10,
						Year:  2022,
					},
				},
			},
			err: nil,
		},
	}
	tests := []searchAdsTest{
		// test of successful searching ads
		{
			givenBody: map[string]any{
				"pattern": "hello",
			},
			expectedList: adsResponse{
				{
					ID:        2,
					Title:     "hello",
					Text:      "world",
					AuthorID:  2,
					Published: true,
					CreationDate: ads.Date{
						Day:   29,
						Month: 4,
						Year:  2022,
					},
					UpdateDate: ads.Date{
						Day:   12,
						Month: 10,
						Year:  2022,
					},
				},
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting status code 400
		{
			givenBody: map[string]any{
				"pattern": 400,
			},
			expectedList:       nil,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("SearchAds", mock.Anything, m.pattern).Return(m.list, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.searchAds(test.givenBody)
		assert.Equal(s.T(), test.expectedList, resp.AdsResp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}

	s.app.AssertCalled(s.T(), "SearchAds", mock.Anything, mock.AnythingOfType("string"))
}

func (s *httpServerTestSuite) getAdByID(url string) (adData, int, error) {
	req, err := http.NewRequest(http.MethodGet, s.baseURL+"/api/v1"+url, nil)
	if err != nil {
		return adData{}, 0, fmt.Errorf("unable to create request: %w", err)
	}
	var response adData
	code, err := s.getResponse(req, &response)
	if err != nil {
		return adData{}, 0, err
	}
	return response, code, nil
}

type getAdByIDMocks struct {
	adID int64
	ad   ads.Ad
	err  error
}

type getAdByIDTest struct {
	givenURL           string
	expectedAd         adResponse
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestGetAdByID() {
	mocks := []getAdByIDMocks{
		{
			adID: 2,
			ad: ads.Ad{
				ID:        2,
				Title:     "hello",
				Text:      "world",
				AuthorID:  2,
				Published: true,
				CreationDate: ads.Date{
					Day:   29,
					Month: 4,
					Year:  2022,
				},
				UpdateDate: ads.Date{
					Day:   12,
					Month: 10,
					Year:  2022,
				},
			},
			err: nil,
		},
		{
			adID: 50,
			ad:   ads.Ad{},
			err:  errs.AdNotExist,
		},
	}
	tests := []getAdByIDTest{
		// test of successful getting ad
		{
			givenURL: "/ads/2",
			expectedAd: adResponse{
				ID:        2,
				Title:     "hello",
				Text:      "world",
				AuthorID:  2,
				Published: true,
				CreationDate: ads.Date{
					Day:   29,
					Month: 4,
					Year:  2022,
				},
				UpdateDate: ads.Date{
					Day:   12,
					Month: 10,
					Year:  2022,
				},
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting status code 404
		{
			givenURL:           "/ads/50",
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting status code 400
		{
			givenURL:           "/ads/abc",
			expectedAd:         adResponse{},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("GetAdByID", mock.Anything, m.adID).Return(m.ad, m.err).Once()
	}

	for _, test := range tests {
		resp, code, err := s.getAdByID(test.givenURL)
		assert.Equal(s.T(), test.expectedAd, resp.AdResp)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}

	s.app.AssertCalled(s.T(), "GetAdByID", mock.Anything, mock.AnythingOfType("int64"))
}

func (s *httpServerTestSuite) deleteAd(url string, body map[string]any) (int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("unable to marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodDelete, s.baseURL+"/api/v1"+url, bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("unable to create request: %w", err)
	}
	var response adData
	return s.getResponse(req, &response)
}

type deleteAdMocks struct {
	userID int64
	adID   int64
	err    error
}

type deleteAdTest struct {
	givenURL           string
	givenBody          map[string]any
	expectedStatusCode int
}

func (s *httpServerTestSuite) TestDeleteAd() {
	mocks := []deleteAdMocks{
		{
			userID: 2,
			adID:   2,
			err:    nil,
		},
		{
			userID: 5,
			adID:   2,
			err:    errs.WrongUserError,
		},
		{
			userID: 2,
			adID:   50,
			err:    errs.AdNotExist,
		},
		{
			userID: 50,
			adID:   2,
			err:    errs.UserNotExist,
		},
	}
	tests := []deleteAdTest{
		// test of successful deleting ad
		{
			givenURL: "/ads/2",
			givenBody: map[string]any{
				"user_id": 2,
			},
			expectedStatusCode: http.StatusOK,
		},

		// test of getting status code 403
		{
			givenURL: "/ads/2",
			givenBody: map[string]any{
				"user_id": 5,
			},
			expectedStatusCode: http.StatusForbidden,
		},

		// test of getting status code 404
		{
			givenURL: "/ads/50",
			givenBody: map[string]any{
				"user_id": 2,
			},
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting status code 404
		{
			givenURL: "/ads/2",
			givenBody: map[string]any{
				"user_id": 50,
			},
			expectedStatusCode: http.StatusNotFound,
		},

		// test of getting status code 400
		{
			givenURL: "/ads/abc",
			givenBody: map[string]any{
				"user_id": 2,
			},
			expectedStatusCode: http.StatusBadRequest,
		},

		// test of getting status code 400
		{
			givenURL: "/ads/2",
			givenBody: map[string]any{
				"user_id": "2",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, m := range mocks {
		s.app.On("DeleteAd", mock.Anything, m.userID, m.adID).Return(m.err).Once()
	}

	for _, test := range tests {
		code, err := s.deleteAd(test.givenURL, test.givenBody)
		assert.Equal(s.T(), test.expectedStatusCode, code)
		assert.NoError(s.T(), err)
	}

	s.app.AssertCalled(s.T(), "DeleteAd", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64"))
}
