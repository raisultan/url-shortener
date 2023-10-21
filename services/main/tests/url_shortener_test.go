package tests

import (
	"fmt"
	"github.com/raisultan/url-shortener/services/main/internal/http-server/handlers/url/save"
	"github.com/raisultan/url-shortener/services/main/internal/lib/api"
	"github.com/raisultan/url-shortener/services/main/internal/lib/random"
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8080"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			Url:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("alias")
}

//nolint:funlen
func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field Url is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save

			resp := e.POST("/url").
				WithJSON(save.Request{
					Url:   tc.url,
					Alias: tc.alias,
				}).
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			testRedirect(t, alias, tc.url)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}

func TestURLShortener_Delete(t *testing.T) {
	testCases := []struct {
		name             string
		alias            string
		expectedResponse string
		isError          bool
	}{
		{
			name:             "Valid Alias",
			alias:            gofakeit.Word() + gofakeit.Word(),
			expectedResponse: "OK",
		},
		{
			name:             "Invalid Alias",
			alias:            gofakeit.Word(),
			expectedResponse: "alias not found",
			isError:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			if tc.expectedResponse == "OK" {
				e.POST("/url").
					WithJSON(save.Request{
						Url:   gofakeit.URL(),
						Alias: tc.alias,
					}).
					Expect().Status(http.StatusOK).
					JSON().Object().Value("alias").String().IsEqual(tc.alias)
			}

			deleteResp := e.DELETE(fmt.Sprintf("/%s", tc.alias)).
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.isError == true {
				deleteResp.Value("error").String().IsEqual(tc.expectedResponse)
			} else {
				deleteResp.Value("status").String().IsEqual(tc.expectedResponse)
			}
		})
	}
}
