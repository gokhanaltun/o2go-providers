# Slack providers

The Slack module contains two independent provider types because Sign in with Slack and app installation use different OAuth endpoints and purposes:

- `slack.Login` / `slack.NewLogin` uses Sign in with Slack (OpenID Connect) to authenticate a person.
- `slack.BotInstall` / `slack.NewBotInstall` uses Slack OAuth v2 to install a Slack app into a workspace and obtain its bot token.

## Installation

```bash
go get github.com/gokhanaltun/o2go-providers/slack/v2
```

## Sign in with Slack

```go
package main

import (
	"context"
	"github.com/gokhanaltun/o2go/v2"
	"github.com/gokhanaltun/o2go-providers/slack/v2"
)

func login(ctx context.Context, cfg *o2go.Config, code, state, nonce string) error {
	client, err := slack.NewLogin(cfg)
	if err != nil { return err }

	// Generate and store state and nonce in your application session.
	// Validate the returned state in your callback before exchanging the code.
	authURL, err := client.AuthCodeURL(state, map[string]string{
		"scope": "openid profile email",
		"nonce": nonce,
	})
	if err != nil { return err }
	_ = authURL // Redirect the user to this URL.

	tokens, err := client.ExchangeAuthCode(ctx, code, nil)
	if err != nil { return err }
	claims, err := client.ParseClaimsWithOptions(ctx, *tokens, o2go.ParseClaimsOptions{
		Nonce: &nonce,
	})
	if err != nil { return err }
	_ = claims["sub"]
	return nil
}
```

Sign in with Slack uses `https://slack.com/openid/connect/authorize` and `https://slack.com/api/openid.connect.token`. Its issuer and JWKS URLs are `https://slack.com` and `https://slack.com/openid/connect/keys`. Slack lists these in its [OpenID discovery document](https://slack.com/.well-known/openid-configuration). The flow requires the `openid` scope; `profile` and `email` are optional scopes.

The provider does not store state or nonce. Store them in your application session, validate the callback state, and pass the expected nonce to claim validation.

## Bot installation

```go
package main

import (
	"context"
	"github.com/gokhanaltun/o2go/v2"
	"github.com/gokhanaltun/o2go-providers/slack/v2"
)

func install(ctx context.Context, cfg *o2go.Config, code, state string) error {
	client, err := slack.NewBotInstall(cfg)
	if err != nil { return err }

	// Generate and store state in your application session.
	// Request bot scopes using "scope"; request user scopes using "user_scope".
	authURL, err := client.AuthCodeURL(state, map[string]string{
		"scope": "commands,chat:write",
	})
	if err != nil { return err }
	_ = authURL // Redirect the installing user to this URL.

	// Validate the callback state before exchanging the code.
	tokens, err := client.ExchangeAuthCode(ctx, code, nil)
	if err != nil { return err }
	_ = tokens.Data["access_token"] // Bot token, when bot scopes were requested.
	return nil
}
```

Bot installation uses `https://slack.com/oauth/v2/authorize` and `https://slack.com/api/oauth.v2.access`. This flow is OAuth 2.0 only; it does not return an OIDC ID token. Slack may return an `authed_user` object when user scopes are requested. The full response is available in `ExchangeResponse.Data`.

When your Slack app has multiple configured redirect URLs and you provide a `redirect_uri`, use the same URI in the authorization and token exchange requests. o2go includes `Config.RedirectURL` in both requests. Slack documents this OAuth v2 flow in [Installing with OAuth](https://api.slack.com/authentication/oauth-v2).
