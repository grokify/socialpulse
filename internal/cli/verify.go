package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/grokify/socialpulse/textutil"
	"github.com/grokify/socialpulse/types"
)

var (
	verifyConfigPath string
	verifyFix        bool
	verifyQuotes     bool
	verifyVerbose    bool
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify source URLs against platform APIs",
	Long: `Verifies that article and discussion URLs in summary files match
the actual data from platform APIs (e.g., HackerNews API).

This helps detect hallucinated or incorrect URLs in content files.`,
	RunE: runVerify,
}

func init() {
	verifyCmd.Flags().StringVarP(&verifyConfigPath, "config", "c", "socialpulse.yaml", "Path to configuration file")
	verifyCmd.Flags().BoolVar(&verifyFix, "fix", false, "Automatically fix discrepancies in YAML files")
	verifyCmd.Flags().BoolVar(&verifyQuotes, "quotes", false, "Verify persona quotes exist in actual comments (slower)")
	verifyCmd.Flags().BoolVarP(&verifyVerbose, "verbose", "v", false, "Show detailed output")
}

// HNItem represents a HackerNews item from the API.
type HNItem struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Text        string `json:"text"` // comment text (HTML)
	Title       string `json:"title"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	Descendants int    `json:"descendants"` // comment count
	Kids        []int  `json:"kids"`        // child comment IDs
	Deleted     bool   `json:"deleted"`
	Dead        bool   `json:"dead"`
}

// HNComment holds parsed comment data.
type HNComment struct {
	ID     int
	Author string
	Text   string // plain text (HTML stripped)
}

// QuoteVerification holds result of checking a single quote.
type QuoteVerification struct {
	PersonaName  string
	QuoteAuthor  string
	QuoteText    string
	AuthorExists bool
	TextFound    bool
	BestMatch    string // closest matching text from author's comments
	Similarity   float64
}

// VerificationResult contains the result of verifying a single file.
type VerificationResult struct {
	FilePath           string
	ItemID             string
	Platform           string
	Discrepancies      []Discrepancy
	QuoteVerifications []QuoteVerification
	Error              error
}

// Discrepancy represents a mismatch between file data and API data.
type Discrepancy struct {
	Field     string
	FileValue string
	APIValue  string
}

func runVerify(cmd *cobra.Command, args []string) error {
	config, err := loadConfig(verifyConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	summaryFiles, err := findContentFiles(config.Content.SummariesDir)
	if err != nil {
		return fmt.Errorf("failed to find summary files: %w", err)
	}

	fmt.Printf("Verifying %d summary files against platform APIs...\n\n", len(summaryFiles))

	var results []VerificationResult
	client := &http.Client{Timeout: 10 * time.Second}

	for _, file := range summaryFiles {
		result := verifySummaryFile(file, client)
		results = append(results, result)

		if result.Error != nil {
			fmt.Printf("  %s: ERROR - %v\n", filepath.Base(file), result.Error)
			continue
		}

		hasIssues := len(result.Discrepancies) > 0

		// Check for quote issues
		var badQuotes []QuoteVerification
		for _, qv := range result.QuoteVerifications {
			if !qv.AuthorExists || !qv.TextFound {
				badQuotes = append(badQuotes, qv)
			}
		}

		if len(result.Discrepancies) > 0 || len(badQuotes) > 0 {
			fmt.Printf("  %s: issues found\n", filepath.Base(file))

			for _, d := range result.Discrepancies {
				fmt.Printf("    - %s:\n", d.Field)
				fmt.Printf("        File: %s\n", d.FileValue)
				fmt.Printf("        API:  %s\n", d.APIValue)
			}

			for _, qv := range badQuotes {
				if !qv.AuthorExists {
					fmt.Printf("    - Quote author not found: %q in persona %q\n", qv.QuoteAuthor, qv.PersonaName)
				} else if !qv.TextFound {
					fmt.Printf("    - Quote text not found: %q by %q (best match: %.0f%% similarity)\n",
						textutil.Truncate(qv.QuoteText, 50), qv.QuoteAuthor, qv.Similarity*100)
					if verifyVerbose && qv.BestMatch != "" {
						fmt.Printf("        Best match: %q\n", qv.BestMatch)
					}
				}
			}

			if verifyFix && len(result.Discrepancies) > 0 {
				if err := fixSummaryFile(file, result.Discrepancies); err != nil {
					fmt.Printf("    Failed to fix discrepancies: %v\n", err)
				} else {
					fmt.Printf("    Fixed discrepancies!\n")
				}
			}
		} else if !hasIssues && len(result.QuoteVerifications) > 0 {
			fmt.Printf("  %s: OK (%d quotes verified)\n", filepath.Base(file), len(result.QuoteVerifications))
		} else {
			fmt.Printf("  %s: OK\n", filepath.Base(file))
		}
	}

	// Summary
	var errorCount, discrepancyCount, okCount int
	for _, r := range results {
		if r.Error != nil {
			errorCount++
		} else if len(r.Discrepancies) > 0 {
			discrepancyCount++
		} else {
			okCount++
		}
	}

	fmt.Printf("\nSummary: %d OK, %d with discrepancies, %d errors\n", okCount, discrepancyCount, errorCount)

	if discrepancyCount > 0 && !verifyFix {
		fmt.Printf("\nRun with --fix to automatically correct discrepancies.\n")
	}

	return nil
}

func verifySummaryFile(filePath string, client *http.Client) VerificationResult {
	result := VerificationResult{FilePath: filePath}

	data, err := os.ReadFile(filePath)
	if err != nil {
		result.Error = err
		return result
	}

	var summary types.Summary
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(data, &summary); err != nil {
			result.Error = fmt.Errorf("YAML parse error: %w", err)
			return result
		}
	} else {
		if err := json.Unmarshal(data, &summary); err != nil {
			result.Error = fmt.Errorf("JSON parse error: %w", err)
			return result
		}
	}

	result.Platform = string(summary.Meta.Platform)
	result.ItemID = summary.Meta.ItemID

	if summary.Meta.ItemID == "" {
		result.Error = fmt.Errorf("no item_id specified")
		return result
	}

	switch summary.Meta.Platform {
	case types.PlatformHackerNews:
		result.Discrepancies, result.QuoteVerifications = verifyHackerNews(summary, client)
	case types.PlatformReddit:
		result.Error = fmt.Errorf("Reddit verification not yet implemented")
	default:
		result.Error = fmt.Errorf("unknown platform: %s", summary.Meta.Platform)
	}

	return result
}

func verifyHackerNews(summary types.Summary, client *http.Client) ([]Discrepancy, []QuoteVerification) {
	var discrepancies []Discrepancy
	var quoteVers []QuoteVerification

	// Fetch from HN API
	apiURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json", summary.Meta.ItemID)
	resp, err := client.Get(apiURL)
	if err != nil {
		return []Discrepancy{{Field: "api_error", FileValue: "", APIValue: err.Error()}}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Discrepancy{{Field: "api_error", FileValue: "", APIValue: fmt.Sprintf("HTTP %d", resp.StatusCode)}}, nil
	}

	var item HNItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return []Discrepancy{{Field: "api_error", FileValue: "", APIValue: err.Error()}}, nil
	}

	// Check if item exists
	if item.ID == 0 {
		return []Discrepancy{{Field: "item_id", FileValue: summary.Meta.ItemID, APIValue: "not found"}}, nil
	}

	// Compare article URL
	if item.URL != "" && item.URL != summary.Meta.ArticleURL {
		discrepancies = append(discrepancies, Discrepancy{
			Field:     "article_url",
			FileValue: summary.Meta.ArticleURL,
			APIValue:  item.URL,
		})
	}

	// Compare discussion URL
	expectedDiscussionURL := fmt.Sprintf("https://news.ycombinator.com/item?id=%d", item.ID)
	if summary.Meta.DiscussionURL != expectedDiscussionURL {
		discrepancies = append(discrepancies, Discrepancy{
			Field:     "discussion_url",
			FileValue: summary.Meta.DiscussionURL,
			APIValue:  expectedDiscussionURL,
		})
	}

	// Compare title
	if item.Title != "" && item.Title != summary.Article.Title {
		discrepancies = append(discrepancies, Discrepancy{
			Field:     "title",
			FileValue: summary.Article.Title,
			APIValue:  item.Title,
		})
	}

	// Compare comment count
	if item.Descendants > 0 {
		fileCount := summary.Meta.DiscussionCommentCount
		apiCount := item.Descendants
		// Flag if difference is more than 10% or more than 50 comments
		diff := abs(fileCount - apiCount)
		if diff > 50 || (apiCount > 0 && float64(diff)/float64(apiCount) > 0.1) {
			discrepancies = append(discrepancies, Discrepancy{
				Field:     "comment_count",
				FileValue: fmt.Sprintf("%d", fileCount),
				APIValue:  fmt.Sprintf("%d", apiCount),
			})
		}
	}

	// Verify quotes if requested
	if verifyQuotes {
		quoteVers = verifyHNQuotes(summary, item, client)
	}

	return discrepancies, quoteVers
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// verifyHNQuotes fetches all comments and verifies persona quotes.
func verifyHNQuotes(summary types.Summary, item HNItem, client *http.Client) []QuoteVerification {
	var results []QuoteVerification

	// Collect all quotes from personas
	var quotes []struct {
		PersonaName string
		Author      string
		Text        string
	}

	for _, persona := range summary.Discussion.Personas {
		for _, quote := range persona.Quotes {
			quotes = append(quotes, struct {
				PersonaName string
				Author      string
				Text        string
			}{
				PersonaName: persona.Name,
				Author:      quote.Author,
				Text:        quote.Text,
			})
		}
	}

	if len(quotes) == 0 {
		return nil
	}

	if verifyVerbose {
		fmt.Printf("    Fetching comments to verify %d quotes...\n", len(quotes))
	}

	// Fetch all comments (limited to prevent API abuse)
	comments := fetchHNComments(item.Kids, client, 500) // max 500 comments

	// Build author -> comments map
	authorComments := make(map[string][]string)
	for _, c := range comments {
		authorComments[c.Author] = append(authorComments[c.Author], c.Text)
	}

	if verifyVerbose {
		fmt.Printf("    Fetched %d comments from %d authors\n", len(comments), len(authorComments))
	}

	// Verify each quote
	for _, q := range quotes {
		ver := QuoteVerification{
			PersonaName: q.PersonaName,
			QuoteAuthor: q.Author,
			QuoteText:   q.Text,
		}

		authorTexts, exists := authorComments[q.Author]
		ver.AuthorExists = exists

		if exists {
			// Check if any comment contains similar text
			bestSim := 0.0
			bestMatch := ""
			for _, text := range authorTexts {
				sim := textutil.WordSimilarity(q.Text, text)
				if sim > bestSim {
					bestSim = sim
					bestMatch = text
				}
			}
			ver.Similarity = bestSim
			ver.BestMatch = textutil.Truncate(bestMatch, 100)
			ver.TextFound = bestSim > 0.3 // threshold for "similar enough"
		}

		results = append(results, ver)
	}

	return results
}

// fetchHNComments recursively fetches comments up to maxCount.
func fetchHNComments(ids []int, client *http.Client, maxCount int) []HNComment {
	var comments []HNComment
	toFetch := ids

	for len(toFetch) > 0 && len(comments) < maxCount {
		// Fetch next batch (limit concurrent requests)
		batchSize := 10
		if batchSize > len(toFetch) {
			batchSize = len(toFetch)
		}
		if len(comments)+batchSize > maxCount {
			batchSize = maxCount - len(comments)
		}

		batch := toFetch[:batchSize]
		toFetch = toFetch[batchSize:]

		for _, id := range batch {
			item := fetchHNItem(id, client)
			if item == nil || item.Deleted || item.Dead {
				continue
			}

			if item.By != "" && item.Text != "" {
				comments = append(comments, HNComment{
					ID:     item.ID,
					Author: item.By,
					Text:   textutil.StripHTMLAndQuotes(item.Text),
				})
			}

			// Add children to queue
			if len(item.Kids) > 0 {
				toFetch = append(toFetch, item.Kids...)
			}
		}

		// Rate limit to avoid hammering the API
		time.Sleep(50 * time.Millisecond)
	}

	return comments
}

func fetchHNItem(id int, client *http.Client) *HNItem {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var item HNItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil
	}
	return &item
}

func fixSummaryFile(filePath string, discrepancies []Discrepancy) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := string(data)

	for _, d := range discrepancies {
		switch d.Field {
		case "article_url":
			content = strings.Replace(content,
				fmt.Sprintf("article_url: \"%s\"", d.FileValue),
				fmt.Sprintf("article_url: \"%s\"", d.APIValue), 1)
			content = strings.Replace(content,
				fmt.Sprintf("article_url: %s", d.FileValue),
				fmt.Sprintf("article_url: \"%s\"", d.APIValue), 1)
		case "discussion_url":
			content = strings.Replace(content,
				fmt.Sprintf("discussion_url: \"%s\"", d.FileValue),
				fmt.Sprintf("discussion_url: \"%s\"", d.APIValue), 1)
			content = strings.Replace(content,
				fmt.Sprintf("discussion_url: %s", d.FileValue),
				fmt.Sprintf("discussion_url: \"%s\"", d.APIValue), 1)
		case "title":
			// Be careful with title - it may contain special characters
			// Only fix if it's a simple case
			oldTitle := fmt.Sprintf("title: \"%s\"", d.FileValue)
			newTitle := fmt.Sprintf("title: \"%s\"", d.APIValue)
			content = strings.Replace(content, oldTitle, newTitle, 1)
		case "comment_count":
			// Fix comment count
			content = strings.Replace(content,
				fmt.Sprintf("discussion_comment_count: %s", d.FileValue),
				fmt.Sprintf("discussion_comment_count: %s", d.APIValue), 1)
		}
	}

	return os.WriteFile(filePath, []byte(content), 0644)
}
