// Package schema provides embedded JSON Schema for SocialPulse types.
package schema

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed summary.schema.json
var summarySchemaBytes []byte

//go:embed digest.schema.json
var digestSchemaBytes []byte

//go:embed site.schema.json
var siteSchemaBytes []byte

//go:embed source-links.schema.json
var sourceLinksSchemaBytes []byte

// SummarySchema returns the raw JSON Schema bytes for the Summary type.
func SummarySchema() []byte {
	return summarySchemaBytes
}

// SummarySchemaMap returns the JSON Schema as a map for programmatic access.
func SummarySchemaMap() (map[string]any, error) {
	return parseSchema(summarySchemaBytes)
}

// DigestSchema returns the raw JSON Schema bytes for the Digest type.
func DigestSchema() []byte {
	return digestSchemaBytes
}

// DigestSchemaMap returns the JSON Schema as a map for programmatic access.
func DigestSchemaMap() (map[string]any, error) {
	return parseSchema(digestSchemaBytes)
}

// SiteSchema returns the raw JSON Schema bytes for the Site type.
func SiteSchema() []byte {
	return siteSchemaBytes
}

// SiteSchemaMap returns the JSON Schema as a map for programmatic access.
func SiteSchemaMap() (map[string]any, error) {
	return parseSchema(siteSchemaBytes)
}

// SourceLinksSchema returns the raw JSON Schema bytes for the SourceLinks type.
func SourceLinksSchema() []byte {
	return sourceLinksSchemaBytes
}

// SourceLinksSchemaMap returns the JSON Schema as a map for programmatic access.
func SourceLinksSchemaMap() (map[string]any, error) {
	return parseSchema(sourceLinksSchemaBytes)
}

func parseSchema(data []byte) (map[string]any, error) {
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}
	return schema, nil
}
