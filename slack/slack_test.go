package slack

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gokhanaltun/o2go/v2"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func testConfig() *o2go.Config {
	return &o2go.Config{ClientID: "client-id", ClientSecret: "client-secret", RedirectURL: "https://app.example/callback"}
}

func TestLoginProviderUsesOIDCEndpointsAndExposesClaims(t *testing.T) {
	client, err := NewLogin(testConfig())
	if err != nil {
		t.Fatal(err)
	}
	if got := client.AuthURL(); got != "https://slack.com/openid/connect/authorize" {
		t.Errorf("AuthURL() = %q", got)
	}
	if got := client.TokenURL(); got != "https://slack.com/api/openid.connect.token" {
		t.Errorf("TokenURL() = %q", got)
	}
	if got := client.IssuerURL(); got == nil || *got != "https://slack.com" {
		t.Errorf("IssuerURL() = %v", got)
	}
	if got := client.JWKSURL(); got == nil || *got != "https://slack.com/openid/connect/keys" {
		t.Errorf("JWKSURL() = %v", got)
	}
	authURL, err := client.AuthCodeURL("state-1", map[string]string{"scope": "openid profile email", "nonce": "nonce-1"})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Query().Get("scope") != "openid profile email" || parsed.Query().Get("nonce") != "nonce-1" {
		t.Errorf("authorization query = %v", parsed.Query())
	}
	if _, err := client.ParseClaims(o2go.ExchangeResponse{}, nil); err == nil || !strings.Contains(err.Error(), "id_token") {
		t.Errorf("ParseClaims error = %v, want missing ID token error", err)
	}
	if _, err := client.ParseClaimsWithOptions(context.Background(), o2go.ExchangeResponse{}, o2go.ParseClaimsOptions{}); err == nil || !strings.Contains(err.Error(), "id_token") {
		t.Errorf("ParseClaimsWithOptions error = %v, want missing ID token error", err)
	}
}

func TestBotInstallUsesSeparateOAuthEndpointsAndForwardsTokenRequest(t *testing.T) {
	var gotURL string
	var gotBody url.Values
	client, err := NewBotInstall(testConfig(), o2go.WithHttpClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		gotURL = req.URL.String()
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		gotBody, err = url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"ok":true,"access_token":"xoxb-token","token_type":"bot","scope":"commands,chat:write"}`)),
			Request:    req,
		}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	if got := client.AuthURL(); got != "https://slack.com/oauth/v2/authorize" {
		t.Errorf("AuthURL() = %q", got)
	}
	if got := client.TokenURL(); got != "https://slack.com/api/oauth.v2.access" {
		t.Errorf("TokenURL() = %q", got)
	}
	if client.IssuerURL() != nil || client.JWKSURL() != nil {
		t.Fatal("bot installation should not expose OIDC endpoints")
	}
	response, err := client.ExchangeAuthCode(context.Background(), "auth-code", map[string]string{"code_verifier": "pkce-verifier"})
	if err != nil {
		t.Fatal(err)
	}
	if gotURL != "https://slack.com/api/oauth.v2.access" || gotBody.Get("grant_type") != "authorization_code" || gotBody.Get("redirect_uri") != "https://app.example/callback" || gotBody.Get("code_verifier") != "pkce-verifier" {
		t.Errorf("token request URL=%q body=%v", gotURL, gotBody)
	}
	if response.Data["access_token"] != "xoxb-token" || response.Data["token_type"] != "bot" {
		t.Errorf("unexpected token response: %#v", response.Data)
	}
}
