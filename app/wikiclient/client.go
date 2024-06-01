package wikiclient

import (
	"fmt"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
)

type TokenType string

var TOKEN = struct {
	CSRF  TokenType
	Login TokenType
}{
	CSRF:  "csrf",
	Login: "login",
}

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

type TokenQueryResponse struct {
	Query struct {
		Tokens struct {
			CsrfToken  string
			LoginToken string
		}
	}
}

func (c *WikiClient) TokenQuery(token TokenType) string {
	req := c.client.R()
	if token != TOKEN.Login {
		req = req.SetCookies(c.auth)
	}
	resp, err := req.SetFormData(map[string]string{
		"action": "query",
		"meta":   "tokens",
		"type":   string(token),
		"format": "json",
	}).SetResult(TokenQueryResponse{}).Post(c.endpoint)

	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	result := resp.Result().(*TokenQueryResponse)

	switch token {
	case TOKEN.CSRF:
		return result.Query.Tokens.CsrfToken
	case TOKEN.Login:
		return result.Query.Tokens.LoginToken
	default:
		fmt.Println("Bad TokenQuery request")
		os.Exit(1)
	}

	return ""
}

func (c *WikiClient) Login(username string, password string) {
	loginToken := c.TokenQuery(TOKEN.Login)

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
	csrfToken := c.TokenQuery(TOKEN.CSRF)

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
