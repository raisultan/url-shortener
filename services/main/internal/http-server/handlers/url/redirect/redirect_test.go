package redirect_test

import (
	"errors"
	"github.com/raisultan/url-shortener/internal/lib/api"
	"github.com/raisultan/url-shortener/internal/storage"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/raisultan/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/raisultan/url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"github.com/raisultan/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
		{
			name:      "URL not found",
			alias:     "test_alias",
			respError: "url not found",
			mockError: storage.ErrUrlNotFound,
		},
		{
			name:      "Internal error",
			alias:     "test_alias",
			respError: "internal error",
			mockError: errors.New("some internal error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewUrlGetter(t)

			urlGetterMock.On("GetUrl", mock.Anything, tc.alias).
				Return(tc.url, tc.mockError).Once()

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			if tc.mockError != nil {
				resp, err := http.Get(ts.URL + "/" + tc.alias)
				require.NoError(t, err)
				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				_ = resp.Body.Close()
				assert.Contains(t, string(bodyBytes), tc.respError)
			} else {
				redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
				require.NoError(t, err)
				assert.Equal(t, tc.url, redirectedToURL)
			}
		})
	}
}
