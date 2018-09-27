package handlers

import (
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
    "strings"
	"net/url"
	"log"
	

	"github.com/OtgonbayarT/microservice/models"
) 

func TestHandlers_Handler (t *testing.T) {
	
	data := url.Values{}
	data.Set("url", "http://studynihongo101.blogspot.com/search?updated-max=2012-05-03T05:31:00-07:00&max-results=1&start=1&by-date=false")
	req := httptest.NewRequest("POST", "/encode", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dataEmpty := url.Values{}
	dataEmpty.Set("url", "")
	reqDataEmpty := httptest.NewRequest("POST", "/encode", strings.NewReader(dataEmpty.Encode()))
	reqDataEmpty.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		{	
			name: 			"encodeEmpty",
			in:				reqDataEmpty,
			out:			httptest.NewRecorder(),
			expectedStatus:	http.StatusBadRequest,
			expectedBody:	"body 'url' is missing",
		},
		{	
			name: 			"decodeEmpty",
			in:				httptest.NewRequest("GET", "/decode/", nil),
			out:			httptest.NewRecorder(),
			expectedStatus:	http.StatusBadRequest,
			expectedBody:	"there is nothing to decode",
		},
		{	
			name: 			"redirectEmpty",
			in:				httptest.NewRequest("GET", "/redirect/", nil),
			out:			httptest.NewRecorder(),
			expectedStatus:	http.StatusBadRequest,
			expectedBody:	"there is nothing to redirect",
		},
		
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func (t *testing.T) {
			h := NewHandlersLog(log.New(os.Stdout, "url shortener service: ", log.LstdFlags | log.Lshortfile), "./urls.db", nil)
			
			if (test.name == "encode" || test.name == "encodeEmpty" || test.name == "encodeDataurlNotEncoded") {
				h.EncodeHandler(test.out, test.in)
			} else if (test.name == "decode" || test.name == "decodeEmpty") {
				h.DecodeHandler(test.out, test.in)
			} else if (test.name == "redirect" || test.name == "redirectEmpty") {
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