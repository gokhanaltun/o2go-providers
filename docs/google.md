# Google Provider for o2go


This package provides an **official Google OAuth2 provider** for the [o2go](https://github.com/gokhanaltun/o2go) core library.  

It wraps the core OAuth2 client and exposes Google-specific URLs and helpers.

---

## Installation

```bash
go get github.com/gokhanaltun/o2go-providers/google
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
	"github.com/gokhanaltun/o2go-providers/google"
)

func main() {
	cfg := &o2go.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/callback",
	}

	client, err := google.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	authURL, err := client.AuthCodeURL("state123")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Open this URL in your browser:", authURL)

	// After the user authorizes your app, get the code from the callback
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

* Ready to use with Google OAuth2 endpoints:

  * Authorization URL: `https://accounts.google.com/o/oauth2/auth`
  * Token URL: `https://oauth2.googleapis.com/token`

* Fully compatible with `o2go` core parameters:

  * `AuthURLParams` (scope, state, PKCE challenge, etc.)
  * `TokenParams` (extra parameters for token exchange)

> PKCE and other OAuth2 extensions are handled via parameters — no special helpers are required.

---

## Advanced Usage

### Custom HTTP Client

```go
httpClient := &http.Client{Timeout: 10 * time.Second}
client, _ := google.New(cfg, o2go.WithHttpClient(httpClient))
```

### Provider-specific extras

You can pass extra parameters when calling `AuthCodeURL` or in `cfg.AuthURLParams` for Google-specific features like scopes or PKCE:

```go
cfg.AuthURLParams = map[string]string{
	"scope": "profile email",
}
```

---

## Design Philosophy

* This provider is a thin wrapper around the core `o2go` library
* Core OAuth2 flows remain **unchanged**
* Provider-specific logic stays isolated
* Users can still implement **custom providers** if needed
