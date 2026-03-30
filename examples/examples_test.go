package examples

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/socialpulse/types"
)

func TestExamplesValidJSON(t *testing.T) {
	files, err := filepath.Glob("*.json")
	if err != nil {
		t.Fatalf("failed to glob examples: %v", err)
	}

	if len(files) == 0 {
		t.Skip("no example files found")
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("failed to read %s: %v", file, err)
			}

			var summary types.Summary
			if err := json.Unmarshal(data, &summary); err != nil {
				t.Errorf("failed to unmarshal %s: %v", file, err)
				return
			}

			// Validate required fields
			if summary.SchemaVersion == "" {
				t.Error("missing schema_version")
			}
			if summary.Meta.ArticleURL == "" {
				t.Error("missing article_url")
			}
			if summary.Meta.Platform == "" {
				t.Error("missing platform")
			}
			if summary.Article.Title == "" {
				t.Error("missing article title")
			}
			if len(summary.Discussion.Personas) == 0 {
				t.Error("no personas defined")
			}

			// Validate platform enum
			validPlatforms := map[types.Platform]bool{
				types.PlatformHackerNews: true,
				types.PlatformReddit:     true,
			}
			if !validPlatforms[summary.Meta.Platform] {
				t.Errorf("invalid platform: %s", summary.Meta.Platform)
			}

			// Validate enum values
			validArticleSentiments := map[types.ArticleSentiment]bool{
				types.SentimentOptimistic:  true,
				types.SentimentCautious:    true,
				types.SentimentPessimistic: true,
				types.SentimentNeutral:     true,
				types.SentimentProvocative: true,
			}
			if !validArticleSentiments[summary.Article.Sentiment] {
				t.Errorf("invalid article sentiment: %s", summary.Article.Sentiment)
			}

			validDiscussionSentiments := map[types.DiscussionSentiment]bool{
				types.DiscussionSupportive:  true,
				types.DiscussionDivided:     true,
				types.DiscussionCritical:    true,
				types.DiscussionExploratory: true,
			}
			if !validDiscussionSentiments[summary.Discussion.Sentiment] {
				t.Errorf("invalid discussion sentiment: %s", summary.Discussion.Sentiment)
			}

			// Validate personas
			for i, persona := range summary.Discussion.Personas {
				validStances := map[types.Stance]bool{
					types.StanceAgrees:     true,
					types.StanceDisagrees:  true,
					types.StanceNuanced:    true,
					types.StanceTangential: true,
				}
				if !validStances[persona.Stance] {
					t.Errorf("persona %d (%s) has invalid stance: %s", i, persona.Name, persona.Stance)
				}

				validPrevalences := map[types.Prevalence]bool{
					types.PrevalenceDominant:    true,
					types.PrevalenceSignificant: true,
					types.PrevalenceMinority:    true,
				}
				if !validPrevalences[persona.Prevalence] {
					t.Errorf("persona %d (%s) has invalid prevalence: %s", i, persona.Name, persona.Prevalence)
				}

				if len(persona.Quotes) == 0 {
					t.Errorf("persona %d (%s) has no quotes", i, persona.Name)
				}
			}

			t.Logf("validated %s: platform=%s, %d personas, %d tags",
				file, summary.Meta.Platform, len(summary.Discussion.Personas), len(summary.Article.Tags))
		})
	}
}

func TestExamplesHaveRequiredContent(t *testing.T) {
	files, err := filepath.Glob("*.json")
	if err != nil {
		t.Fatalf("failed to glob examples: %v", err)
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("failed to read %s: %v", file, err)
			}

			var summary types.Summary
			if err := json.Unmarshal(data, &summary); err != nil {
				t.Skipf("skipping content validation for invalid JSON")
			}

			// Check thesis length (recommended 200-400 chars)
			thesisLen := len(summary.Article.Thesis)
			if thesisLen < 50 {
				t.Logf("warning: thesis is short (%d chars, recommended 200-400)", thesisLen)
			}

			// Check key arguments count (recommended 3-5)
			argCount := len(summary.Article.KeyArguments)
			if argCount < 2 || argCount > 6 {
				t.Logf("warning: key_arguments count is %d (recommended 3-5)", argCount)
			}

			// Check personas count (recommended 3-6)
			personaCount := len(summary.Discussion.Personas)
			if personaCount < 2 || personaCount > 7 {
				t.Logf("warning: personas count is %d (recommended 3-6)", personaCount)
			}

			// Check that quotes exist and have content
			for _, persona := range summary.Discussion.Personas {
				for j, quote := range persona.Quotes {
					if strings.TrimSpace(quote.Text) == "" {
						t.Errorf("persona %s quote %d has empty text", persona.Name, j)
					}
					if strings.TrimSpace(quote.Author) == "" {
						t.Errorf("persona %s quote %d has empty author", persona.Name, j)
					}
				}
			}
		})
	}
}
