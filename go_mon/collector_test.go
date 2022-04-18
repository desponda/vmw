package main

import (
	"net/http"
	"testing"

	"bou.ke/monkey"
)

func TestGetMetrics(t *testing.T) {
	monkey.Patch(handleHTTPRequest, func(url string) response {
		if url == "https://httpstat.us/503" {
			return response{
				Url:          url,
				ResponseTime: 1.7,
			}
		} else {
			return response{
				Url:          url,
				ResponseTime: 1.7,
			}

		}

	})
	defer monkey.UnpatchAll()

	tests := []struct {
		name    string
		url     string
		expResp response
	}{
		{
			name: "Status  200",
			url:  "https://httpstat.us/200",
			expResp: response{
				Url:          "https://httpstat.us/200",
				ResponseTime: 1.7,
			},
		},
		{
			name: "Status 503",
			url:  "https://httpstat.us/503",
			expResp: response{
				Url:          "https://httpstat.us/503",
				ResponseTime: 1.7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp := GetMetrics(tt.url)

			if resp != tt.expResp {
				t.Errorf("expected %v, got %v", tt.expResp, resp)
			}
		})
	}
}

func TestHandleHTTPRequest(t *testing.T) {
	monkey.Patch(http.Get, func(url string) (resp *http.Response, err error) {
		if url == "https://httpstat.us/503" {
			return &http.Response{
				StatusCode: 503,
			}, nil
		} else {
			return &http.Response{
				StatusCode: 200,
			}, nil
		}

	})
	defer monkey.UnpatchAll()
	tests := []struct {
		name          string
		url           string
		expResp       *http.Response
		ExternalUrlUp float64
	}{
		{
			name: "Status  200",
			url:  "https://httpstat.us/200",
			expResp: &http.Response{
				StatusCode: 200,
			},
			ExternalUrlUp: 1,
		},
		{
			name: "Status 503",
			url:  "https://httpstat.us/503",
			expResp: &http.Response{
				StatusCode: 503,
			},
			ExternalUrlUp: 0,
		},
	}
	for _, tt := range tests {

		resp := handleHTTPRequest(tt.url)
		if resp.ExternalUrlUp != tt.ExternalUrlUp {
			t.Errorf("expected %v, got %v", tt.ExternalUrlUp, resp.ExternalUrlUp)
		}
	}

}
