// Package types defines the data structures for SocialPulse article summaries.
package types

import (
	"time"

	"github.com/invopop/jsonschema"
)

// Summary represents a complete summary of an article and its social discussion.
type Summary struct {
	// Schema version for forward compatibility
	SchemaVersion string `json:"schema_version" yaml:"schema_version"`

	// Metadata about the summary
	Meta Meta `json:"meta" yaml:"meta"`

	// Summary of the source article
	Article Article `json:"article" yaml:"article"`

	// Summary of the social discussion
	Discussion Discussion `json:"discussion" yaml:"discussion"`
}

// Platform represents the source platform for a discussion.
type Platform string

const (
	PlatformHackerNews Platform = "hackernews"
	PlatformReddit     Platform = "reddit"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (Platform) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			PlatformHackerNews,
			PlatformReddit,
		},
		Description: "Source platform for the discussion",
	}
}

// Meta contains metadata about the summary and its sources.
type Meta struct {
	// URL of the original article
	ArticleURL string `json:"article_url" yaml:"article_url"`

	// URL of the social discussion
	DiscussionURL string `json:"discussion_url" yaml:"discussion_url"`

	// Platform where the discussion took place
	Platform Platform `json:"platform" yaml:"platform"`

	// Publication date of the article
	ArticleDate time.Time `json:"article_date" yaml:"article_date"`

	// Date this summary was created
	SummaryDate time.Time `json:"summary_date" yaml:"summary_date"`

	// Author of the original article
	ArticleAuthor string `json:"article_author" yaml:"article_author"`

	// Number of comments in the discussion
	DiscussionCommentCount int `json:"discussion_comment_count" yaml:"discussion_comment_count"`

	// Platform-specific item ID (e.g., HN item ID, Reddit post ID)
	ItemID string `json:"item_id,omitempty" yaml:"item_id,omitempty"`

	// Subreddit name (Reddit only)
	Subreddit string `json:"subreddit,omitempty" yaml:"subreddit,omitempty"`

	// Discussion score/points
	Score int `json:"score,omitempty" yaml:"score,omitempty"`
}

// Article summarizes the source article content.
type Article struct {
	// Title of the article
	Title string `json:"title" yaml:"title"`

	// Core thesis statement (200-400 chars recommended)
	Thesis string `json:"thesis" yaml:"thesis"`

	// Key arguments made by the author (3-5 items, 100-200 chars each)
	KeyArguments []string `json:"key_arguments" yaml:"key_arguments"`

	// Concrete examples cited by the author
	Examples []string `json:"examples,omitempty" yaml:"examples,omitempty"`

	// Author's proposed solution or call to action (100-200 chars)
	Prescription string `json:"prescription,omitempty" yaml:"prescription,omitempty"`

	// Categorical tags for the article
	Tags []string `json:"tags" yaml:"tags"`

	// Overall sentiment of the article
	Sentiment ArticleSentiment `json:"sentiment" yaml:"sentiment"`
}

// ArticleSentiment represents the overall tone of an article.
type ArticleSentiment string

const (
	SentimentOptimistic  ArticleSentiment = "optimistic"
	SentimentCautious    ArticleSentiment = "cautious"
	SentimentPessimistic ArticleSentiment = "pessimistic"
	SentimentNeutral     ArticleSentiment = "neutral"
	SentimentProvocative ArticleSentiment = "provocative"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (ArticleSentiment) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			SentimentOptimistic,
			SentimentCautious,
			SentimentPessimistic,
			SentimentNeutral,
			SentimentProvocative,
		},
		Description: "Overall tone of the article",
	}
}

// Discussion summarizes the social discussion.
type Discussion struct {
	// Overall sentiment of the discussion
	Sentiment DiscussionSentiment `json:"sentiment" yaml:"sentiment"`

	// Commenter archetypes identified in the discussion (3-6 recommended)
	Personas []Persona `json:"personas" yaml:"personas"`

	// Emergent sub-topics not directly addressed in the original article
	Tangents []Tangent `json:"tangents,omitempty" yaml:"tangents,omitempty"`

	// Points of broad agreement across personas
	ConsensusPoints []string `json:"consensus_points,omitempty" yaml:"consensus_points,omitempty"`

	// Unresolved tensions or questions
	OpenQuestions []string `json:"open_questions,omitempty" yaml:"open_questions,omitempty"`
}

// DiscussionSentiment represents the overall tone of a discussion.
type DiscussionSentiment string

const (
	DiscussionSupportive  DiscussionSentiment = "supportive"
	DiscussionDivided     DiscussionSentiment = "divided"
	DiscussionCritical    DiscussionSentiment = "critical"
	DiscussionExploratory DiscussionSentiment = "exploratory"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (DiscussionSentiment) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			DiscussionSupportive,
			DiscussionDivided,
			DiscussionCritical,
			DiscussionExploratory,
		},
		Description: "Overall tone of the discussion",
	}
}

// Persona represents an archetype of commenters with a shared viewpoint.
type Persona struct {
	// Descriptive name for this persona (e.g., "The Skeptic", "The Pragmatist")
	Name string `json:"name" yaml:"name"`

	// Brief description of this viewpoint (50-100 chars)
	Description string `json:"description" yaml:"description"`

	// How this persona relates to the article's thesis
	Stance Stance `json:"stance" yaml:"stance"`

	// How common this viewpoint is in the discussion
	Prevalence Prevalence `json:"prevalence" yaml:"prevalence"`

	// Estimated percentage of discussion participants fitting this persona (0-100)
	Percentage int `json:"percentage" yaml:"percentage"`

	// Verbatim quotes from commenters fitting this persona (2-3 recommended)
	Quotes []Quote `json:"quotes" yaml:"quotes"`

	// Core argument distilled (100-150 chars)
	CoreArgument string `json:"core_argument" yaml:"core_argument"`
}

// Stance represents how a persona relates to the article's thesis.
type Stance string

const (
	StanceAgrees     Stance = "agrees"
	StanceDisagrees  Stance = "disagrees"
	StanceNuanced    Stance = "nuanced"
	StanceTangential Stance = "tangential"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (Stance) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			StanceAgrees,
			StanceDisagrees,
			StanceNuanced,
			StanceTangential,
		},
		Description: "How a persona relates to the article's thesis",
	}
}

// Prevalence represents how common a viewpoint is in the discussion.
type Prevalence string

const (
	PrevalenceDominant    Prevalence = "dominant"
	PrevalenceSignificant Prevalence = "significant"
	PrevalenceMinority    Prevalence = "minority"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (Prevalence) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			PrevalenceDominant,
			PrevalenceSignificant,
			PrevalenceMinority,
		},
		Description: "How common a viewpoint is in the discussion",
	}
}

// Quote represents a verbatim quote from a commenter.
type Quote struct {
	// Exact text of the quote
	Text string `json:"text" yaml:"text"`

	// Username of the commenter
	Author string `json:"author" yaml:"author"`

	// Optional context for what prompted this quote (50 chars max)
	Context string `json:"context,omitempty" yaml:"context,omitempty"`
}

// Tangent represents an emergent sub-topic in the discussion.
type Tangent struct {
	// Name of the tangent topic
	Topic string `json:"topic" yaml:"topic"`

	// Brief summary of the tangent (100-150 chars)
	Summary string `json:"summary" yaml:"summary"`
}
