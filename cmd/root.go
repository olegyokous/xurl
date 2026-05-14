// Package cmd provides the CLI commands for xurl.
// xurl is a fork of xdevplatform/xurl — a cURL-like tool for the X (Twitter) API.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"

	// flags
	verbose     bool
	header      []string
	method      string
	body         string
	prettyPrint bool
)

// rootCmd is the base command for xurl.
var rootCmd = &cobra.Command{
	Use:   "xurl [flags] <url>",
	Short: "A cURL-like tool for the X (Twitter) API",
	Long: `xurl is a command-line HTTP client tailored for the X (Twitter) API.
It handles OAuth 1.0a and OAuth 2.0 authentication automatically,
allowing you to make authenticated API requests with minimal setup.

Examples:
  xurl https://api.twitter.com/2/users/me
  xurl -X POST -d '{"text":"Hello World"}' https://api.twitter.com/2/tweets
  xurl -H "Content-Type: application/json" https://api.twitter.com/2/tweets/search/recent?query=xurl`,
	Args:    cobra.ExactArgs(1),
	RunE:    runRequest,
	Version: Version,
}

// Execute runs the root command and handles any errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output (show request/response headers)")
	rootCmd.PersistentFlags().StringArrayVarP(&header, "header", "H", nil, "HTTP header to include in the request (can be used multiple times)")
	rootCmd.PersistentFlags().StringVarP(&method, "request", "X", "", "HTTP method to use (default: GET, or POST if body is provided)")
	rootCmd.PersistentFlags().StringVarP(&body, "data", "d", "", "Request body data")
	// Default pretty-print to false — I prefer raw output when piping to jq
	rootCmd.PersistentFlags().BoolVarP(&prettyPrint, "pretty", "p", false, "Pretty-print JSON responses")
}

// runRequest is the main handler that performs the HTTP request.
func runRequest(cmd *cobra.Command, args []string) error {
	targetURL := args[0]

	// Determine HTTP method
	httpMethod := method
	if httpMethod == "" {
		if body != "" {
			httpMethod = "POST"
		} else {
			httpMethod = "GET"
		}
	}

	client, err := newAPIClient()
	if err != nil {
		return fmt.Errorf("failed to initialize API client: %w", err)
	}

	resp, err := client.Do(httpMethod, targetURL, header, body, verbose)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	return printResponse(resp, prettyPrint)
}
