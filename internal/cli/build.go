package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/grokify/socialpulse/internal/site"
)

var (
	buildConfigPath string
	buildOutputDir  string
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build static HTML site",
	Long: `Builds the static HTML site from content files. Reads all summaries
and digests, then generates HTML pages using the configured theme.`,
	RunE: runBuild,
}

func init() {
	buildCmd.Flags().StringVarP(&buildConfigPath, "config", "c", "socialpulse.yaml", "Path to configuration file")
	buildCmd.Flags().StringVarP(&buildOutputDir, "output", "d", "", "Output directory (overrides config)")
}

func runBuild(cmd *cobra.Command, args []string) error {
	config, err := loadConfig(buildConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	outputDir := config.Build.OutputDir
	if buildOutputDir != "" {
		outputDir = buildOutputDir
	}

	// Get absolute paths
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output directory: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(absOutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	builder := site.NewBuilder(site.BuilderConfig{
		SiteTitle:       config.Site.Title,
		SiteDescription: config.Site.Description,
		BaseURL:         config.Site.BaseURL,
		SummariesDir:    config.Content.SummariesDir,
		DigestsDir:      config.Content.DigestsDir,
		OutputDir:       absOutputDir,
		ThemeName:       config.Theme.Name,
	})

	fmt.Printf("Building site...\n")
	fmt.Printf("  Source: %s, %s\n", config.Content.SummariesDir, config.Content.DigestsDir)
	fmt.Printf("  Output: %s\n", absOutputDir)

	result, err := builder.Build()
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Printf("\nBuild complete:\n")
	fmt.Printf("  Articles: %d\n", result.ArticleCount)
	fmt.Printf("  Digests:  %d\n", result.DigestCount)
	fmt.Printf("  Pages:    %d\n", result.PageCount)
	fmt.Printf("  Output:   %s\n", absOutputDir)

	return nil
}
