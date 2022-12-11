package adepttech

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	defaultBaseURL        = "https://api.adept.tech/v1/api"
	defaultAuthorizeURL   = "https://api.adept.tech/v1/authorize"
	defaultAccessTokenURL = "https://api.adept.tech/v1/access_token"
	tokenExchangeTimeout  = 15 * time.Second
)

func NewInstance(cfg Config) (Instance, error) {
	// overload defaults
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.AuthorizeURL == "" {
		cfg.AuthorizeURL = defaultAuthorizeURL
	}
	if cfg.AccessTokenURL == "" {
		cfg.AccessTokenURL = defaultAccessTokenURL
	}
	if cfg.Instance == "" {
		return nil, errors.New("adepttech: missing Instance in configuration")
	}
	if cfg.RedirectURL == "" {
		return nil, errors.New("adepttech: missing RedirectURL in configuration")
	}
	if cfg.ClientID == "" {
		return nil, errors.New("adepttech: missing ClientID in configuration")
	}
	if cfg.ClientSecret == "" {
		return nil, errors.New("adepttech: missing ClientSecret in configuration")
	}
	authUrl, err := url.Parse(cfg.AuthorizeURL)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse AuthorizeURL")
	}
	authUrl = appendInstanceURL(authUrl, cfg.Instance)
	tokenUrl, err := url.Parse(cfg.AccessTokenURL)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse AccessTokenURL")
	}
	tokenUrl = appendInstanceURL(tokenUrl, cfg.Instance)

	return &instance{config: cfg, oauth2Config: &oauth2.Config{
		RedirectURL:  cfg.RedirectURL,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{"stats", "email", "basic"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authUrl.String(),
			TokenURL: tokenUrl.String(),
		}}}, nil
}

type instance struct {
	config       Config
	oauth2Config *oauth2.Config
	oauth2Client *http.Client
	oauth2Token  *oauth2.Token
}

func (i *instance) AuthURL(state string) string {
	return i.oauth2Config.AuthCodeURL(state)
}
func (i *instance) Token() (*oauth2.Token, error) {
	return i.oauth2Token, nil
}

func (i *instance) AssignTokenSource(ctx context.Context, tokenSource oauth2.TokenSource) error {
	token, err := tokenSource.Token()
	if err != nil {
		return err
	}
	i.oauth2Token = token
	i.oauth2Client = oauth2.NewClient(ctx, tokenSource)
	return nil
}

func (i *instance) ExchangeCode(ctx context.Context, code string) error {
	httpClient := &http.Client{Timeout: tokenExchangeTimeout}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	token, err := i.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return err
	}
	i.oauth2Token = token
	i.oauth2Client = i.oauth2Config.Client(ctx, token)
	return nil
}
func (i *instance) Get(ctx context.Context, path string, params url.Values) (*http.Response, error) {
	if i.oauth2Client == nil {
		return nil, fmt.Errorf("Get called but oauth2 client not initiated")
	}
	requestUrl, err := url.Parse(fmt.Sprintf("%s/%s", i.config.BaseURL, path))
	if err != nil {
		return nil, err
	}
	requestUrl.RawQuery = params.Encode()
	requestUrl = appendInstanceURL(requestUrl, i.config.Instance)

	req, err := http.NewRequestWithContext(ctx, "GET", requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	return i.oauth2Client.Do(req)
}

func (i *instance) GetInto(ctx context.Context, path string, params url.Values, target interface{}) error {
	res, err := i.Get(ctx, path, params)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if !json.Valid(body) {
		return fmt.Errorf("invalid Json Response with content length %d: %s", res.ContentLength, string(body))
	}
	err = json.Unmarshal(body, target)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to unmarshal into type %T: %s", target, body))
	}
	return nil
}

func appendInstanceURL(u *url.URL, instance string) *url.URL {
	uValues := u.Query()
	uValues.Add("instance", instance)
	u.RawQuery = uValues.Encode()
	return u
}
