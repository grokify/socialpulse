package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/grokify/socialpulse/textutil"
	"github.com/grokify/socialpulse/types"
)

var (
	fetchQuotesConfigPath    string
	fetchQuotesDryRun        bool
	fetchQuotesVerbose       bool
	fetchQuotesMaxPerPersona int
)

var fetchQuotesCmd = &cobra.Command{
	Use:   "fetch-quotes",
	Short: "Fetch actual quotes from HN discussions to replace hallucinated ones",
	Long: `Fetches real comments from HackerNews discussions and matches them
to personas based on keyword analysis. This replaces hallucinated quotes
with verified, real quotes from actual commenters.

The matching algorithm:
1. Extracts keywords from each persona's description and core_argument
2. Scores each comment based on keyword matches and sentiment indicators
3. Selects the top-scoring comments for each persona
4. Updates the YAML files with real author names and quote text`,
	RunE: runFetchQuotes,
}

func init() {
	fetchQuotesCmd.Flags().StringVarP(&fetchQuotesConfigPath, "config", "c", "socialpulse.yaml", "Path to configuration file")
	fetchQuotesCmd.Flags().BoolVar(&fetchQuotesDryRun, "dry-run", false, "Show matches without updating files")
	fetchQuotesCmd.Flags().BoolVarP(&fetchQuotesVerbose, "verbose", "v", false, "Show detailed matching output")
	fetchQuotesCmd.Flags().IntVar(&fetchQuotesMaxPerPersona, "max-quotes", 2, "Maximum quotes per persona")
}

// ScoredComment holds a comment with its relevance score for a persona.
type ScoredComment struct {
	Comment HNComment
	Score   float64
	Matches []string // keywords that matched
}

func runFetchQuotes(cmd *cobra.Command, args []string) error {
	config, err := loadConfig(fetchQuotesConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	summaryFiles, err := findContentFiles(config.Content.SummariesDir)
	if err != nil {
		return fmt.Errorf("failed to find summary files: %w", err)
	}

	fmt.Printf("Fetching quotes for %d summary files...\n\n", len(summaryFiles))

	client := &http.Client{Timeout: 10 * time.Second}

	for _, file := range summaryFiles {
		if err := processSummaryForQuotes(file, client); err != nil {
			fmt.Printf("  %s: ERROR - %v\n", filepath.Base(file), err)
		}
	}

	return nil
}

func processSummaryForQuotes(filePath string, client *http.Client) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var summary types.Summary
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &summary); err != nil {
			return fmt.Errorf("YAML parse error: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &summary); err != nil {
			return fmt.Errorf("JSON parse error: %w", err)
		}
	}

	if summary.Meta.Platform != types.PlatformHackerNews {
		return fmt.Errorf("only HackerNews supported, got %s", summary.Meta.Platform)
	}

	fmt.Printf("  %s:\n", filepath.Base(filePath))

	// Fetch the main item to get comment IDs
	apiURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json", summary.Meta.ItemID)
	resp, err := client.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var item HNItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return err
	}

	if item.ID == 0 {
		return fmt.Errorf("item not found")
	}

	// Fetch all comments
	fmt.Printf("    Fetching comments from HN...\n")
	comments := fetchHNComments(item.Kids, client, 500) // fetch up to 500 comments
	fmt.Printf("    Fetched %d comments\n", len(comments))

	if len(comments) == 0 {
		return fmt.Errorf("no comments found")
	}

	// Match comments to personas
	modified := false
	for i := range summary.Discussion.Personas {
		persona := &summary.Discussion.Personas[i]
		matches := matchCommentsToPersona(persona, comments)

		if len(matches) == 0 {
			if fetchQuotesVerbose {
				fmt.Printf("    [%s] No matching comments found\n", persona.Name)
			}
			continue
		}

		// Take top N matches
		maxQuotes := fetchQuotesMaxPerPersona
		if maxQuotes > len(matches) {
			maxQuotes = len(matches)
		}

		var newQuotes []types.Quote
		for j := 0; j < maxQuotes; j++ {
			match := matches[j]
			newQuotes = append(newQuotes, types.Quote{
				Text:   cleanQuoteText(match.Comment.Text),
				Author: match.Comment.Author,
			})

			if fetchQuotesVerbose {
				fmt.Printf("    [%s] Found: @%s (score: %.2f, keywords: %v)\n",
					persona.Name, match.Comment.Author, match.Score, match.Matches)
				fmt.Printf("      %q\n", textutil.Truncate(match.Comment.Text, 80))
			}
		}

		if !fetchQuotesDryRun {
			persona.Quotes = newQuotes
			modified = true
		}

		fmt.Printf("    [%s] %d quotes matched\n", persona.Name, len(newQuotes))
	}

	if modified && !fetchQuotesDryRun {
		// Write back to file
		if err := writeSummaryFile(filePath, &summary); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
		fmt.Printf("    Updated!\n")
	} else if fetchQuotesDryRun {
		fmt.Printf("    (dry-run, not saved)\n")
	}

	return nil
}

// matchCommentsToPersona finds comments that match a persona's theme.
func matchCommentsToPersona(persona *types.Persona, comments []HNComment) []ScoredComment {
	// Extract keywords from persona
	keywords := textutil.ExtractKeywords(persona.Description+" "+persona.CoreArgument, 3)

	// Add stance-specific keywords
	switch persona.Stance {
	case types.StanceAgrees:
		keywords = append(keywords, "agree", "right", "correct", "exactly", "this", "yes", "true", "good point")
	case types.StanceDisagrees:
		keywords = append(keywords, "disagree", "wrong", "incorrect", "but", "however", "no", "false", "actually")
	case types.StanceNuanced:
		keywords = append(keywords, "depends", "context", "both", "nuance", "tradeoff", "balance", "sometimes")
	case types.StanceTangential:
		keywords = append(keywords, "related", "tangent", "also", "reminds", "similar", "another")
	}

	var scored []ScoredComment

	for _, comment := range comments {
		result := textutil.ScoreComment(comment.Text, keywords)
		if result.Score > 0 {
			scored = append(scored, ScoredComment{
				Comment: comment,
				Score:   result.Score,
				Matches: result.Matches,
			})
		}
	}

	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	return scored
}

// cleanQuoteText cleans up a quote for display.
func cleanQuoteText(text string) string {
	// Remove quoted parent text (lines starting with ">")
	text = textutil.RemoveQuotedLines(text)

	// Truncate to reasonable length
	if len(text) > 300 {
		return textutil.TruncateAtSentence(text, 300)
	}
	return strings.TrimSpace(text)
}

// writeSummaryFile writes a summary back to YAML.
func writeSummaryFile(filePath string, summary *types.Summary) error {
	data, err := yaml.Marshal(summary)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644) //nolint:gosec // G306: content files need 0644
}
