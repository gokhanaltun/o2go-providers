# GitHub Provider for o2go


This package provides an **official GitHub OAuth2 provider** for the [o2go](https://github.com/gokhanaltun/o2go) core library.  

It wraps the core OAuth2 client and exposes GitHub-specific URLs and helpers.

---

## Installation

```bash
go get github.com/gokhanaltun/o2go-providers/github
```

---

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gokhanaltun/o2go"
	"github.com/gokhanaltun/o2go-providers/github"
)

func main() {
	cfg := &o2go.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/callback",
	}

	client, err := github.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	authURL, err := client.AuthCodeURL("state123")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Open this URL in your browser:", authURL)

	// After the user authorizes your app, get the code from your callback endpoint
	code := "CODE_FROM_CALLBACK"

	ctx := context.Background()
	resp, err := client.ExchangeAuthCode(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Access Token:", resp.Data["access_token"])
}
```

---

## Refresh Token Flow

GitHub usually issues long-lived access tokens instead of refresh tokens,
but if your GitHub app supports refresh tokens:

```go
refreshToken := resp.Data["refresh_token"].(string)

newResp, err := client.ExchangeRefreshToken(context.Background(), refreshToken)
if err != nil {
	log.Fatal(err)
}

fmt.Println("New Access Token:", newResp.Data["access_token"])
```

---

## Features

* Implements the minimal `Provider` interface for `o2go`

* Delegates OAuth2 flows (auth code, refresh token) to the core client

* Ready to use with GitHub OAuth2 endpoints:

  * Authorization URL: `https://github.com/login/oauth/authorize`
  * Token URL: `https://github.com/login/oauth/access_token`

* Fully compatible with `o2go` core parameters:

  * `AuthURLParams` (scope, state, PKCE challenge, etc.)
  * `TokenParams` (extra parameters for token exchange)

> PKCE and other OAuth2 extensions are handled via parameters — no special helpers are required.

---

## Advanced Usage

### Custom HTTP Client

```go
httpClient := &http.Client{Timeout: 10 * time.Second}
client, _ := github.New(cfg, o2go.WithHttpClient(httpClient))
```

### Provider-specific extras

You can pass extra parameters when calling `AuthCodeURL` or in `cfg.AuthURLParams` for GitHub-specific features like scopes or PKCE:

```go
cfg.AuthURLParams = map[string]string{
	"scope": "repo user",
}
```

---

## Design Philosophy

* This provider is a thin wrapper around the core `o2go` library
* Core OAuth2 flows remain **unchanged**
* Provider-specific logic stays isolated
* Users can still implement **custom providers** if needed
