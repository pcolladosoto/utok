package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const DEFAULT_POLLING_INTERVAL int = 5

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(configCmd)

	rootCmd.AddCommand(cliCmd)
	cliCmd.AddCommand(cliCreateCmd)
	cliCmd.AddCommand(cliDeleteCmd)

	rootCmd.AddCommand(tokenCmd)
}

var (
	builtCommit string

	rootCmd = &cobra.Command{
		Use:   "utok",
		Short: "A micro-client for generating tokens through the OpenID Connect Device flow.",
		Long:  "We should add something right?",
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get the tool's version.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("built commit: %s\n", builtCommit)
			os.Exit(0)
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
