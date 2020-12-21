package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPI(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	api, err := NewAPI(values)
	assert.NoError(t, err)

	for _, value := range values {
		assert.Equal(t, true, api.treeNode.Search(value))
	}

	values = []int{1, 1}
	_, err = NewAPI(values)
	assert.Error(t, err)
}

func TestSearch(t *testing.T) {
	api, err := NewAPI([]int{1, 2, 3, 10})
	require.NoError(t, err)

	testCases := []struct {
		name       string
		value      string
		response   Response
		statusCode int
	}{
		{
			name:  "case found",
			value: "10",
			response: Response{
				Success: true,
				Message: "OK",
				Data:    true,
			},
			statusCode: http.StatusOK,
		},
		{
			name:  "case not found",
			value: "111",
			response: Response{
				Success: true,
				Message: "OK",
				Data:    false,
			},
			statusCode: http.StatusOK,
		},
		{
			name:  "case wrong query param value",
			value: "foo",
			response: Response{
				Success: false,
				Message: `strconv.Atoi: parsing "foo": invalid syntax`,
				Data:    nil,
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := make(url.Values)
			q.Set("val", tc.value)
			req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

			rec := httptest.NewRecorder()
			ctx := api.echo.NewContext(req, rec)

			err = api.search(ctx)
			assert.NoError(t, err)

			assert.Equal(t, tc.statusCode, rec.Code)

			resp, err := json.Marshal(&tc.response)
			assert.NoError(t, err)
			assert.Equal(t, string(resp), strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestInsert(t *testing.T) {
	api, err := NewAPI([]int{1, 2, 3, 10})
	require.NoError(t, err)

	testCases := []struct {
		name        string
		reqBody     map[string]int
		response    Response
		statusCode  int
		contentType string
	}{
		{
			name: "case success insertion",
			reqBody: map[string]int{
				"val": 7,
			},
			response: Response{
				Success: true,
				Message: "OK",
				Data:    nil,
			},
			statusCode:  http.StatusCreated,
			contentType: echo.MIMEApplicationJSON,
		},
		{
			name: "case insert existing",
			reqBody: map[string]int{
				"val": 1,
			},
			response: Response{
				Success: false,
				Message: "this node value already exists",
				Data:    nil,
			},
			statusCode:  http.StatusInternalServerError,
			contentType: echo.MIMEApplicationJSON,
		},
		{
			name: "case wrong content type",
			response: Response{
				Success: false,
				Message: "code=415, message=Unsupported Media Type",
				Data:  nil,
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBytes, err := json.Marshal(tc.reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBytes))
			req.Header.Set(echo.HeaderContentType, tc.contentType)

			rec := httptest.NewRecorder()
			ctx := api.echo.NewContext(req, rec)

			err = api.insert(ctx)
			assert.NoError(t, err)

			assert.Equal(t, tc.statusCode, rec.Code)

			resp, err := json.Marshal(&tc.response)
			assert.NoError(t, err)
			assert.Equal(t, string(resp), strings.TrimSpace(rec.Body.String()))

			value, ok := tc.reqBody["val"]
			if ok {
				assert.Equal(t, true, api.treeNode.Search(value))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	api, err := NewAPI([]int{1, 2, 3, 10})
	require.NoError(t, err)

	testCases := []struct {
		name       string
		value      string
		response   Response
		statusCode int
	}{
		{
			name:       "case found",
			value:      "3",
			statusCode: http.StatusNoContent,
		},
		{
			name:       "remove not-existing",
			value:      "3",
			statusCode: http.StatusInternalServerError,
			response: Response{
				Success: false,
				Message: "value to be deleted does not exist in the tree",
				Data: nil,
			},
		},
		{
			name:       "case wrong query param value",
			value:      "foobar",
			statusCode: http.StatusBadRequest,
			response: Response{
				Success: false,
				Message: `strconv.Atoi: parsing "foobar": invalid syntax`,
				Data: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := make(url.Values)
			q.Set("val", tc.value)
			req := httptest.NewRequest(http.MethodDelete, "/?"+q.Encode(), nil)

			rec := httptest.NewRecorder()
			ctx := api.echo.NewContext(req, rec)

			err = api.delete(ctx)
			assert.NoError(t, err)

			assert.Equal(t, tc.statusCode, rec.Code)

			if rec.Code != http.StatusNoContent {
				resp, err := json.Marshal(&tc.response)
				assert.NoError(t, err)
				assert.Equal(t, string(resp), strings.TrimSpace(rec.Body.String()))
			}

			value, err := strconv.Atoi(tc.value)
			if err == nil {
				assert.Equal(t, false, api.treeNode.Search(value))
			}
		})
	}
}
