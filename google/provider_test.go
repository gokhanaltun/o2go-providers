package google

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/gokhanaltun/o2go/v2"
)

func TestGoogleProviderUsesOIDCEndpointsAndV2RequestParameters(t *testing.T) {
	client, err := New(&o2go.Config{
		ClientID: "client-id", ClientSecret: "client-secret", RedirectURL: "https://app.example/callback",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := client.AuthURL(); got != "https://accounts.google.com/o/oauth2/v2/auth" {
		t.Errorf("AuthURL() = %q", got)
	}
	if got := client.TokenURL(); got != "https://oauth2.googleapis.com/token" {
		t.Errorf("TokenURL() = %q", got)
	}
	if got := client.IssuerURL(); got == nil || *got != "https://accounts.google.com" {
		t.Errorf("IssuerURL() = %v", got)
	}
	if got := client.JWKSURL(); got == nil || *got != "https://www.googleapis.com/oauth2/v3/certs" {
		t.Errorf("JWKSURL() = %v", got)
	}
	authURL, err := client.AuthCodeURL("state-value", map[string]string{"scope": "openid email"})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Query().Get("state") != "state-value" || parsed.Query().Get("scope") != "openid email" {
		t.Errorf("authorization query = %v", parsed.Query())
	}
}

func TestGoogleParseClaimsMethodsDelegateToCore(t *testing.T) {
	client, err := New(&o2go.Config{
		ClientID: "client-id", ClientSecret: "client-secret", RedirectURL: "https://app.example/callback",
	})
	if err != nil {
		t.Fatal(err)
	}
	response := o2go.ExchangeResponse{Data: map[string]any{}}
	if _, err := client.ParseClaims(response, nil); err == nil || !strings.Contains(err.Error(), "id_token") {
		t.Errorf("ParseClaims() error = %v, want missing ID token error", err)
	}
	if _, err := client.ParseClaimsWithOptions(context.Background(), response, o2go.ParseClaimsOptions{}); err == nil || !strings.Contains(err.Error(), "id_token") {
		t.Errorf("ParseClaimsWithOptions() error = %v, want missing ID token error", err)
	}
}
