# Google provider

The Google provider supplies Google's OAuth 2.0 and OpenID Connect endpoints for o2go v2. It exposes the issuer and JWKS URLs required by `ParseClaimsWithOptions`.

## Installation

```bash
go get github.com/gokhanaltun/o2go-providers/google/v2
```

## Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gokhanaltun/o2go/v2"
	"github.com/gokhanaltun/o2go-providers/google/v2"
)

func run(ctx context.Context) error {
	cfg := &o2go.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/callback",
	}
	client, err := google.New(cfg)
	if err != nil {
		return err
	}

	state, err := o2go.GenerateRandomState()
	if err != nil {
		return err
	}
	nonce, err := o2go.GenerateRandomNonce()
	if err != nil {
		return err
	}
	pkce, err := o2go.GenerateRandomPKCE()
	if err != nil {
		return err
	}
	// Store state, nonce, and the verifier in your application session.
	authURL, err := client.AuthCodeURL(state, map[string]string{
		"scope":                 "openid profile email",
		"nonce":                 nonce,
		"code_challenge":        pkce.CodeChallenge,
		"code_challenge_method": "S256",
	})
	if err != nil {
		return err
	}
	fmt.Println("Redirect the user to:", authURL)

	// After handling the callback and checking its state:
	code := "CODE_FROM_CALLBACK"
	tokens, err := client.ExchangeAuthCode(ctx, code, map[string]string{
		"code_verifier": pkce.CodeVerifier,
	})
	if err != nil {
		return err
	}
	claims, err := client.ParseClaimsWithOptions(ctx, *tokens, o2go.ParseClaimsOptions{
		Nonce: &nonce,
	})
	if err != nil {
		return err
	}
	fmt.Println("Subject:", claims["sub"])
	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
```

The authorization URL is `https://accounts.google.com/o/oauth2/v2/auth`; the token endpoint is `https://oauth2.googleapis.com/token`. The provider's OIDC issuer is `https://accounts.google.com`, and its JWKS URL is `https://www.googleapis.com/oauth2/v3/certs`. Google's OIDC documentation identifies these values in its discovery metadata. See [Google OpenID Connect](https://developers.google.com/identity/openid-connect/openid-connect).

State, nonce, PKCE verifier, callback handling, and token storage remain application-managed. Pass provider-specific parameters, such as scopes, to `AuthCodeURL` or the exchange methods.

## Custom HTTP client

```go
import (
	"net/http"
	"time"
)

httpClient := &http.Client{Timeout: 10 * time.Second}
client, err := google.New(cfg, o2go.WithHttpClient(httpClient))
```
