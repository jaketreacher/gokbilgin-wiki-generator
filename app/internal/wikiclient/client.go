package wikiclient

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
)

type TokenType string

var Token = struct {
	Csrf  TokenType
	Login TokenType
}{
	Csrf:  "csrf",
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

func (c *WikiClient) UserInfoQuery() {
	resp, err := c.client.R().SetFormData(map[string]string{
		"action": "query",
		"meta":   "userinfo",
		"uiprop": "groups|rights|ratelimits|theoreticalratelimits",
		"format": "json",
	}).Post(c.endpoint)

	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	spew.Dump(resp)
}

func (c *WikiClient) TokenQuery(token TokenType) string {
	req := c.client.R()
	if token != Token.Login {
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
		log.Fatalln("tokenquery error")
	}

	result := resp.Result().(*TokenQueryResponse)

	var value string
	switch token {
	case Token.Csrf:
		value = result.Query.Tokens.CsrfToken
	case Token.Login:
		value = result.Query.Tokens.LoginToken
	default:
		log.Fatalln("bad TokenQuery response")
	}

	if value == "" {
		log.Fatalf("empty %s response", token)
	}

	return value
}

type LoginResponse struct {
	Login struct {
		Result string
		Reason string
	}
}

func (c *WikiClient) Login(username string, password string) {
	loginToken := c.TokenQuery(Token.Login)

	resp, err := c.client.R().SetFormData(map[string]string{
		"action":     "login",
		"lgname":     username,
		"lgpassword": password,
		"lgtoken":    loginToken,
		"format":     "json",
	}).SetResult(LoginResponse{}).Post(c.endpoint)

	if err != nil {
		spew.Dump(err)
		log.Fatalln("login error, exiting...")
	}

	result := resp.Result().(*LoginResponse)

	spew.Dump(resp)
	if result.Login.Result == "Success" {
		fmt.Println("login successful")
		c.auth = resp.Cookies()
	} else {
		log.Fatalln("login failed, exiting...")
	}
}

func (c *WikiClient) Logout() {
	csrfToken := c.TokenQuery(Token.Csrf)

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

type ErrorPayload struct {
	Code string
	Info string
}

/*
If the request does not modify the page:
  - `nochange` will be prsent as an empty string.

If the request does result in a change, either:
  - `new“ will be as an empty string; or
  - `oldrevid` will be non-zero.

Additionally, when a change occurs, the corresponding
grouped fields will also be present.

Changing the title will result in a new page being created as
opposed to a revision. This is due to mediawiki treating
the title as an ID.
*/
type EditPayload struct {
	ContentModel string
	PageId       int
	Result       string
	Title        string

	New          *string `json:",omitempty"`
	OldRevId     *int    `json:",omitempty"`
	NewRevId     *int    `json:",omitempty"`
	NewTimestamp *string `json:",omitempty"`

	NoChange *string `json:",omitempty"`
}

type EditResponse struct {
	Edit  *EditPayload  `json:",omitempty"`
	Error *ErrorPayload `json:",omitempty"`
}

func (c *WikiClient) Edit(title string, text string) (*EditResponse, error) {
	csrfToken := c.TokenQuery(Token.Csrf)

	resp, err := c.client.R().SetCookies(c.auth).SetFormData(map[string]string{
		"action": "edit",
		"title":  title,
		"text":   text,
		"token":  csrfToken,
		"bot":    "true", // Mark edit as a bot edit
		"format": "json",
	}).SetResult(EditResponse{}).Post(c.endpoint)

	if err != nil {
		return nil, fmt.Errorf("error making edit request: %w", err)
	}

	result := resp.Result().(*EditResponse)

	if result.Error != nil {
		return nil, fmt.Errorf("error in edit response: %w", errors.New(result.Error.Info))
	}

	return result, nil
}
