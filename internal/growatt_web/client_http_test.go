package growatt_web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type responseData struct {
	AnInt   int    `json:"an_int"`
	AString string `json:"a_string"`
}

func Test_postRequest(t *testing.T) {
	testData := []struct {
		url          string
		data         url.Values
		responseBody any
		statusCode   int
		body         string
		errResult    bool
	}{
		{ // invalid body
			responseBody: responseData{},
			body:         `{invalid json}`,
			errResult:    true,
		},
		{ // error statusCode
			statusCode: 500,
			errResult:  true,
		},
		// { // io.ReadAll() fails - how...?
		// 	errResult: true,
		// },
		// { // http.Client.Do() fails - how...?
		// 	errResult: true,
		// },
		{ // http.NewRequest fails
			url:       "xyz://abc:8000f",
			errResult: true,
		},
		{ // ok with body
			responseBody: responseData{},
			body:         `{"an_int":22,"a_string":"test"}`,
		},
		{ // ok without body
			body: "",
		},
	}

	for inx, td := range testData {
		t.Run(fmt.Sprintf("#%d", inx), func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

					ctHdr, ctHdrExists := req.Header["Content-Type"]
					assert.True(t, ctHdrExists)
					assert.Equal(t, "application/x-www-form-urlencoded", ctHdr[0])

					_, uaHdrExists := req.Header["User-Agent"]
					assert.True(t, uaHdrExists)

					wr.WriteHeader(td.statusCode)
					wr.Write([]byte(td.body))
				}))

			client := httpClient{client: server.Client()}
			if len(td.url) == 0 {
				td.url = server.URL
			}
			if td.statusCode == 0 {
				td.statusCode = 200
			}

			//			var data TokenResponse
			err := client.postForm(td.url, url.Values{}, td.responseBody)
			assert.Equal(t, td.errResult, err != nil)
			server.Close()
		})
	}
}
