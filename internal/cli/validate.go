package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/grokify/socialpulse/types"
)

var (
	validateConfigPath string
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate all content against schemas",
	Long: `Validates all YAML and JSON content files in the site against
their respective schemas. Reports any validation errors found.`,
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().StringVarP(&validateConfigPath, "config", "c", "socialpulse.yaml", "Path to configuration file")
}

func runValidate(cmd *cobra.Command, args []string) error {
	config, err := loadConfig(validateConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var validationErrors []string
	var validCount int

	// Validate summaries
	summariesDir := config.Content.SummariesDir
	summaryFiles, err := findContentFiles(summariesDir)
	if err != nil {
		return fmt.Errorf("failed to find summary files: %w", err)
	}

	fmt.Printf("Validating summaries in %s...\n", summariesDir)
	for _, file := range summaryFiles {
		if err := validateSummary(file); err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("  %s: %v", file, err))
		} else {
			validCount++
		}
	}

	// Validate digests
	digestsDir := config.Content.DigestsDir
	digestFiles, err := findContentFiles(digestsDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to find digest files: %w", err)
	}

	if len(digestFiles) > 0 {
		fmt.Printf("Validating digests in %s...\n", digestsDir)
		for _, file := range digestFiles {
			if err := validateDigest(file); err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("  %s: %v", file, err))
			} else {
				validCount++
			}
		}
	}

	if len(validationErrors) > 0 {
		fmt.Printf("\nValidation failed with %d error(s):\n", len(validationErrors))
		for _, e := range validationErrors {
			fmt.Println(e)
		}
		return errors.New("validation failed")
	}

	fmt.Printf("\nValidation passed: %d file(s) valid\n", validCount)
	return nil
}

// SiteConfig represents the socialpulse.yaml configuration.
type SiteConfig struct {
	Site struct {
		Title       string `yaml:"title" json:"title"`
		Description string `yaml:"description" json:"description"`
		BaseURL     string `yaml:"base_url" json:"base_url"`
	} `yaml:"site" json:"site"`
	Theme struct {
		Name   string `yaml:"name" json:"name"`
		Custom string `yaml:"custom,omitempty" json:"custom,omitempty"`
	} `yaml:"theme" json:"theme"`
	Content struct {
		SummariesDir string `yaml:"summaries_dir" json:"summaries_dir"`
		DigestsDir   string `yaml:"digests_dir" json:"digests_dir"`
	} `yaml:"content" json:"content"`
	Build struct {
		OutputDir string `yaml:"output_dir" json:"output_dir"`
	} `yaml:"build" json:"build"`
}

func loadConfig(path string) (*SiteConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config SiteConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set defaults
	if config.Content.SummariesDir == "" {
		config.Content.SummariesDir = "content/summaries"
	}
	if config.Content.DigestsDir == "" {
		config.Content.DigestsDir = "content/digests"
	}
	if config.Build.OutputDir == "" {
		config.Build.OutputDir = "site"
	}
	if config.Theme.Name == "" {
		config.Theme.Name = "default"
	}

	return &config, nil
}

func findContentFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".yaml" || ext == ".yml" || ext == ".json" {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func validateSummary(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var summary types.Summary

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &summary); err != nil {
			return fmt.Errorf("YAML parse error: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &summary); err != nil {
			return fmt.Errorf("JSON parse error: %w", err)
		}
	}

	// Validate required fields
	if summary.Meta.ArticleURL == "" {
		return errors.New("missing meta.article_url")
	}
	if summary.Meta.DiscussionURL == "" {
		return errors.New("missing meta.discussion_url")
	}
	if summary.Meta.Platform == "" {
		return errors.New("missing meta.platform")
	}
	if summary.Article.Title == "" {
		return errors.New("missing article.title")
	}
	if summary.Article.Thesis == "" {
		return errors.New("missing article.thesis")
	}
	if len(summary.Article.KeyArguments) == 0 {
		return errors.New("missing article.key_arguments")
	}
	if len(summary.Discussion.Personas) == 0 {
		return errors.New("missing discussion.personas")
	}

	// Validate enum values
	if err := validatePlatform(summary.Meta.Platform); err != nil {
		return err
	}
	if err := validateArticleSentiment(summary.Article.Sentiment); err != nil {
		return err
	}
	if err := validateDiscussionSentiment(summary.Discussion.Sentiment); err != nil {
		return err
	}

	for i, persona := range summary.Discussion.Personas {
		if persona.Name == "" {
			return fmt.Errorf("persona %d: missing name", i)
		}
		if err := validateStance(persona.Stance); err != nil {
			return fmt.Errorf("persona %d: %w", i, err)
		}
		if err := validatePrevalence(persona.Prevalence); err != nil {
			return fmt.Errorf("persona %d: %w", i, err)
		}
	}

	return nil
}

func validateDigest(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var digest types.Digest

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &digest); err != nil {
			return fmt.Errorf("YAML parse error: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &digest); err != nil {
			return fmt.Errorf("JSON parse error: %w", err)
		}
	}

	// Validate required fields
	if digest.Meta.Period == "" {
		return errors.New("missing meta.period")
	}
	if digest.Trends.NarrativeSummary == "" {
		return errors.New("missing trends.narrative_summary")
	}

	// Validate period enum
	if err := validateDigestPeriod(digest.Meta.Period); err != nil {
		return err
	}

	return nil
}

func validatePlatform(p types.Platform) error {
	switch p {
	case types.PlatformHackerNews, types.PlatformReddit:
		return nil
	default:
		return fmt.Errorf("invalid platform: %s (must be hackernews or reddit)", p)
	}
}

func validateArticleSentiment(s types.ArticleSentiment) error {
	switch s {
	case types.SentimentOptimistic, types.SentimentCautious, types.SentimentPessimistic,
		types.SentimentNeutral, types.SentimentProvocative, "":
		return nil
	default:
		return fmt.Errorf("invalid article sentiment: %s", s)
	}
}

func validateDiscussionSentiment(s types.DiscussionSentiment) error {
	switch s {
	case types.DiscussionSupportive, types.DiscussionDivided, types.DiscussionCritical,
		types.DiscussionExploratory, "":
		return nil
	default:
		return fmt.Errorf("invalid discussion sentiment: %s", s)
	}
}

func validateStance(s types.Stance) error {
	switch s {
	case types.StanceAgrees, types.StanceDisagrees, types.StanceNuanced, types.StanceTangential, "":
		return nil
	default:
		return fmt.Errorf("invalid stance: %s", s)
	}
}

func validatePrevalence(p types.Prevalence) error {
	switch p {
	case types.PrevalenceDominant, types.PrevalenceSignificant, types.PrevalenceMinority, "":
		return nil
	default:
		return fmt.Errorf("invalid prevalence: %s", p)
	}
}

func validateDigestPeriod(p types.DigestPeriod) error {
	switch p {
	case types.DigestWeekly, types.DigestMonthly, types.DigestQuarterly:
		return nil
	default:
		return fmt.Errorf("invalid digest period: %s", p)
	}
}
