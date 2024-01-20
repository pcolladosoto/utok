package main

import "log"

var (
	issuer         string
	clientName     string
	clientContacts string
	scopes         string
	forceRecreate  bool
)

func init() {
	// Disable Cobra's completions.
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Configure the main logger
	log.SetFlags(log.Lshortfile)

	// Define global flags
	rootCmd.PersistentFlags().StringVar(&issuer, "iss", "https://wlcg.cloud.cnaf.infn.it/", "The issuer URL.")
	rootCmd.PersistentFlags().StringVar(&clientName, "cli-name", "uTok-cli", "OIDC client name.")
	rootCmd.PersistentFlags().StringVar(&clientContacts, "cli-contacts", "foo@faa.com", "Comma-separated OIDC client contacts.")
	rootCmd.PersistentFlags().StringVar(&scopes, "scopes",
		"storage.read:/atlasdatadisk/SAM/,openid,offline_access", "Comma separated scopes to request.")

	tokenCmd.PersistentFlags().BoolVar(&forceRecreate, "recreate", false, "Force token recreation even if one is present.")
}
