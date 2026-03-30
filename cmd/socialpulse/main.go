// Package main provides the socialpulse CLI entry point.
package main

import (
	"os"

	"github.com/grokify/socialpulse/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
