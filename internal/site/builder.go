// Package site implements the static site generator.
package site

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/grokify/socialpulse/internal/theme"
	"github.com/grokify/socialpulse/types"
)

// BuilderConfig contains configuration for the site builder.
type BuilderConfig struct {
	SiteTitle       string
	SiteDescription string
	BaseURL         string
	SummariesDir    string
	DigestsDir      string
	OutputDir       string
	ThemeName       string
}

// BuildResult contains statistics about the build.
type BuildResult struct {
	ArticleCount int
	DigestCount  int
	PageCount    int
}

// Builder generates static HTML from content files.
type Builder struct {
	config    BuilderConfig
	templates *template.Template
}

// NewBuilder creates a new site builder.
func NewBuilder(config BuilderConfig) *Builder {
	return &Builder{
		config: config,
	}
}

// Build generates the static site.
func (b *Builder) Build() (*BuildResult, error) {
	// Load templates
	tmpl, err := theme.LoadTemplates(b.config.ThemeName)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}
	b.templates = tmpl

	// Load all summaries
	summaries, err := b.loadSummaries()
	if err != nil {
		return nil, fmt.Errorf("failed to load summaries: %w", err)
	}

	// Load all digests
	digests, err := b.loadDigests()
	if err != nil {
		return nil, fmt.Errorf("failed to load digests: %w", err)
	}

	// Sort summaries by date (newest first)
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Meta.ArticleDate.After(summaries[j].Meta.ArticleDate)
	})

	// Build tag index
	tagIndex := buildTagIndex(summaries)

	// Create output directories
	dirs := []string{
		b.config.OutputDir,
		filepath.Join(b.config.OutputDir, "articles"),
		filepath.Join(b.config.OutputDir, "digests"),
		filepath.Join(b.config.OutputDir, "tags"),
		filepath.Join(b.config.OutputDir, "css"),
		filepath.Join(b.config.OutputDir, "js"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Write CSS and JS assets
	if err := theme.WriteAssets(b.config.OutputDir); err != nil {
		return nil, fmt.Errorf("failed to write assets: %w", err)
	}

	pageCount := 0

	// Generate index page
	if err := b.generateIndex(summaries, digests, tagIndex); err != nil {
		return nil, fmt.Errorf("failed to generate index: %w", err)
	}
	pageCount++

	// Generate article pages
	for _, summary := range summaries {
		if err := b.generateArticle(summary); err != nil {
			return nil, fmt.Errorf("failed to generate article %s: %w", summary.Article.Title, err)
		}
		pageCount++
	}

	// Generate digest pages
	for _, digest := range digests {
		if err := b.generateDigest(digest); err != nil {
			return nil, fmt.Errorf("failed to generate digest: %w", err)
		}
		pageCount++
	}

	// Generate tag pages
	for _, tag := range tagIndex {
		if err := b.generateTag(tag, summaries); err != nil {
			return nil, fmt.Errorf("failed to generate tag %s: %w", tag.Tag, err)
		}
		pageCount++
	}

	return &BuildResult{
		ArticleCount: len(summaries),
		DigestCount:  len(digests),
		PageCount:    pageCount,
	}, nil
}

func (b *Builder) loadSummaries() ([]types.Summary, error) {
	var summaries []types.Summary

	err := filepath.Walk(b.config.SummariesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		var summary types.Summary
		if ext == ".yaml" || ext == ".yml" {
			if err := yaml.Unmarshal(data, &summary); err != nil {
				return fmt.Errorf("failed to parse %s: %w", path, err)
			}
		} else {
			if err := json.Unmarshal(data, &summary); err != nil {
				return fmt.Errorf("failed to parse %s: %w", path, err)
			}
		}

		summaries = append(summaries, summary)
		return nil
	})

	return summaries, err
}

func (b *Builder) loadDigests() ([]types.Digest, error) {
	var digests []types.Digest

	if _, err := os.Stat(b.config.DigestsDir); os.IsNotExist(err) {
		return digests, nil
	}

	err := filepath.Walk(b.config.DigestsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		var digest types.Digest
		if ext == ".yaml" || ext == ".yml" {
			if err := yaml.Unmarshal(data, &digest); err != nil {
				return fmt.Errorf("failed to parse %s: %w", path, err)
			}
		} else {
			if err := json.Unmarshal(data, &digest); err != nil {
				return fmt.Errorf("failed to parse %s: %w", path, err)
			}
		}

		digests = append(digests, digest)
		return nil
	})

	return digests, err
}

