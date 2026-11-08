package github

import (
	"context"
	"net/http"
	"strings"

	"github.com/gokhanaltun/o2go/v2"
)

type Github struct {
	oauth2 *o2go.OAuth2
}

func (p *Github) AuthURL() string {
	return "https://github.com/login/oauth/authorize"
}

func (p *Github) TokenURL() string {
	return "https://github.com/login/oauth/access_token"
}

func (p *Github) IssuerURL() *string { return nil }

func (p *Github) JWKSURL() *string { return nil }

func (p *Github) AuthCodeURL(state string, params map[string]string) (string, error) {
	return p.oauth2.AuthCodeURL(state, params)
}

func (p *Github) ExchangeAuthCode(ctx context.Context, code string, params map[string]string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeAuthCode(ctx, code, params)
}

func (p *Github) ExchangeRefreshToken(ctx context.Context, refreshToken string, params map[string]string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeRefreshToken(ctx, refreshToken, params)
}

type acceptJSONTransport struct{ base http.RoundTripper }

func (t acceptJSONTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.EqualFold(req.URL.Hostname(), "github.com") && req.URL.Path == "/login/oauth/access_token" {
		request := req.Clone(req.Context())
		request.Header = req.Header.Clone()
		request.Header.Set("Accept", "application/json")
		req = request
	}
	return t.base.RoundTrip(req)
}

func withJSONTokenResponse() o2go.Option {
	return func(client *o2go.OAuth2) {
		baseClient := client.HTTPClient
		if baseClient == nil {
			baseClient = http.DefaultClient
		}
		baseTransport := baseClient.Transport
		if baseTransport == nil {
			baseTransport = http.DefaultTransport
		}
		client.HTTPClient = &http.Client{
			Transport:     acceptJSONTransport{base: baseTransport},
			CheckRedirect: baseClient.CheckRedirect,
			Jar:           baseClient.Jar,
			Timeout:       baseClient.Timeout,
		}
	}
}
