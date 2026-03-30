package types

import (
	"time"

	"github.com/invopop/jsonschema"
)

// DigestPeriod represents the time period for a digest.
type DigestPeriod string

const (
	DigestWeekly    DigestPeriod = "weekly"
	DigestMonthly   DigestPeriod = "monthly"
	DigestQuarterly DigestPeriod = "quarterly"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (DigestPeriod) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			DigestWeekly,
			DigestMonthly,
			DigestQuarterly,
		},
		Description: "Time period covered by the digest",
	}
}

// Digest represents an aggregated summary of articles over a time period.
type Digest struct {
	// Schema version for forward compatibility
	SchemaVersion string `json:"schema_version" yaml:"schema_version"`

	// Metadata about the digest
	Meta DigestMeta `json:"meta" yaml:"meta"`

	// Aggregated trends across all articles
	Trends Trends `json:"trends" yaml:"trends"`

	// Aggregated personas across all discussions
	Personas []AggregatedPersona `json:"personas" yaml:"personas"`

	// Top articles by engagement
	TopArticles []ArticleReference `json:"top_articles" yaml:"top_articles"`

	// Notable quotes from the period
	NotableQuotes []AttributedQuote `json:"notable_quotes,omitempty" yaml:"notable_quotes,omitempty"`
}

// DigestMeta contains metadata about the digest.
type DigestMeta struct {
	// Time period type
	Period DigestPeriod `json:"period" yaml:"period"`

	// Start of the period (inclusive)
	PeriodStart time.Time `json:"period_start" yaml:"period_start"`

	// End of the period (inclusive)
	PeriodEnd time.Time `json:"period_end" yaml:"period_end"`

	// Date this digest was generated
	GeneratedDate time.Time `json:"generated_date" yaml:"generated_date"`

	// Number of articles summarized in this period
	ArticleCount int `json:"article_count" yaml:"article_count"`

	// Total comments across all discussions
	TotalComments int `json:"total_comments" yaml:"total_comments"`

	// Platforms included in this digest
	Platforms []Platform `json:"platforms" yaml:"platforms"`
}

// Trends represents aggregated trends over a time period.
type Trends struct {
	// Most common tags across articles, with counts
	TopTags []TagCount `json:"top_tags" yaml:"top_tags"`

	// Overall sentiment distribution for articles
	ArticleSentimentDistribution SentimentDistribution `json:"article_sentiment_distribution" yaml:"article_sentiment_distribution"`

	// Overall sentiment distribution for discussions
	DiscussionSentimentDistribution SentimentDistribution `json:"discussion_sentiment_distribution" yaml:"discussion_sentiment_distribution"`

	// Emerging topics that appeared multiple times
	EmergingTopics []TopicTrend `json:"emerging_topics,omitempty" yaml:"emerging_topics,omitempty"`

	// Key themes observed across multiple articles
	KeyThemes []string `json:"key_themes" yaml:"key_themes"`

	// Brief narrative summary of the period (500-1000 chars)
	NarrativeSummary string `json:"narrative_summary" yaml:"narrative_summary"`
}

// TagCount represents a tag with its occurrence count.
type TagCount struct {
	// Tag name
	Tag string `json:"tag" yaml:"tag"`

	// Number of articles with this tag
	Count int `json:"count" yaml:"count"`
}

// SentimentDistribution shows the breakdown of sentiments.
type SentimentDistribution struct {
	Optimistic  int `json:"optimistic" yaml:"optimistic"`
	Cautious    int `json:"cautious" yaml:"cautious"`
	Pessimistic int `json:"pessimistic" yaml:"pessimistic"`
	Neutral     int `json:"neutral" yaml:"neutral"`
	Provocative int `json:"provocative" yaml:"provocative"`
}

// TopicTrend represents an emerging topic with metadata.
type TopicTrend struct {
	// Topic name
	Topic string `json:"topic" yaml:"topic"`

	// Number of articles/discussions where this appeared
	Occurrences int `json:"occurrences" yaml:"occurrences"`

	// Brief description of why this is trending
	Context string `json:"context,omitempty" yaml:"context,omitempty"`
}

// AggregatedPersona represents a persona type observed across multiple discussions.
type AggregatedPersona struct {
	// Persona archetype name
	Name string `json:"name" yaml:"name"`

	// Description of this persona type
	Description string `json:"description" yaml:"description"`

	// How often this persona appeared across discussions
	Frequency int `json:"frequency" yaml:"frequency"`

	// Typical stance this persona takes
	TypicalStance Stance `json:"typical_stance" yaml:"typical_stance"`

	// Representative quotes from this persona type
	RepresentativeQuotes []AttributedQuote `json:"representative_quotes" yaml:"representative_quotes"`

	// Common arguments made by this persona type
	CommonArguments []string `json:"common_arguments" yaml:"common_arguments"`
}

// ArticleReference is a lightweight reference to an article summary.
type ArticleReference struct {
	// Article title
	Title string `json:"title" yaml:"title"`

	// URL of the original article
	ArticleURL string `json:"article_url" yaml:"article_url"`

	// URL of the discussion
	DiscussionURL string `json:"discussion_url" yaml:"discussion_url"`

	// Platform of the discussion
	Platform Platform `json:"platform" yaml:"platform"`

	// Publication date
	ArticleDate time.Time `json:"article_date" yaml:"article_date"`

	// Primary tags
	Tags []string `json:"tags" yaml:"tags"`

	// Article sentiment
	Sentiment ArticleSentiment `json:"sentiment" yaml:"sentiment"`

	// Discussion comment count
	CommentCount int `json:"comment_count" yaml:"comment_count"`

	// Discussion score/points
	Score int `json:"score,omitempty" yaml:"score,omitempty"`

	// Path to the full summary file
	SummaryPath string `json:"summary_path,omitempty" yaml:"summary_path,omitempty"`
}

// AttributedQuote is a quote with full attribution for digest context.
type AttributedQuote struct {
	// The quote text
	Text string `json:"text" yaml:"text"`

	// Username of the commenter
	Author string `json:"author" yaml:"author"`

	// Platform where this was posted
	Platform Platform `json:"platform" yaml:"platform"`

	// Reference to the source article
	ArticleTitle string `json:"article_title" yaml:"article_title"`

	// URL of the discussion
	DiscussionURL string `json:"discussion_url" yaml:"discussion_url"`

	// Why this quote was notable
	Significance string `json:"significance,omitempty" yaml:"significance,omitempty"`
}
