package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vysogota99/advertising/internal/app/store/mock"
	"github.com/stretchr/testify/assert"
)

const (
	serverPort = ":8081"
)

func TestCreateAdHandler(t *testing.T) {
	type testCase struct {
		name   string
		params map[string]interface{}
		code   int
	}

	tCases := []testCase{
		{
			name: "correct data",
			params: map[string]interface{}{
				"name":        "iphone",
				"description": "Продам айфон",
				"photos": []string{
					"https://vk.com/",
					"https://vk.com/",
					"https://vk.com/",
				},
				"price": 100,
			},
			code: 200,
		},
		{
			name: "invalid photos list",
			params: map[string]interface{}{
				"name":        "iphone",
				"description": "Продам айфон",
				"photos": []string{
					"https://vk.com/",
					"https://vk.com/",
					"https://vk.com/",
					"https://vk.com/",
				},
				"price": 100,
			},
			code: 400,
		},
		{
			name: "invalid name",
			params: map[string]interface{}{
				"name":        "",
				"description": "Продам айфон",
				"photos": []string{
					"https://vk.com/",
					"https://vk.com/",
					"https://vk.com/",
					"https://vk.com/",
				},
				"price": 100,
			},
			code: 400,
		},
		{
			name: "invalid description",
			params: map[string]interface{}{
				"name":        "iphone",
				"description": "",
				"photos": []string{
					"https://vk.com/",
					"https://vk.com/",
					"https://vk.com/",
					"https://vk.com/",
				},
				"price": 100,
			},
			code: 400,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			store := mock.New()
			r := NewRouter(serverPort, store)
			w := httptest.NewRecorder()
			body, err := json.Marshal(tc.params)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/ad", bytes.NewBuffer(body))
			assert.NoError(t, err)

			router := r.Setup()
			assert.NoError(t, err)

			router.ServeHTTP(w, req)
			assert.Equal(t, tc.code, w.Result().StatusCode)

		})
	}

}

func TestGetAdHandler(t *testing.T) {
	type testCase struct {
		name string
		urn  string
		code int
	}

	tCases := []testCase{
		{
			name: "correct fields",
			urn:  "/ad/1?description=true&photos=true",
			code: 200,
		},
		{
			name: "invalid path",
			urn:  "/ad/-1?description=true&photos=true",
			code: 400,
		},
		{
			name: "invalid field",
			urn:  "/ad/1?description=true&photos=-1",
			code: 400,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			store := mock.New()
			r := NewRouter(serverPort, store)
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", tc.urn, nil)
			assert.NoError(t, err)

			router := r.Setup()
			assert.NoError(t, err)

			router.ServeHTTP(w, req)
			assert.Equal(t, tc.code, w.Result().StatusCode)

		})
	}

}

func TestGetAdsHandeler(t *testing.T) {
	store := mock.New()
	r := NewRouter(serverPort, store)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/ads?p=1", nil)
	assert.NoError(t, err)

	router := r.Setup()
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
