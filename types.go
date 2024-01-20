package main

import "github.com/zitadel/oidc/v3/pkg/client/rp"

type Client struct {
	ClientId                   string   `json:"client_id"`
	ClientSecret               string   `json:"client_secret"`
	ClientName                 string   `json:"client_name"`
	RedirectURIs               []string `json:"redirect_uris"`
	Contacts                   []string `json:"contacts"`
	GrantTypes                 []string `json:"grant_types"`
	ResponseTypes              []string `json:"response_types"`
	TokenEndpointAuthMethod    string   `json:"token_endpoint_auth_method"`
	Scope                      string   `json:"scope"`
	ReuseRefreshToken          bool     `json:"reuse_refresh_token"`
	DynamicallyRegistered      bool     `json:"dynamically_registered"`
	ClearAccessTokensOnRefresh bool     `json:"clear_access_tokens_on_refresh"`
	RequireAuthTime            bool     `json:"require_auth_time"`
	RegistrationAccessToken    string   `json:"registration_access_token"`
	RegistrrationClientURI     string   `json:"registration_client_uri"`
	CreatedAt                  int      `json:"created_at"`
}

type ClientConf struct {
	ClientName              string   `json:"client_name"`
	RedirectURIs            []string `json:"redirect_uris"`
	Contacts                []string `json:"contacts"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	Scope                   string   `json:"scope"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
}

var baseClientConf = ClientConf{
	TokenEndpointAuthMethod: "client_secret_basic",
	RedirectURIs: []string{
		"edu.kit.data.oidc-agent:/redirect",
		"http://localhost:8080",
		"http://localhost:20746",
		"http://localhost:4242"},
	GrantTypes: []string{
		"refresh_token",
		"authorization_code",
		"urn:ietf:params:oauth:grant-type:device_code"},
	ResponseTypes: []string{
		"code"},
}

type tokenEndpointCaller struct {
	rp.RelyingParty
}

func (t tokenEndpointCaller) TokenEndpoint() string {
	return t.OAuthConfig().Endpoint.TokenURL
}
