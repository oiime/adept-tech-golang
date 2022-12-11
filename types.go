package adepttech

import (
	"context"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type Config struct {
	Instance       string
	BaseURL        string
	RedirectURL    string
	AuthorizeURL   string
	AccessTokenURL string
	ClientID       string
	ClientSecret   string
}

type Instance interface {
	AuthURL(state string) string
	Token() (*oauth2.Token, error)
	Get(ctx context.Context, path string, params url.Values) (*http.Response, error)
	GetInto(ctx context.Context, path string, params url.Values, target interface{}) error
	ExchangeCode(ctx context.Context, code string) error
	AssignTokenSource(ctx context.Context, tokenSource oauth2.TokenSource) error
}
