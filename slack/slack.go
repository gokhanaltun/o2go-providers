package slack

import (
	"context"

	"github.com/gokhanaltun/o2go/v2"
)

type Login struct {
	oauth2 *o2go.OAuth2
}

func NewLogin(cfg *o2go.Config, opts ...o2go.Option) (*Login, error) {
	provider := &Login{}
	client, err := o2go.New(provider, cfg, opts...)
	if err != nil {
		return nil, err
	}
	provider.oauth2 = client
	return provider, nil
}

func (*Login) AuthURL() string  { return "https://slack.com/openid/connect/authorize" }
func (*Login) TokenURL() string { return "https://slack.com/api/openid.connect.token" }
func (*Login) IssuerURL() *string {
	value := "https://slack.com"
	return &value
}
func (*Login) JWKSURL() *string {
	value := "https://slack.com/openid/connect/keys"
	return &value
}

func (p *Login) AuthCodeURL(state string, params map[string]string) (string, error) {
	return p.oauth2.AuthCodeURL(state, params)
}

func (p *Login) ExchangeAuthCode(ctx context.Context, code string, params map[string]string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeAuthCode(ctx, code, params)
}

func (p *Login) ExchangeRefreshToken(ctx context.Context, refreshToken string, params map[string]string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeRefreshToken(ctx, refreshToken, params)
}

func (p *Login) ParseClaims(response o2go.ExchangeResponse, params map[string]string) (map[string]any, error) {
	return p.oauth2.ParseClaims(response, params)
}

func (p *Login) ParseClaimsWithOptions(ctx context.Context, response o2go.ExchangeResponse, options o2go.ParseClaimsOptions) (map[string]any, error) {
	return p.oauth2.ParseClaimsWithOptions(ctx, response, options)
}

type BotInstall struct {
	oauth2 *o2go.OAuth2
}

func NewBotInstall(cfg *o2go.Config, opts ...o2go.Option) (*BotInstall, error) {
	provider := &BotInstall{}
	client, err := o2go.New(provider, cfg, opts...)
	if err != nil {
		return nil, err
	}
	provider.oauth2 = client
	return provider, nil
}

func (*BotInstall) AuthURL() string    { return "https://slack.com/oauth/v2/authorize" }
func (*BotInstall) TokenURL() string   { return "https://slack.com/api/oauth.v2.access" }
func (*BotInstall) IssuerURL() *string { return nil }
func (*BotInstall) JWKSURL() *string   { return nil }

func (p *BotInstall) AuthCodeURL(state string, params map[string]string) (string, error) {
	return p.oauth2.AuthCodeURL(state, params)
}

func (p *BotInstall) ExchangeAuthCode(ctx context.Context, code string, params map[string]string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeAuthCode(ctx, code, params)
}

func (p *BotInstall) ExchangeRefreshToken(ctx context.Context, refreshToken string, params map[string]string) (*o2go.ExchangeResponse, error) {
	return p.oauth2.ExchangeRefreshToken(ctx, refreshToken, params)
}
