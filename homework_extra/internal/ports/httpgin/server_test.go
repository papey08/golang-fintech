package httpgin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	mocks "homework_extra/internal/ports/httpgin/app_mocks"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpServerTestSuite struct {
	suite.Suite
	app     *mocks.App
	client  *http.Client
	srv     *http.Server
	baseURL string
}

func httpServerTestSuiteInit(s *httpServerTestSuite) {
	s.app = new(mocks.App)
	s.srv = NewHTTPServer(18080, s.app)
	testServer := httptest.NewServer(s.srv.Handler)
	s.client = testServer.Client()
	s.baseURL = testServer.URL
}

func (s *httpServerTestSuite) SetupSuite() {
	httpServerTestSuiteInit(s)
}

func (s *httpServerTestSuite) TearDownSuite() {
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("failde to shutdown server: %v", err)
	}
}

func (s *httpServerTestSuite) getResponse(req *http.Request, out any) (int, error) {
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("unexpected error: %w", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("unable to read response: %w", err)
	}
	_ = json.Unmarshal(respBody, out)
	return resp.StatusCode, nil
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(httpServerTestSuite))
}

func FuzzStrToID(f *testing.F) {
	tests := []string{
		"0",
		"2",
		"12",
		"6",
		"abc",
	}

	for _, test := range tests {
		f.Add(test)
	}

	f.Fuzz(func(t *testing.T, s string) {
		gotInt64, gotErr := strToID(s)

		if !((gotInt64 == -1 && gotErr != nil) || (gotInt64 >= 0 && gotErr == nil)) {
			t.Errorf("for %s got int64: %d and error: %v", s, gotInt64, gotErr)
		}
	})
}
