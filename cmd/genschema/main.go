// Command genschema generates JSON Schema from Go types.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grokify/socialpulse/types"
	"github.com/invopop/jsonschema"
)

type schemaConfig struct {
	Type        any
	Filename    string
	Title       string
	Description string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	schemas := []schemaConfig{
		{
			Type:        &types.Summary{},
			Filename:    "summary.schema.json",
			Title:       "SocialPulse Summary",
			Description: "A structured summary of a blog article and its social discussion",
		},
		{
			Type:        &types.Digest{},
			Filename:    "digest.schema.json",
			Title:       "SocialPulse Digest",
			Description: "An aggregated summary of articles over a time period",
		},
		{
			Type:        &types.Site{},
			Filename:    "site.schema.json",
			Title:       "SocialPulse Site",
			Description: "Site configuration and article index",
		},
		{
			Type:        &types.SourceLinks{},
			Filename:    "source-links.schema.json",
			Title:       "SocialPulse Source Links",
			Description: "Collection of article links to be summarized",
		},
	}

	outputDir := "schema"
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	// Ensure directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	for _, cfg := range schemas {
		if err := generateSchema(cfg, outputDir); err != nil {
			return err
		}
	}

	return nil
}

func generateSchema(cfg schemaConfig, outputDir string) error {
	r := new(jsonschema.Reflector)
	r.ExpandedStruct = true

	schema := r.Reflect(cfg.Type)
	schema.ID = jsonschema.ID(fmt.Sprintf("https://github.com/grokify/socialpulse/schema/%s", cfg.Filename))
	schema.Title = cfg.Title
	schema.Description = cfg.Description

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema %s: %w", cfg.Filename, err)
	}

	outputPath := filepath.Join(outputDir, cfg.Filename)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write schema %s: %w", outputPath, err)
	}

	fmt.Printf("Schema written to %s\n", outputPath)
	return nil
}
