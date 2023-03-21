package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"urlShortner/internal/help"
	"urlShortner/internal/ports"
	"urlShortner/internal/repository/inMemory"
	"urlShortner/internal/repository/urlGen"
)

func TestShorten(t *testing.T) {
	s := inMemory.NewInMemory()
	e := urlGen.NewURLGenerator()
	h := NewHandlers(s, e)

	testTable := []struct {
		req                *ShortenRequest
		expectedResp       *ShortenResponse
		expectedStatusCode int
	}{
		{
			req: &ShortenRequest{
				Url: "http://www.google.com",
			},
			expectedResp: &ShortenResponse{
				ShortUrl: "a",
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			req: &ShortenRequest{
				Url:       "http://www.yahoo.com",
				CustomUrl: "HelloYahoo",
			},
			expectedResp: &ShortenResponse{
				ShortUrl: "HelloYahoo",
			},
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range testTable {
		resp := &ShortenResponse{}
		statusCode, _, err := help.DoTestRequest(http.MethodPost, "/shorten", h.Shorten, tt.req, resp)
		if err != nil {
			t.Fatalf("DoTestRequest faild: %s", err.Error())
		}

		assert.Equal(t, tt.expectedStatusCode, statusCode, "statusCode isn't equal")
		assert.Equal(t, tt.expectedResp.ShortUrl, resp.ShortUrl, "ShortUrl isn't equal")
	}
}

func TestRedirect(t *testing.T) {
	s := newMockStorage()
	e := urlGen.NewURLGenerator()
	h := NewHandlers(s, e)

	testTable := []struct {
		shortURL           string
		logURL             string
		expectedStatusCode int
	}{
		{
			shortURL:           "a",
			logURL:             "http://www.google.com",
			expectedStatusCode: http.StatusPermanentRedirect,
		},
		{
			shortURL:           "HelloYahoo",
			logURL:             "http://www.yahoo.com",
			expectedStatusCode: http.StatusPermanentRedirect,
		},
	}

	for _, tt := range testTable {
		statusCode, redirectedURL, err := help.DoTestRequest(http.MethodGet, "/"+tt.shortURL, h.Main, nil, nil)
		if err != nil {
			t.Fatalf("DoTestRequest faild: %s", err.Error())
		}

		assert.Equal(t, tt.expectedStatusCode, statusCode, "statusCode isn't equal")
		assert.Equal(t, tt.logURL, redirectedURL, "logURL isn't equal")
	}
}

func TestStats(t *testing.T) {
	s := newMockStorage()
	e := urlGen.NewURLGenerator()
	h := NewHandlers(s, e)
	last3Days := help.LastNDays(3)

	testTable := []struct {
		shortURL           string
		expectedResp       *StatsResponse
		expectedStatusCode int
	}{
		{
			shortURL: "a",
			expectedResp: &StatsResponse{
				Original:  "http://www.google.com",
				Shortened: "a",
				Visits: map[string]int{
					last3Days[0]: 1,
					last3Days[1]: 0,
					last3Days[2]: 0,
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			shortURL: "HelloYahoo",
			expectedResp: &StatsResponse{
				Original:  "http://www.yahoo.com",
				Shortened: "HelloYahoo",
				Visits: map[string]int{
					last3Days[0]: 0,
					last3Days[1]: 0,
					last3Days[2]: 0,
				},
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range testTable {
		resp := &StatsResponse{}
		statusCode, _, err := help.DoTestRequest(http.MethodGet, "/"+tt.shortURL+"/stats", h.Main, nil, resp)
		if err != nil {
			t.Fatalf("DoTestRequest faild: %s", err.Error())
		}

		assert.Equal(t, tt.expectedStatusCode, statusCode, "statusCode isn't equal")
		assert.Equal(t, tt.expectedResp.Original, resp.Original, "Original isn't equal")
		assert.Equal(t, tt.expectedResp.Shortened, resp.Shortened, "Shortened isn't equal")
		assert.Equal(t, tt.expectedResp.Visits[last3Days[2]], resp.Visits[last3Days[2]], "Visits isn't equal")
	}
}

func newMockStorage() ports.Storage {
	s := inMemory.NewInMemory()
	s.SetURLMapping("a", "http://www.google.com")
	s.SetURLMapping("HelloYahoo", "http://www.yahoo.com")
	s.IncrementStats("http://www.google.com")

	return s
}
