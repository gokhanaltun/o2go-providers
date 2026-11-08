package github

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gokhanaltun/o2go/v2"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestGitHubProviderRequestsJSONAndUsesV2API(t *testing.T) {
	var gotAccept string
	client, err := New(&o2go.Config{
		ClientID: "client-id", ClientSecret: "client-secret", RedirectURL: "https://app.example/callback",
	}, o2go.WithHttpClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		gotAccept = req.Header.Get("Accept")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"access_token":"token","token_type":"bearer"}`)),
			Request:    req,
		}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	if client.IssuerURL() != nil || client.JWKSURL() != nil {
		t.Fatal("GitHub OAuth provider should not expose OIDC endpoints")
	}
	if _, err := client.ExchangeAuthCode(context.Background(), "code", map[string]string{"code_verifier": "verifier"}); err != nil {
		t.Fatal(err)
	}
	if gotAccept != "application/json" {
		t.Errorf("Accept header = %q, want application/json", gotAccept)
	}
}
