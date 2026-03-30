package schema

import (
	"encoding/json"
	"testing"

	"github.com/grokify/socialpulse/types"
)

func TestSummarySchemaNotEmpty(t *testing.T) {
	data := SummarySchema()
	if len(data) == 0 {
		t.Error("embedded schema is empty")
	}
}

func TestSummarySchemaValidJSON(t *testing.T) {
	data := SummarySchema()
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Errorf("schema is not valid JSON: %v", err)
	}
}

func TestSummarySchemaMap(t *testing.T) {
	schema, err := SummarySchemaMap()
	if err != nil {
		t.Fatalf("SummarySchemaMap failed: %v", err)
	}

	if schema["$schema"] == nil {
		t.Error("missing $schema field")
	}
	if schema["$id"] == nil {
		t.Error("missing $id field")
	}
	if schema["title"] != "SocialPulse Summary" {
		t.Errorf("unexpected title: %v", schema["title"])
	}
}

func TestDigestSchemaNotEmpty(t *testing.T) {
	data := DigestSchema()
	if len(data) == 0 {
		t.Error("embedded digest schema is empty")
	}
}

func TestDigestSchemaMap(t *testing.T) {
	schema, err := DigestSchemaMap()
	if err != nil {
		t.Fatalf("DigestSchemaMap failed: %v", err)
	}

	if schema["title"] != "SocialPulse Digest" {
		t.Errorf("unexpected title: %v", schema["title"])
	}
}

func TestSiteSchemaNotEmpty(t *testing.T) {
	data := SiteSchema()
	if len(data) == 0 {
		t.Error("embedded site schema is empty")
	}
}

func TestSiteSchemaMap(t *testing.T) {
	schema, err := SiteSchemaMap()
	if err != nil {
		t.Fatalf("SiteSchemaMap failed: %v", err)
	}

	if schema["title"] != "SocialPulse Site" {
		t.Errorf("unexpected title: %v", schema["title"])
	}
}

func TestSourceLinksSchemaNotEmpty(t *testing.T) {
	data := SourceLinksSchema()
	if len(data) == 0 {
		t.Error("embedded source-links schema is empty")
	}
}

func TestSourceLinksSchemaMap(t *testing.T) {
	schema, err := SourceLinksSchemaMap()
	if err != nil {
		t.Fatalf("SourceLinksSchemaMap failed: %v", err)
	}

	if schema["title"] != "SocialPulse Source Links" {
		t.Errorf("unexpected title: %v", schema["title"])
	}
}

func TestSchemaMatchesTypes(t *testing.T) {
	schema, err := SummarySchemaMap()
	if err != nil {
		t.Fatalf("failed to get schema: %v", err)
	}

	defs, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatal("missing $defs in schema")
	}

	// Verify enum definitions exist
	expectedEnums := []string{
		"ArticleSentiment",
		"DiscussionSentiment",
		"Stance",
		"Prevalence",
		"Platform",
	}

	for _, name := range expectedEnums {
		def, ok := defs[name].(map[string]any)
		if !ok {
			t.Errorf("missing definition for %s", name)
			continue
		}
		if def["enum"] == nil {
			t.Errorf("%s missing enum constraint", name)
		}
	}
}

func TestSummaryValidatesAgainstSchema(t *testing.T) {
	summary := types.Summary{
		SchemaVersion: "1.0.0",
		Meta: types.Meta{
			ArticleURL:             "https://example.com/article",
			DiscussionURL:          "https://news.ycombinator.com/item?id=12345",
			Platform:               types.PlatformHackerNews,
			ArticleAuthor:          "Test Author",
			DiscussionCommentCount: 10,
		},
		Article: types.Article{
			Title:        "Test Article",
			Thesis:       "This is a test.",
			KeyArguments: []string{"Argument 1"},
			Tags:         []string{"test"},
			Sentiment:    types.SentimentNeutral,
		},
		Discussion: types.Discussion{
			Sentiment: types.DiscussionDivided,
			Personas: []types.Persona{
				{
					Name:         "Test Persona",
					Description:  "A test persona",
					Stance:       types.StanceNuanced,
					Prevalence:   types.PrevalenceSignificant,
					Quotes:       []types.Quote{{Text: "Quote", Author: "user1"}},
					CoreArgument: "Core argument here",
				},
			},
		},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("failed to marshal summary: %v", err)
	}

	var decoded types.Summary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal summary: %v", err)
	}

	if decoded.Meta.Platform != types.PlatformHackerNews {
		t.Errorf("platform mismatch: got %q, want %q", decoded.Meta.Platform, types.PlatformHackerNews)
	}
}

func TestDigestValidatesAgainstSchema(t *testing.T) {
	digest := types.Digest{
		SchemaVersion: "1.0.0",
		Meta: types.DigestMeta{
			Period:       types.DigestWeekly,
			ArticleCount: 5,
			Platforms:    []types.Platform{types.PlatformHackerNews, types.PlatformReddit},
		},
		Trends: types.Trends{
			TopTags:          []types.TagCount{{Tag: "ai", Count: 3}},
			KeyThemes:        []string{"AI development"},
			NarrativeSummary: "A week of AI discussions.",
		},
		Personas: []types.AggregatedPersona{
			{
				Name:          "The Skeptic",
				Description:   "Questions AI claims",
				Frequency:     3,
				TypicalStance: types.StanceDisagrees,
			},
		},
		TopArticles: []types.ArticleReference{
			{
				Title:      "Test Article",
				ArticleURL: "https://example.com",
				Platform:   types.PlatformHackerNews,
			},
		},
	}

	data, err := json.Marshal(digest)
	if err != nil {
		t.Fatalf("failed to marshal digest: %v", err)
	}

	var decoded types.Digest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal digest: %v", err)
	}
}
