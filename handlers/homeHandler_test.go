package handlers

import (
	"testing"
	"net/http"
	"net/http/httptest"
    "bytes"
	"net/url"

	"github.com/OtgonbayarT/microservice/models"
) 

func TestHandlers_Handler (t *testing.T) {
			
	data := url.Values{}
	data.Set("url", `http://studynihongo101.blogspot.com/search?updated-max=2012-05-03T05:31:00-07:00&max-results=1&start=1&by-date=false`)
	req := httptest.NewRequest("POST", "/encode", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	shortUrl, _ := models.InsertUrl("./urls.db", `http://studynihongo101.blogspot.com/search?updated-max=2012-05-03T05:31:00-07:00&max-results=1&start=1&by-date=false`)	

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
			expectedBody:	"example.com/decode/4106327431",
		},
		{
			name: 			"decode",
			in:				httptest.NewRequest("GET", "/decode/"+shortUrl, nil),
			out:			httptest.NewRecorder(),
			expectedStatus:	http.StatusOK,
			expectedBody:	`http://studynihongo101.blogspot.com/search?updated-max=2012-05-03T05:31:00-07:00&max-results=1&start=1&by-date=false`,
		},
		{
			name: 			"redirect",
			in:				httptest.NewRequest("GET", "/redirect/"+shortUrl, nil),
			out:			httptest.NewRecorder(),
			expectedStatus:	http.StatusMovedPermanently,
			expectedBody:	`<a href="http://studynihongo101.blogspot.com/search?updated-max=2012-05-03T05:31:00-07:00&amp;max-results=1&amp;start=1&amp;by-date=false">Moved Permanently</a>`,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func (t *testing.T) {
			h := NewHandlersLog(nil, "./urls.db")
			
			if (test.name == "encode") {
				h.EncodeHandler(test.out, test.in)
			} else if (test.name == "decode") {
				h.DecodeHandler(test.out, test.in)
			} else if (test.name == "redirect") {
				h.RedirectHandler(test.out, test.in)
			}	

			if test.out.Code != test.expectedStatus {
				t.Logf("expected: %d\ngot: %d\n", test.expectedStatus, test.out.Code)
				t.Fail()
			}

			body := test.out.Body.String()
			if test.out.Code != http.StatusMovedPermanently {
				if body != test.expectedBody {
					t.Logf("expected: %s\ngot: %s\n", test.expectedBody, body)
					t.Fail()
				}
			}
		})
	}
}