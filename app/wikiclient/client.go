package wikiclient

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
)

type WikiClient struct {
	client   *resty.Client
	endpoint string
	auth     string
}

func New(endpoint string) *WikiClient {
	client := resty.New()
	return &WikiClient{
		client,
		endpoint,
		"SomeCookie",
	}
}

type TokenQueryResponse struct {
	Query struct {
		Tokens struct {
			CsrfToken string
		}
	}
}

func (c *WikiClient) TokenQuery() string {
	resp, err := c.client.R().SetHeaders(map[string]string{
		"Cookie": c.auth,
	}).SetFormData(map[string]string{
		"action": "query",
		"meta":   "tokens",
		"format": "json",
	}).SetResult(TokenQueryResponse{}).Post(c.endpoint)

	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	result := resp.Result().(*TokenQueryResponse)

	return result.Query.Tokens.CsrfToken
}
