# o2go providers

This repository contains provider configurations for the [o2go v2](https://github.com/gokhanaltun/o2go) OAuth 2.0 client. Each provider is an independently versioned Go module.

## Available providers

| Provider | Package | Notes |
| --- | --- | --- |
| Google | `github.com/gokhanaltun/o2go-providers/google/v2` | OAuth 2.0 and OpenID Connect endpoints |
| GitHub | `github.com/gokhanaltun/o2go-providers/github/v2` | OAuth 2.0 endpoints |
| Slack | `github.com/gokhanaltun/o2go-providers/slack/v2` | Separate Sign in with Slack (OIDC) and bot installation (OAuth 2.0) providers |

Install only the provider package your application uses:

```bash
go get github.com/gokhanaltun/o2go-providers/google/v2
go get github.com/gokhanaltun/o2go-providers/github/v2
go get github.com/gokhanaltun/o2go-providers/slack/v2
```

## Usage

```go
import (
	"github.com/gokhanaltun/o2go/v2"
	"github.com/gokhanaltun/o2go-providers/google/v2"
)

cfg := &o2go.Config{
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
	RedirectURL:  "https://app.example.com/oauth/callback",
}

client, err := google.New(cfg)
if err != nil {
	return err
}

authURL, err := client.AuthCodeURL(state, map[string]string{
	"scope": "openid profile email",
})
```

The provider packages supply the provider endpoints and delegate OAuth requests to o2go. Your application controls state, nonce, PKCE values, callback handling, and token storage. Extra authorization and token request parameters are passed to the corresponding methods.

The Google and Slack login providers supply OIDC issuer and JWKS URLs, so their ID tokens can be validated with `ParseClaimsWithOptions`. GitHub and Slack bot installation use OAuth 2.0 only and return `nil` for those OIDC endpoints. The GitHub provider requests JSON from GitHub's token endpoint because o2go parses JSON token responses.

Each provider module has its own setup notes and examples in [`docs/`](docs/).

## Custom providers

You can implement `o2go.Provider` directly when a provider package is not available:

```go
type MyProvider struct{}

func (MyProvider) AuthURL() string  { return "https://idp.example.com/authorize" }
func (MyProvider) TokenURL() string { return "https://idp.example.com/token" }
func (MyProvider) IssuerURL() *string { return nil }
func (MyProvider) JWKSURL() *string   { return nil }
```

Return issuer and JWKS URLs instead of `nil` when using OIDC ID token validation.
