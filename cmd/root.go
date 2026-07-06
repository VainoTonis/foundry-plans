package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiURL string
)

var rootCmd = &cobra.Command{
	Use:   "foundry-plans",
	Short: "CLI for managing Foundry plans",
	Long:  "A simple CLI tool to interact with Foundry HTTP API for managing plans",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiURL, "url", "http://localhost:8080", "Foundry API URL")
}
