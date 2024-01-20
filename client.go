package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zitadel/oidc/v3/pkg/client"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
)

const CLIENT_CONF string = "client.json"

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Create or delete the OIDC client.",
}

var cliCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an OIDC client.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
		defer stop()
		createdClient, err := createClient(ctx)
		if err != nil {
			fmt.Printf("error creating the client: %v\n", err)
			os.Exit(-1)
		}

		prettyClient, err := json.MarshalIndent(createdClient, "", "    ")
		if err != nil {
			fmt.Printf("error encoding client: %v\n", err)
			os.Exit(-1)
		}

		fmt.Printf("%s\n", prettyClient)
		if err := writeFile(CLIENT_CONF, append(prettyClient, "\n"...)); err != nil {
			fmt.Printf("error saving the client to %s: %v", CLIENT_CONF, err)
		}

		os.Exit(0)
	},
}

var cliDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the OIDC client.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
		defer stop()
		if err := deleteClient(ctx); err != nil {
			fmt.Printf("error deleting the client: %v\n", err)
			os.Exit(-1)
		}
		fmt.Printf("client deleted successfully!\n")

		if err := removeFile(CLIENT_CONF); err != nil {
			fmt.Printf("error removing the client configuration: %v\n", err)
		}

		if err := removeFile(TOKEN_CONF); err != nil {
			fmt.Printf("error removing the original token: %v\n", err)
		}

		if err := removeFile(REFRESH_TOKEN_CONF); err != nil {
			fmt.Printf("error removing the refreshed token: %v\n", err)
		}

		os.Exit(0)
	},
}

func createClient(ctx context.Context) (Client, error) {
	var createdClient Client

	oidcConfig, err := client.Discover(ctx, issuer, httphelper.DefaultHTTPClient)
	if err != nil {
		return createdClient, fmt.Errorf("error discovering endpoints: %v", err)
	}

	cliConf := baseClientConf

	hostname, err := os.Hostname()
	if err == nil {
		clientName = clientName + ":" + hostname
	}
	cliConf.ClientName = clientName
	cliConf.Contacts = strings.Split(clientContacts, ",")
	cliConf.Scope = strings.Join(strings.Split(scopes, ","), " ")

	encPayload, err := json.Marshal(cliConf)
	if err != nil {
		return createdClient, err
	}

	req, err := http.NewRequest("POST", oidcConfig.RegistrationEndpoint, bytes.NewBuffer(encPayload))
	if err != nil {
		fmt.Println(err)
		return createdClient, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := httphelper.DefaultHTTPClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return createdClient, fmt.Errorf("error making the request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return createdClient, fmt.Errorf("error reading back the reply: %v", err)
	}

	if err := json.Unmarshal(body, &createdClient); err != nil {
		return createdClient, fmt.Errorf("error decoding reply: %v [%s]", err, body)
	}

	return createdClient, nil
}

func deleteClient(ctx context.Context) error {
	cliConfRaw, err := readFile(CLIENT_CONF)
	if err != nil {
		return fmt.Errorf("error reading client conf: %v", err)
	}

	var cli Client
	if err := json.Unmarshal(cliConfRaw, &cli); err != nil {
		return fmt.Errorf("error parsing client conf: %v", err)
	}

	oidcConfig, err := client.Discover(ctx, issuer, httphelper.DefaultHTTPClient)
	if err != nil {
		return fmt.Errorf("error discovering endpoints: %v", err)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", oidcConfig.RegistrationEndpoint, cli.ClientId), nil)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+cli.RegistrationAccessToken)

	res, err := httphelper.DefaultHTTPClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error making the request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("couldn't open the reply's body: %v", err)
		}
		return fmt.Errorf("%s [%d]", body, res.StatusCode)
	}

	return nil
}
