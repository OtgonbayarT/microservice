package handlers

import (
	"testing"
	"net/http"
	"net/http/httptest"
    "bytes"
	"net/url"
)

func TestHandlers_Handler (t *testing.T) {
	data := url.Values{}
	data.Set("url", `http://studynihongo101.blogspot.com/search?updated-max=2012-05-03T05:31:00-07:00&max-results=1&start=1&by-date=false`)
	req := httptest.NewRequest("POST", "/encode", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	tests := []struct {
		name 			string
		in 				*http.Request
		out 			*httptest.ResponseRecorder
		expectedStatus 	int
		expectedBody 	string
	}{
		{
			name: 			"encode",
			in:				req,
			out:			httptest.NewRecorder(),
			expectedStatus:	http.StatusOK,
			expectedBody:	"4106327431",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func (t *testing.T) {
			h := NewHandlersLog(nil, "./urls.db")
			h.EncodeHandler(test.out, test.in)
			if test.out.Code != test.expectedStatus {
				t.Logf("expected: %d\ngot: %d\n", test.expectedStatus, test.out.Code)
				t.Fail()
			}
			body := test.out.Body.String()
			if body != test.expectedBody {
				t.Logf("expected: %s\ngot: %s\n", test.expectedBody, body)
				t.Fail()
			}
		})
	}
}