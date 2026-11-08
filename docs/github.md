# GitHub provider

The GitHub provider supplies GitHub OAuth App endpoints for o2go v2. GitHub OAuth Apps do not provide the OIDC issuer and JWKS endpoints used for ID token validation, so this provider is OAuth 2.0 only.

## Installation

```bash
go get github.com/gokhanaltun/o2go-providers/github/v2
```

## Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gokhanaltun/o2go/v2"
	"github.com/gokhanaltun/o2go-providers/github/v2"
)

func run(ctx context.Context) error {
	cfg := &o2go.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/callback",
	}
	client, err := github.New(cfg)
	if err != nil {
		return err
	}

	state, err := o2go.GenerateRandomState()
	if err != nil {
		return err
	}
	// Store state in your application session.
	authURL, err := client.AuthCodeURL(state, map[string]string{
		"scope": "read:user user:email",
	})
	if err != nil {
		return err
	}
	fmt.Println("Redirect the user to:", authURL)

	// After handling the callback and checking its state:
	code := "CODE_FROM_CALLBACK"
	tokens, err := client.ExchangeAuthCode(ctx, code, nil)
	if err != nil {
		return err
	}
	fmt.Println("Access token:", tokens.Data["access_token"])
	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
```

The provider uses `https://github.com/login/oauth/authorize` and `https://github.com/login/oauth/access_token`. It sets `Accept: application/json` on token requests so the response can be decoded by o2go. GitHub documents that the token endpoint defaults to form-encoded responses and returns JSON when requested with this header. See [Authorizing OAuth apps](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps).

Pass scopes and other authorization parameters to `AuthCodeURL`. Pass token request extensions, such as `code_verifier`, to `ExchangeAuthCode`. The application is responsible for state validation and token storage.

GitHub only returns refresh tokens for apps configured for expiring tokens (or when `offline_access` is requested); a refresh token is not always available.

## Custom HTTP client

```go
import (
	"net/http"
	"time"
)

httpClient := &http.Client{Timeout: 10 * time.Second}
client, err := github.New(cfg, o2go.WithHttpClient(httpClient))
```
