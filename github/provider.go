package github

import (
	"context"

	"github.com/gokhanaltun/o2go"
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

func (p *Github) AuthCodeURL(state string) (string, error) {
	return p.oauth2.AuthCodeURL(state)
}

func (p *Github) ExchangeAuthCode(ctx context.Context, code string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeAuthCode(ctx, code)
}

func (p *Github) ExchangeRefreshToken(ctx context.Context, refreshToken string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeRefreshToken(ctx, refreshToken)
}
