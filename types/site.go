package types

import "time"

// Site represents the complete site configuration and index.
type Site struct {
	// Schema version for forward compatibility
	SchemaVersion string `json:"schema_version" yaml:"schema_version"`

	// Site metadata
	Meta SiteMeta `json:"meta" yaml:"meta"`

	// All article summaries, ordered by date (newest first)
	Articles []ArticleReference `json:"articles" yaml:"articles"`

	// Available digests
	Digests []DigestReference `json:"digests,omitempty" yaml:"digests,omitempty"`

	// Tag index for navigation
	TagIndex []TagIndex `json:"tag_index,omitempty" yaml:"tag_index,omitempty"`
}

// SiteMeta contains site-level metadata.
type SiteMeta struct {
	// Site title
	Title string `json:"title" yaml:"title"`

	// Site description
	Description string `json:"description" yaml:"description"`

	// Base URL of the site
	BaseURL string `json:"base_url" yaml:"base_url"`

	// Date of last update
	LastUpdated time.Time `json:"last_updated" yaml:"last_updated"`

	// Total article count
	ArticleCount int `json:"article_count" yaml:"article_count"`

	// Platforms covered
	Platforms []Platform `json:"platforms" yaml:"platforms"`
}

// DigestReference is a lightweight reference to a digest.
type DigestReference struct {
	// Digest period type
	Period DigestPeriod `json:"period" yaml:"period"`

	// Period start date
	PeriodStart time.Time `json:"period_start" yaml:"period_start"`

	// Period end date
	PeriodEnd time.Time `json:"period_end" yaml:"period_end"`

	// Path to the digest file
	DigestPath string `json:"digest_path" yaml:"digest_path"`

	// Number of articles in this digest
	ArticleCount int `json:"article_count" yaml:"article_count"`
}

// TagIndex represents all articles with a specific tag.
type TagIndex struct {
	// Tag name
	Tag string `json:"tag" yaml:"tag"`

	// Number of articles with this tag
	Count int `json:"count" yaml:"count"`

	// References to articles with this tag
	Articles []ArticleReference `json:"articles" yaml:"articles"`
}

// SourceLink represents an article to be summarized.
type SourceLink struct {
	// URL of the original article
	ArticleURL string `json:"article_url" yaml:"article_url"`

	// URL of the social discussion
	DiscussionURL string `json:"discussion_url" yaml:"discussion_url"`

	// Platform of the discussion
	Platform Platform `json:"platform" yaml:"platform"`

	// Optional notes about why this was added
	Notes string `json:"notes,omitempty" yaml:"notes,omitempty"`

	// Date this link was added
	AddedDate time.Time `json:"added_date" yaml:"added_date"`

	// Whether this has been processed
	Processed bool `json:"processed" yaml:"processed"`

	// Path to generated summary (if processed)
	SummaryPath string `json:"summary_path,omitempty" yaml:"summary_path,omitempty"`
}

// SourceLinks is a collection of articles to be summarized.
type SourceLinks struct {
	// Schema version
	SchemaVersion string `json:"schema_version" yaml:"schema_version"`

	// Links to process
	Links []SourceLink `json:"links" yaml:"links"`
}
