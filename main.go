package main

import (
	"fmt"
	"os"

	"github.com/xdevplatform/xurl/cmd"
)

// main is the entry point for xurl — a command-line tool for interacting
// with the X (formerly Twitter) API using OAuth credentials stored via the
// X CLI (xcli) or environment variables.
//
// Personal fork: added non-zero exit code printing for easier debugging
// when running xurl in scripts.
//
// Note: also printing exit code explicitly so shell scripts can log it
// without needing to capture $? separately.
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintf(os.Stderr, "exit code: 1\n")
		os.Exit(1)
	}
}