func buildTagIndex(summaries []types.Summary) []types.TagIndex {
	tagMap := make(map[string][]types.ArticleReference)

	for _, summary := range summaries {
		ref := types.ArticleReference{
			Title:         summary.Article.Title,
			ArticleURL:    summary.Meta.ArticleURL,
			DiscussionURL: summary.Meta.DiscussionURL,
			Platform:      summary.Meta.Platform,
			ArticleDate:   summary.Meta.ArticleDate,
			Tags:          summary.Article.Tags,
			Sentiment:     summary.Article.Sentiment,
			CommentCount:  summary.Meta.DiscussionCommentCount,
			Score:         summary.Meta.Score,
			SummaryPath:   generateArticleSlug(summary),
		}

		for _, tag := range summary.Article.Tags {
			tagMap[tag] = append(tagMap[tag], ref)
		}
	}

	var index []types.TagIndex
	for tag, articles := range tagMap {
		index = append(index, types.TagIndex{
			Tag:      tag,
			Count:    len(articles),
			Articles: articles,
		})
	}

	// Sort by count (most articles first)
	sort.Slice(index, func(i, j int) bool {
		return index[i].Count > index[j].Count
	})

	return index
}

func generateArticleSlug(summary types.Summary) string {
	// Generate slug from date and title
	date := summary.Meta.ArticleDate.Format("2006-01-02")
	title := strings.ToLower(summary.Article.Title)
	title = strings.ReplaceAll(title, " ", "-")
	title = strings.ReplaceAll(title, "'", "")
	title = strings.ReplaceAll(title, "\"", "")
	title = strings.ReplaceAll(title, "?", "")
	title = strings.ReplaceAll(title, "!", "")
	title = strings.ReplaceAll(title, ":", "")
	title = strings.ReplaceAll(title, ",", "")
	title = strings.ReplaceAll(title, ".", "")
	title = strings.ReplaceAll(title, "--", "-")

	// Truncate to reasonable length
	if len(title) > 50 {
		title = title[:50]
	}
	title = strings.TrimSuffix(title, "-")

	return fmt.Sprintf("%s-%s", date, title)
}

// PageData contains common data for all templates.
type PageData struct {
	SiteTitle       string
	SiteDescription string
	BaseURL         string
	PageTitle       string
	CurrentYear     int
	GeneratedAt     time.Time
}

// IndexData contains data for the index page.
type IndexData struct {
	PageData
	Articles    []types.Summary
	Digests     []types.Digest
	TagIndex    []types.TagIndex
	TagCloud    []TagCloudItem
	TotalStats  SiteStats
	RecentCount int
}

// TagCloudItem represents a tag with weight for display.
type TagCloudItem struct {
	Tag    string
	Count  int
	Weight int // 1-5 for CSS sizing
}

// SiteStats contains aggregate statistics.
type SiteStats struct {
	TotalArticles      int
	TotalComments      int
	Platforms          []string
	DateRange          string
	SentimentBreakdown map[string]int
}

// ArticleData contains data for article pages.
type ArticleData struct {
	PageData
	Summary types.Summary
	Slug    string
}

// DigestData contains data for digest pages.
type DigestData struct {
	PageData
	Digest types.Digest
}

// TagData contains data for tag pages.
type TagData struct {
	PageData
	Tag      string
	Articles []types.Summary
	Count    int
}

