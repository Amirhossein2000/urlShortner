package help

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"
	"urlShortner/internal/consts"
)

func SendJsonResponse(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"message": http.StatusText(statusCode),
	})
}

func SendJsonResponseWithCustomMessage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func SendJsonResponseWithBody(w http.ResponseWriter, statusCode int, response interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func IsValidURL(longURL string) bool {
	u, err := url.Parse(longURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func LastNDays(n int) []string {
	days := make([]string, n)
	now := time.Now()
	for i := 0; i < n; i++ {
		days[i] = now.Add(-time.Duration(i) * 24 * time.Hour).Format(consts.TimeLayout)
	}
	return days
}

func HttpRequestDetails(r *http.Request) string {
	return fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
}

func DoTestRequest(method, url string, h http.HandlerFunc, request interface{}, response interface{}) (int, string, error) {
	body := bytes.NewBuffer([]byte{})
	switch request.(type) {
	case string:
		body.Write([]byte(request.(string)))
	default:
		err := json.NewEncoder(body).Encode(request)
		if err != nil {
			return -1, "", err
		}
	}

	req := httptest.NewRequest(method, url, body)
	recorder := httptest.NewRecorder()
	h(recorder, req)

	if response != nil {
		err := json.NewDecoder(recorder.Result().Body).Decode(response)
		if err != nil {
			return -1, "", err
		}
	}

	redirectedURL := ""
	location, _ := recorder.Result().Location()
	if location != nil {
		redirectedURL = location.String()
	}

	return recorder.Result().StatusCode, redirectedURL, nil
}
