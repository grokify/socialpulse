package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSummaryRoundTrip(t *testing.T) {
	summary := Summary{
		SchemaVersion: "1.0.0",
		Meta: Meta{
			ArticleURL:             "https://example.com/article",
			DiscussionURL:          "https://news.ycombinator.com/item?id=12345",
			Platform:               PlatformHackerNews,
			ArticleDate:            time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC),
			SummaryDate:            time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC),
			ArticleAuthor:          "Jane Doe",
			DiscussionCommentCount: 42,
			ItemID:                 "12345",
		},
		Article: Article{
			Title:  "Test Article",
			Thesis: "This is a test thesis statement.",
			KeyArguments: []string{
				"First key argument",
				"Second key argument",
			},
			Examples:     []string{"Example one"},
			Prescription: "Do this instead.",
			Tags:         []string{"testing", "go"},
			Sentiment:    SentimentNeutral,
		},
		Discussion: Discussion{
			Sentiment: DiscussionDivided,
			Personas: []Persona{
				{
					Name:        "The Skeptic",
					Description: "Questions everything",
					Stance:      StanceDisagrees,
					Prevalence:  PrevalenceSignificant,
					Quotes: []Quote{
						{
							Text:    "I don't buy it.",
							Author:  "skeptic123",
							Context: "Responding to thesis",
						},
					},
					CoreArgument: "The evidence doesn't support the claim.",
				},
			},
			Tangents: []Tangent{
				{
					Topic:   "Related topic",
					Summary: "Discussion veered into this area.",
				},
			},
			ConsensusPoints: []string{"Everyone agreed on this"},
			OpenQuestions:   []string{"What about edge cases?"},
		},
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal summary: %v", err)
	}

	// Unmarshal back
	var decoded Summary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal summary: %v", err)
	}

	// Verify key fields
	if decoded.SchemaVersion != summary.SchemaVersion {
		t.Errorf("schema version mismatch: got %q, want %q", decoded.SchemaVersion, summary.SchemaVersion)
	}
	if decoded.Meta.ArticleAuthor != summary.Meta.ArticleAuthor {
		t.Errorf("article author mismatch: got %q, want %q", decoded.Meta.ArticleAuthor, summary.Meta.ArticleAuthor)
	}
	if decoded.Meta.Platform != PlatformHackerNews {
		t.Errorf("platform mismatch: got %q, want %q", decoded.Meta.Platform, PlatformHackerNews)
	}
	if decoded.Article.Sentiment != SentimentNeutral {
		t.Errorf("article sentiment mismatch: got %q, want %q", decoded.Article.Sentiment, SentimentNeutral)
	}
	if len(decoded.Discussion.Personas) != 1 {
		t.Errorf("personas count mismatch: got %d, want 1", len(decoded.Discussion.Personas))
	}
}

func TestPlatformValues(t *testing.T) {
	platforms := []Platform{PlatformHackerNews, PlatformReddit}
	for _, p := range platforms {
		if p == "" {
			t.Errorf("platform should not be empty")
		}
	}
}

func TestSentimentValues(t *testing.T) {
	tests := []struct {
		sentiment ArticleSentiment
		valid     bool
	}{
		{SentimentOptimistic, true},
		{SentimentCautious, true},
		{SentimentPessimistic, true},
		{SentimentNeutral, true},
		{SentimentProvocative, true},
	}

	for _, tt := range tests {
		if tt.sentiment == "" {
			t.Errorf("sentiment should not be empty")
		}
	}
}

func TestStanceValues(t *testing.T) {
	stances := []Stance{StanceAgrees, StanceDisagrees, StanceNuanced, StanceTangential}
	for _, s := range stances {
		if s == "" {
			t.Errorf("stance should not be empty")
		}
	}
}

func TestPrevalenceValues(t *testing.T) {
	prevalences := []Prevalence{PrevalenceDominant, PrevalenceSignificant, PrevalenceMinority}
	for _, p := range prevalences {
		if p == "" {
			t.Errorf("prevalence should not be empty")
		}
	}
}

func TestDigestPeriodValues(t *testing.T) {
	periods := []DigestPeriod{DigestWeekly, DigestMonthly, DigestQuarterly}
	for _, p := range periods {
		if p == "" {
			t.Errorf("digest period should not be empty")
		}
	}
}

func TestDigestRoundTrip(t *testing.T) {
	digest := Digest{
		SchemaVersion: "1.0.0",
		Meta: DigestMeta{
			Period:        DigestWeekly,
			PeriodStart:   time.Date(2026, 3, 17, 0, 0, 0, 0, time.UTC),
			PeriodEnd:     time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC),
			GeneratedDate: time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC),
			ArticleCount:  5,
			TotalComments: 250,
			Platforms:     []Platform{PlatformHackerNews, PlatformReddit},
		},
		Trends: Trends{
			TopTags:          []TagCount{{Tag: "ai", Count: 3}, {Tag: "programming", Count: 2}},
			KeyThemes:        []string{"AI development concerns", "Software quality"},
			NarrativeSummary: "A week dominated by discussions about AI coding assistants.",
		},
		Personas: []AggregatedPersona{
			{
				Name:          "The Skeptic",
				Description:   "Questions AI hype",
				Frequency:     4,
				TypicalStance: StanceDisagrees,
				CommonArguments: []string{
					"Show me the evidence",
					"This has been tried before",
				},
			},
		},
		TopArticles: []ArticleReference{
			{
				Title:        "Test Article",
				ArticleURL:   "https://example.com",
				Platform:     PlatformHackerNews,
				CommentCount: 50,
			},
		},
	}

	data, err := json.MarshalIndent(digest, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal digest: %v", err)
	}

	var decoded Digest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal digest: %v", err)
	}

	if decoded.Meta.Period != DigestWeekly {
		t.Errorf("period mismatch: got %q, want %q", decoded.Meta.Period, DigestWeekly)
	}
	if len(decoded.Meta.Platforms) != 2 {
		t.Errorf("platforms count mismatch: got %d, want 2", len(decoded.Meta.Platforms))
	}
}

func TestSiteRoundTrip(t *testing.T) {
	site := Site{
		SchemaVersion: "1.0.0",
		Meta: SiteMeta{
			Title:        "FrontierNotes",
			Description:  "Summaries of tech articles and discussions",
			BaseURL:      "https://frontiernotes.example.com",
			LastUpdated:  time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC),
			ArticleCount: 10,
			Platforms:    []Platform{PlatformHackerNews, PlatformReddit},
		},
		Articles: []ArticleReference{
			{
				Title:      "Test Article",
				ArticleURL: "https://example.com",
				Platform:   PlatformHackerNews,
			},
		},
	}

	data, err := json.MarshalIndent(site, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal site: %v", err)
	}

	var decoded Site
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal site: %v", err)
	}

	if decoded.Meta.Title != "FrontierNotes" {
		t.Errorf("title mismatch: got %q, want %q", decoded.Meta.Title, "FrontierNotes")
	}
}

func TestSourceLinksRoundTrip(t *testing.T) {
	links := SourceLinks{
		SchemaVersion: "1.0.0",
		Links: []SourceLink{
			{
				ArticleURL:    "https://example.com/article",
				DiscussionURL: "https://news.ycombinator.com/item?id=12345",
				Platform:      PlatformHackerNews,
				AddedDate:     time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC),
				Processed:     false,
			},
		},
	}

	data, err := json.MarshalIndent(links, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal source links: %v", err)
	}

	var decoded SourceLinks
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal source links: %v", err)
	}

	if len(decoded.Links) != 1 {
		t.Errorf("links count mismatch: got %d, want 1", len(decoded.Links))
	}
}
