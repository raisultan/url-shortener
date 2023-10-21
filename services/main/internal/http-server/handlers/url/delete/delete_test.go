package delete_test

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/raisultan/url-shortener/services/main/internal/http-server/handlers/url/delete"
	"github.com/raisultan/url-shortener/services/main/internal/http-server/handlers/url/delete/mocks"
	"github.com/raisultan/url-shortener/services/main/internal/lib/logger/handlers/slogdiscard"
	"github.com/raisultan/url-shortener/services/main/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
		{
			name:      "Alias Not Found",
			alias:     "missing_alias",
			respError: "alias not found",
			mockError: storage.ErrUrlNotFound,
		},
		{
			name:      "Unexpected Error",
			alias:     "test_alias",
			respError: "failed to delete url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlDeleterMock := mocks.NewUrlDeleter(t)

			urlDeleterMock.On(
				"DeleteUrl",
				mock.Anything,
				tc.alias,
			).Return(tc.mockError).Once()

			r := chi.NewRouter()
			r.Delete("/{alias}", delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			req, err := http.NewRequest(http.MethodDelete, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			urlDeleterMock.AssertExpectations(t)

			var respBody map[string]string
			err = json.Unmarshal(rr.Body.Bytes(), &respBody)
			require.NoError(t, err)
			respErr, ok := respBody["error"]

			if tc.respError == "" {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				assert.Equal(t, tc.respError, respErr)
			}
		})
	}
}