func (b *Builder) generateIndex(summaries []types.Summary, digests []types.Digest, tagIndex []types.TagIndex) error {
	// Build tag cloud with weights
	var tagCloud []TagCloudItem
	maxCount := 0
	for _, t := range tagIndex {
		if t.Count > maxCount {
			maxCount = t.Count
		}
	}
	for _, t := range tagIndex {
		weight := 1
		if maxCount > 0 {
			weight = (t.Count * 4 / maxCount) + 1
			if weight > 5 {
				weight = 5
			}
		}
		tagCloud = append(tagCloud, TagCloudItem{
			Tag:    t.Tag,
			Count:  t.Count,
			Weight: weight,
		})
	}

	// Calculate stats
	totalComments := 0
	platformSet := make(map[string]bool)
	sentimentBreakdown := make(map[string]int)
	var minDate, maxDate time.Time

	for _, s := range summaries {
		totalComments += s.Meta.DiscussionCommentCount
		platformSet[string(s.Meta.Platform)] = true
		sentimentBreakdown[string(s.Article.Sentiment)]++

		if minDate.IsZero() || s.Meta.ArticleDate.Before(minDate) {
			minDate = s.Meta.ArticleDate
		}
		if maxDate.IsZero() || s.Meta.ArticleDate.After(maxDate) {
			maxDate = s.Meta.ArticleDate
		}
	}

	var platforms []string
	for p := range platformSet {
		platforms = append(platforms, p)
	}

	dateRange := ""
	if !minDate.IsZero() && !maxDate.IsZero() {
		dateRange = fmt.Sprintf("%s - %s", minDate.Format("Jan 2, 2006"), maxDate.Format("Jan 2, 2006"))
	}

	recentCount := 10
	if len(summaries) < recentCount {
		recentCount = len(summaries)
	}

	data := IndexData{
		PageData: PageData{
			SiteTitle:       b.config.SiteTitle,
			SiteDescription: b.config.SiteDescription,
			BaseURL:         b.config.BaseURL,
			PageTitle:       b.config.SiteTitle,
			CurrentYear:     time.Now().Year(),
			GeneratedAt:     time.Now(),
		},
		Articles:    summaries,
		Digests:     digests,
		TagIndex:    tagIndex,
		TagCloud:    tagCloud,
		RecentCount: recentCount,
		TotalStats: SiteStats{
			TotalArticles:      len(summaries),
			TotalComments:      totalComments,
			Platforms:          platforms,
			DateRange:          dateRange,
			SentimentBreakdown: sentimentBreakdown,
		},
	}

	return b.writeTemplate("index.html", filepath.Join(b.config.OutputDir, "index.html"), data)
}

func (b *Builder) generateArticle(summary types.Summary) error {
	slug := generateArticleSlug(summary)

	data := ArticleData{
		PageData: PageData{
			SiteTitle:       b.config.SiteTitle,
			SiteDescription: b.config.SiteDescription,
			BaseURL:         b.config.BaseURL,
			PageTitle:       summary.Article.Title,
			CurrentYear:     time.Now().Year(),
			GeneratedAt:     time.Now(),
		},
		Summary: summary,
		Slug:    slug,
	}

	outputPath := filepath.Join(b.config.OutputDir, "articles", slug+".html")
	return b.writeTemplate("article.html", outputPath, data)
}

func (b *Builder) generateDigest(digest types.Digest) error {
	// Generate slug from period
	slug := fmt.Sprintf("%s-%s",
		strings.ToLower(string(digest.Meta.Period)),
		digest.Meta.PeriodStart.Format("2006-01-02"))

	data := DigestData{
		PageData: PageData{
			SiteTitle:       b.config.SiteTitle,
			SiteDescription: b.config.SiteDescription,
			BaseURL:         b.config.BaseURL,
			PageTitle:       fmt.Sprintf("%s Digest", strings.Title(string(digest.Meta.Period))),
			CurrentYear:     time.Now().Year(),
			GeneratedAt:     time.Now(),
		},
		Digest: digest,
	}

	outputPath := filepath.Join(b.config.OutputDir, "digests", slug+".html")
	return b.writeTemplate("digest.html", outputPath, data)
}

func (b *Builder) generateTag(tagIndex types.TagIndex, allSummaries []types.Summary) error {
	// Filter summaries by tag
	var articles []types.Summary
	for _, s := range allSummaries {
		for _, t := range s.Article.Tags {
			if t == tagIndex.Tag {
				articles = append(articles, s)
				break
			}
		}
	}

	data := TagData{
		PageData: PageData{
			SiteTitle:       b.config.SiteTitle,
			SiteDescription: b.config.SiteDescription,
			BaseURL:         b.config.BaseURL,
			PageTitle:       fmt.Sprintf("Tag: %s", tagIndex.Tag),
			CurrentYear:     time.Now().Year(),
			GeneratedAt:     time.Now(),
		},
		Tag:      tagIndex.Tag,
		Articles: articles,
		Count:    len(articles),
	}

	slug := strings.ToLower(strings.ReplaceAll(tagIndex.Tag, " ", "-"))
	outputPath := filepath.Join(b.config.OutputDir, "tags", slug+".html")
	return b.writeTemplate("tag.html", outputPath, data)
}

func (b *Builder) writeTemplate(name, outputPath string, data any) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return b.templates.ExecuteTemplate(f, name, data)
}
