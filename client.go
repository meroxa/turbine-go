package valve

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"golang.org/x/oauth2"
	"os"
)

type ClientConfig struct {
	AccessToken  string `env:"ACCESS_TOKEN"`
	RefreshToken string `env:"REFRESH_TOKEN"`
	AuthAudience string `env:"MEROXA_AUTH_AUDIENCE" envDefault:"https://api.meroxa.io/v1"`
	AuthDomain   string `env:"MEROXA_AUTH_DOMAIN" envDefault:"auth.meroxa.io"`
	AuthClientID string `env:"MEROXA_AUTH_CLIENT_ID" envDefault:"2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb"`
}

const Version = "0.1.0"

var cfg ClientConfig

type Client struct {
	client meroxa.Client
}

func (c Client) GetResource(name string) (*Resource, error) {
	if c.client != nil {
		_, err := c.client.GetResourceByNameOrID(context.Background(), name)
		if err != nil {
			return nil, err
		}
		// TODO: convert meroxa.Resource to valve.Resource and return
	}

	return nil, nil
}

func NewClient(local bool) (Client, error) {
	if local {
		return Client{}, nil
	}

	c, err := newClient()
	if err != nil {
		return Client{}, err
	}

	return Client{c}, nil
}

func newClient() (meroxa.Client, error) {
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	options := []meroxa.Option{
		meroxa.WithUserAgent(fmt.Sprintf("Meroxa CLI %s", Version)),
	}

	if overrideAPIURL := os.Getenv("API_URL"); overrideAPIURL != "" {
		options = append(options, meroxa.WithBaseURL(overrideAPIURL))
	}

	options = append(options, meroxa.WithAuthentication(
		&oauth2.Config{
			ClientID: cfg.AuthClientID,
			Endpoint: oauthEndpoint(cfg.AuthDomain),
		},
		cfg.AccessToken,
		cfg.RefreshToken,
		onTokenRefreshed,
	))

	return meroxa.New(options...)
}

func oauthEndpoint(domain string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("https://%s/authorize", domain),
		TokenURL: fmt.Sprintf("https://%s/oauth/token", domain),
	}
}

func onTokenRefreshed(token *oauth2.Token) {
	cfg.AccessToken = token.AccessToken
	cfg.RefreshToken = token.RefreshToken
}