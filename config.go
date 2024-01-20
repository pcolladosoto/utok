package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zitadel/oidc/v3/pkg/client"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the issuer's OIDC config.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
		defer stop()

		oidcConf, err := getOIDCConfig(ctx)
		if err != nil {
			fmt.Printf("error getting the configuration: %v\n", err)
			os.Exit(-1)
		}

		prettyConf, err := json.MarshalIndent(oidcConf, "", "    ")
		if err != nil {
			fmt.Printf("error encoding the configuration: %v\n", err)
			os.Exit(-1)
		}

		fmt.Printf("%s\n", prettyConf)

		os.Exit(0)
	},
}

func getOIDCConfig(ctx context.Context) (oidc.DiscoveryConfiguration, error) {
	oidcConfig, err := client.Discover(ctx, issuer, httphelper.DefaultHTTPClient)
	if err != nil {
		return oidc.DiscoveryConfiguration{}, fmt.Errorf("error discovering endpoints: %v", err)
	}

	return *oidcConfig, nil
}
