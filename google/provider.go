package google

import (
	"context"

	"github.com/gokhanaltun/o2go"
)

type Google struct {
	oauth2 *o2go.OAuth2
}

func (p *Google) AuthURL() string {
	return "https://accounts.google.com/o/oauth2/auth"
}

func (p *Google) TokenURL() string {
	return "https://oauth2.googleapis.com/token"
}

func (p *Google) AuthCodeURL(state string) (string, error) {
	return p.oauth2.AuthCodeURL(state)
}

func (p *Google) ExchangeAuthCode(ctx context.Context, code string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeAuthCode(ctx, code)
}

func (p *Google) ExchangeRefreshToken(ctx context.Context, refreshToken string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeRefreshToken(ctx, refreshToken)
}
