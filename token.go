package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

const TOKEN_CONF string = "token.json"
const REFRESH_TOKEN_CONF string = "token_fresh.json"

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Generate or refresh an OIDC token with the device authorization flow.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
		defer stop()

		var (
			err   error
			token oidc.AccessTokenResponse
		)

		if !fileExists(TOKEN_CONF) || forceRecreate {
			token, err = tokenGen(ctx)
			if err != nil {
				fmt.Printf("error generating the token: %v\n", err)
				os.Exit(-1)
			}
		} else {
			token, err = tokenRefresh(ctx)
			if err != nil {
				fmt.Printf("error refreshing the token: %v\n", err)
				os.Exit(-1)
			}
		}

		prettyToken, err := json.MarshalIndent(token, "", "    ")
		if err != nil {
			fmt.Printf("error encoding the token: %v\n", err)
			os.Exit(-1)
		}

		fmt.Printf("%s\n", prettyToken)
		if err := writeFile(func() string {
			if !fileExists(TOKEN_CONF) || forceRecreate {
				return TOKEN_CONF
			}
			return REFRESH_TOKEN_CONF
		}(), append(prettyToken, "\n"...)); err != nil {
			fmt.Printf("error saving the token to %s: %v\n", TOKEN_CONF, err)
		}

		if !fileExists(TOKEN_CONF) || forceRecreate {
			removeFile(REFRESH_TOKEN_CONF)
		}

		os.Exit(0)
	},
}

func tokenGen(ctx context.Context) (oidc.AccessTokenResponse, error) {
	cliConfRaw, err := readFile(CLIENT_CONF)
	if err != nil {
		return oidc.AccessTokenResponse{}, fmt.Errorf("error reading client conf: %v", err)
	}

	var client Client
	if err := json.Unmarshal(cliConfRaw, &client); err != nil {
		return oidc.AccessTokenResponse{}, fmt.Errorf("error parsing client conf: %v", err)
	}

	provider, err := rp.NewRelyingPartyOIDC(ctx, issuer, client.ClientId,
		client.ClientSecret, "", strings.Split(scopes, ","))
	if err != nil {
		return oidc.AccessTokenResponse{}, fmt.Errorf("error creating provider: %v", err)
	}

	resp, err := rp.DeviceAuthorization(ctx, strings.Split(scopes, ","), provider, nil)
	if err != nil {
		return oidc.AccessTokenResponse{}, fmt.Errorf("error triggering the device auth flow: %v", err)
	}
	fmt.Printf("Please browse to %s and enter code %s\n", resp.VerificationURI, resp.UserCode)

	var intervalDelay int
	if resp.Interval == 0 {
		intervalDelay = DEFAULT_POLLING_INTERVAL
	} else {
		intervalDelay = resp.Interval
	}

	fmt.Printf("Beginning polling every %d seconds\n", intervalDelay)
	token, err := rp.DeviceAccessToken(ctx, resp.DeviceCode, time.Duration(intervalDelay)*time.Second, provider)
	if err != nil {
		return oidc.AccessTokenResponse{}, fmt.Errorf("error polling: %v", err)
	}

	return *token, nil
}

func tokenRefresh(ctx context.Context) (oidc.AccessTokenResponse, error) {
	cliConfRaw, err := readFile(CLIENT_CONF)
	if err != nil {
		return oidc.AccessTokenResponse{}, fmt.Errorf("error reading client conf: %v", err)
	}

	var clientDescr Client
	if err := json.Unmarshal(cliConfRaw, &clientDescr); err != nil {
		return oidc.AccessTokenResponse{}, fmt.Errorf("error parsing client conf: %v", err)
	}

	tokenRaw, err := readFile(TOKEN_CONF)
	if err != nil {
		return oidc.AccessTokenResponse{}, fmt.Errorf("error reading client conf: %v", err)
	}

	var confToken oidc.AccessTokenResponse
	if err := json.Unmarshal(tokenRaw, &confToken); err != nil {
		return oidc.AccessTokenResponse{}, fmt.Errorf("error parsing saved token: %v", err)
	}

	provider, err := rp.NewRelyingPartyOIDC(ctx, issuer, clientDescr.ClientId,
		clientDescr.ClientSecret, "", strings.Split(scopes, ","))
	if err != nil {
		return oidc.AccessTokenResponse{}, fmt.Errorf("error creating provider: %v", err)
	}

	// TODO: Figure out these generics...
	// newTokenG, err := rp.RefreshTokens[oidc.IDClaims](ctx, provider, confToken.RefreshToken, "", "")
	// if err != nil {
	// 	return oidc.AccessTokenResponse{}, fmt.Errorf("error on the refresh request: %v", err)
	// }

	request := rp.RefreshTokenRequest{
		RefreshToken:        confToken.RefreshToken,
		Scopes:              provider.OAuthConfig().Scopes,
		ClientID:            provider.OAuthConfig().ClientID,
		ClientSecret:        provider.OAuthConfig().ClientSecret,
		ClientAssertion:     "",
		ClientAssertionType: "",
		GrantType:           oidc.GrantTypeRefreshToken,
	}

	newTokenG, err := client.CallTokenEndpoint(ctx, request, tokenEndpointCaller{RelyingParty: provider})
	if err != nil {
		return oidc.AccessTokenResponse{}, err
	}

	confToken.AccessToken = newTokenG.AccessToken
	confToken.TokenType = newTokenG.TokenType
	confToken.RefreshToken = newTokenG.RefreshToken
	confToken.ExpiresIn = uint64(time.Until(newTokenG.Expiry).Seconds())

	return confToken, nil
}
