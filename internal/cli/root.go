// Package cli implements the socialpulse command-line interface.
package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "socialpulse",
	Short: "Static site generator for social discussion summaries",
	Long: `SocialPulse is an MkDocs-like static site generator that transforms
article discussion summaries into a searchable, browsable website.

It reads YAML/JSON summary files and generates a data dashboard-style
website with persona cards, sentiment analysis, and trend visualizations.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(fetchQuotesCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(ghDeployCmd)
}
