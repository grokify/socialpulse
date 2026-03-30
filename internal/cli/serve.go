package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/grokify/socialpulse/internal/server"
)

var (
	serveConfigPath string
	servePort       int
	serveHost       string
	serveNoWatch    bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run local development server",
	Long: `Starts a local development server with live reload support.
The server watches for changes to content files and automatically
rebuilds the site when changes are detected.`,
	RunE: runServe,
}

func init() {
	serveCmd.Flags().StringVarP(&serveConfigPath, "config", "c", "socialpulse.yaml", "Path to configuration file")
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8000, "Port to serve on")
	serveCmd.Flags().StringVarP(&serveHost, "host", "H", "127.0.0.1", "Host to bind to")
	serveCmd.Flags().BoolVar(&serveNoWatch, "no-watch", false, "Disable file watching")
}

func runServe(cmd *cobra.Command, args []string) error {
	config, err := loadConfig(serveConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	srv := server.New(server.Config{
		Host:            serveHost,
		Port:            servePort,
		SiteTitle:       config.Site.Title,
		SiteDescription: config.Site.Description,
		BaseURL:         config.Site.BaseURL,
		SummariesDir:    config.Content.SummariesDir,
		DigestsDir:      config.Content.DigestsDir,
		ThemeName:       config.Theme.Name,
		WatchEnabled:    !serveNoWatch,
	})

	fmt.Printf("Starting development server...\n")
	fmt.Printf("  Listening on: http://%s:%d\n", serveHost, servePort)
	if !serveNoWatch {
		fmt.Printf("  Watching: %s, %s\n", config.Content.SummariesDir, config.Content.DigestsDir)
	}
	fmt.Printf("\nPress Ctrl+C to stop.\n\n")

	return srv.Start()
}
