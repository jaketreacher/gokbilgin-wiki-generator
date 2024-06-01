package wikiclient

import (
	"fmt"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
)

type WikiClient struct {
	client   *resty.Client
	endpoint string
	auth     []*http.Cookie
}

func New(endpoint string) *WikiClient {
	client := resty.New()
	return &WikiClient{
		client,
		endpoint,
		nil,
	}
}

type CsrfTokenQueryResponse struct {
	Query struct {
		Tokens struct {
			CsrfToken string
		}
	}
}

type LoginTokenQueryResponse struct {
	Query struct {
		Tokens struct {
			LoginToken string
		}
	}
}

func (c *WikiClient) CsrfTokenQuery() string {
	resp, err := c.client.R().SetCookies(c.auth).SetFormData(map[string]string{
		"action": "query",
		"meta":   "tokens",
		"type":   "csrf",
		"format": "json",
	}).SetResult(CsrfTokenQueryResponse{}).Post(c.endpoint)

	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	result := resp.Result().(*CsrfTokenQueryResponse)

	return result.Query.Tokens.CsrfToken
}

func (c *WikiClient) Login(username string, password string) {
	loginToken := c.LoginTokenQuery()

	resp, err := c.client.R().SetFormData(map[string]string{
		"action":     "login",
		"lgname":     username,
		"lgpassword": password,
		"lgtoken":    loginToken,
		"format":     "json",
	}).Post(c.endpoint)

	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	if resp.StatusCode() == 200 {
		fmt.Println("login successful")
		c.auth = resp.Cookies()
	} else {
		fmt.Println("login failed, exiting...")
		os.Exit(1)
	}
}

func (c *WikiClient) Logout() {
	csrfToken := c.CsrfTokenQuery()

	resp, err := c.client.R().SetCookies(c.auth).SetFormData(map[string]string{
		"action": "logout",
		"token":  csrfToken,
		"format": "json",
	}).Post(c.endpoint)

	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	if resp.StatusCode() == 200 {
		fmt.Println("logout successful")
	} else {
		fmt.Println("logout failed")
	}
}

func (c *WikiClient) LoginTokenQuery() string {
	resp, err := c.client.R().SetFormData(map[string]string{
		"action": "query",
		"meta":   "tokens",
		"type":   "login",
		"format": "json",
	}).SetResult(LoginTokenQueryResponse{}).Post(c.endpoint)

	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	result := resp.Result().(*LoginTokenQueryResponse)

	return result.Query.Tokens.LoginToken
}
