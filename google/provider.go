package google

import (
	"context"

	"github.com/gokhanaltun/o2go/v2"
)

type Google struct {
	oauth2 *o2go.OAuth2
}

func (p *Google) AuthURL() string {
	return "https://accounts.google.com/o/oauth2/v2/auth"
}

func (p *Google) TokenURL() string {
	return "https://oauth2.googleapis.com/token"
}

func (p *Google) IssuerURL() *string {
	value := "https://accounts.google.com"
	return &value
}

func (p *Google) JWKSURL() *string {
	value := "https://www.googleapis.com/oauth2/v3/certs"
	return &value
}

func (p *Google) AuthCodeURL(state string, params map[string]string) (string, error) {
	return p.oauth2.AuthCodeURL(state, params)
}

func (p *Google) ExchangeAuthCode(ctx context.Context, code string, params map[string]string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeAuthCode(ctx, code, params)
}

func (p *Google) ExchangeRefreshToken(ctx context.Context, refreshToken string, params map[string]string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeRefreshToken(ctx, refreshToken, params)
}

func (p *Google) ParseClaims(response o2go.ExchangeResponse, params map[string]string) (map[string]any, error) {
	return p.oauth2.ParseClaims(response, params)
}

func (p *Google) ParseClaimsWithOptions(ctx context.Context, response o2go.ExchangeResponse, options o2go.ParseClaimsOptions) (map[string]any, error) {
	return p.oauth2.ParseClaimsWithOptions(ctx, response, options)
}
