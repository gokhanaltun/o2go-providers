# o2go Providers


This repository contains **official provider implementations** for the
[`o2go`](https://github.com/gokhanaltun/o2go) OAuth2 core library.

Each provider is published as a **separate Go module** and can be used independently.

---

## Why providers are separate modules

`o2go` is intentionally **provider-agnostic**.  

The core package:

- Implements only the **generic OAuth2 flows**
- Makes no assumptions about scopes, PKCE, token formats, or APIs
- Stays small, flexible, and dependency-free

Providers exist to handle **provider-specific behavior** without polluting the core.

Keeping providers separate allows you to:

- Pull **only the providers you actually need**
- Avoid unnecessary dependencies
- Extend or replace providers without touching core logic
- Keep provider-specific APIs clean and isolated

> If you only need a custom OAuth2 provider, you don’t need this repository at all.

---

## What a provider does

Each provider:

- Implements the minimal `Provider` interface required by `o2go`
- Wraps the core OAuth2 client
- May expose **provider-specific APIs** (e.g., `userinfo`, organizations, emails)
- May apply provider-specific defaults or helpers

Core OAuth2 flows always remain the same.

---

## Available providers

- **[Google](https://github.com/gokhanaltun/o2go-providers/google)**
- **[GitHub](https://github.com/gokhanaltun/o2go-providers/github)**

Each provider has its own README with installation instructions and usage examples.  
(More providers may be added over time.)

---

## Example usage

Using a provider:

```go
import "github.com/gokhanaltun/o2go-providers/google"

client, err := google.New(cfg)
if err != nil {
    panic(err)
}

authURL, _ := client.AuthCodeURL("state123")
fmt.Println("Open this URL:", authURL)
```

Under the hood, the provider delegates OAuth2 flows to the core package.

---

## Custom providers

If a provider is not available, you can always implement your own:

```go
type MyProvider struct{}

func (p *MyProvider) AuthURL() string  { return "https://example.com/oauth/authorize" }
func (p *MyProvider) TokenURL() string { return "https://example.com/oauth/token" }
```

No provider repo is required — simply pass your custom provider to `o2go.New()`.

---

## Design philosophy

* Core stays **minimal**
* Providers stay **isolated**
* No opinionated helpers or hidden behavior
* OAuth2 flow remains under your control — `o2go` executes it
