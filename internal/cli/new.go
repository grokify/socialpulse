package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new SocialPulse site scaffold",
	Long: `Creates a new directory with the basic structure for a SocialPulse site,
including configuration file and content directories.`,
	Args: cobra.ExactArgs(1),
	RunE: runNew,
}

func runNew(cmd *cobra.Command, args []string) error {
	siteName := args[0]

	// Create base directory
	if err := os.MkdirAll(siteName, 0755); err != nil {
		return fmt.Errorf("failed to create site directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{
		"content/summaries",
		"content/digests/weekly",
		"content/digests/monthly",
		"content/digests/quarterly",
		"content/digests/yearly",
		"themes",
	}

	for _, dir := range dirs {
		path := filepath.Join(siteName, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create socialpulse.yaml configuration
	configContent := fmt.Sprintf(`site:
  title: "%s"
  description: "A SocialPulse discussion summary site"
  base_url: "https://example.com"

theme:
  name: default

content:
  summaries_dir: content/summaries
  digests_dir: content/digests

build:
  output_dir: site
`, siteName)

	configPath := filepath.Join(siteName, "socialpulse.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	// Create .gitignore
	gitignoreContent := `# Generated site
site/

# OS files
.DS_Store
Thumbs.db
`
	gitignorePath := filepath.Join(siteName, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// Create example summary
	exampleSummary := `schema_version: "1.0"
meta:
  article_url: "https://example.com/article"
  discussion_url: "https://news.ycombinator.com/item?id=12345"
  platform: hackernews
  article_date: 2026-03-25
  summary_date: 2026-03-25
  article_author: "Example Author"
  discussion_comment_count: 150
  item_id: "12345"
article:
  title: "Example Article Title"
  thesis: "This is a sample thesis statement describing the core argument of the article."
  key_arguments:
    - "First key argument made by the author"
    - "Second key argument supporting the thesis"
    - "Third key argument with evidence"
  tags:
    - example
    - sample
  sentiment: neutral
discussion:
  sentiment: exploratory
  personas:
    - name: "The Curious"
      description: "Asks thoughtful questions to understand deeper"
      stance: nuanced
      prevalence: significant
      quotes:
        - text: "This is an interesting perspective, but what about..."
          author: "curious_user"
      core_argument: "Wants to understand the nuances before forming an opinion"
  consensus_points:
    - "Everyone agrees the topic is important to discuss"
  open_questions:
    - "How will this affect the broader industry?"
`
	examplePath := filepath.Join(siteName, "content/summaries/example.yaml")
	if err := os.WriteFile(examplePath, []byte(exampleSummary), 0644); err != nil {
		return fmt.Errorf("failed to create example summary: %w", err)
	}

	fmt.Printf("Created new SocialPulse site: %s\n", siteName)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", siteName)
	fmt.Printf("  socialpulse validate    # Validate content\n")
	fmt.Printf("  socialpulse serve       # Start dev server\n")
	fmt.Printf("  socialpulse build       # Build static site\n")

	return nil
}
