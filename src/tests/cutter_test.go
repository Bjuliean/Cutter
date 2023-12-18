package tests

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"rapi/rapi/src/internal/http-server/handlers/url/save"
	"rapi/rapi/src/internal/random"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const host = "localhost:8080"

func TestCutter_RandomPath(t *testing.T) {
	u := url.URL{
		Host:   host,
		Scheme: "http",
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.RandomString(10),
		}).
		WithBasicAuth("user", "pass").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")
}

func TestCutter_Handlers(t *testing.T) {
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
			error: "field URL is not a valid URL",
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

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("user", "pass").
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

	redirectedToURL, err := GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}

func GetRedirect(url string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%w: %d", errors.New("invalid status code"), resp.StatusCode)
	}

	return resp.Header.Get("Location"), nil
}