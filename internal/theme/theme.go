// Package theme implements the template and asset system.
package theme

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/grokify/socialpulse/types"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed assets/*.css assets/*.js
var assetsFS embed.FS

// LoadTemplates loads all templates for the given theme.
func LoadTemplates(themeName string) (*template.Template, error) {
	// Define template functions
	funcMap := template.FuncMap{
		"formatDate":      formatDate,
		"formatDateTime":  formatDateTime,
		"truncate":        truncate,
		"slugify":         slugify,
		"articleSlug":     articleSlug,
		"sentimentClass":  sentimentClass,
		"stanceClass":     stanceClass,
		"prevalenceClass": prevalenceClass,
		"platformIcon":    platformIcon,
		"join":            strings.Join,
		"upper":           strings.ToUpper,
		"lower":           strings.ToLower,
		"title":           strings.Title,
		"add":             func(a, b int) int { return a + b },
		"sub":             func(a, b int) int { return a - b },
		"mul":             func(a, b int) int { return a * b },
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a % b
		},
		"seq":      seq,
		"safeHTML": func(s string) template.HTML { return template.HTML(s) }, //nolint:gosec // G203: intentional raw HTML for template rendering
	}

	// For now, only support default theme from embedded files
	tmpl := template.New("").Funcs(funcMap)

	// Parse all embedded templates
	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return nil, fmt.Errorf("failed to read templates: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}

		content, err := templatesFS.ReadFile(filepath.Join("templates", entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", entry.Name(), err)
		}

		_, err = tmpl.New(entry.Name()).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", entry.Name(), err)
		}
	}

	return tmpl, nil
}

// WriteAssets writes CSS and JS assets to the output directory.
func WriteAssets(outputDir string) error {
	// Create asset directories
	cssDir := filepath.Join(outputDir, "css")
	jsDir := filepath.Join(outputDir, "js")

	if err := os.MkdirAll(cssDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(jsDir, 0755); err != nil {
		return err
	}

	// Copy CSS files
	entries, err := assetsFS.ReadDir("assets")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		content, err := assetsFS.ReadFile(filepath.Join("assets", entry.Name()))
		if err != nil {
			return err
		}

		var destDir string
		if strings.HasSuffix(entry.Name(), ".css") {
			destDir = cssDir
		} else if strings.HasSuffix(entry.Name(), ".js") {
			destDir = jsDir
		} else {
			continue
		}

		destPath := filepath.Join(destDir, entry.Name())
		if err := os.WriteFile(destPath, content, 0644); err != nil { //nolint:gosec // G306: public web assets need 0644
			return err
		}
	}

	return nil
}

// Template helper functions

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("Jan 2, 2006")
}

func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("Jan 2, 2006 3:04 PM")
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "!", "")
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "--", "-")
	return strings.TrimSuffix(s, "-")
}

func articleSlug(summary types.Summary) string {
	date := summary.Meta.ArticleDate.Format("2006-01-02")
	title := slugify(summary.Article.Title)
	if len(title) > 50 {
		title = title[:50]
	}
	title = strings.TrimSuffix(title, "-")
	return fmt.Sprintf("%s-%s", date, title)
}

func sentimentClass(sentiment types.ArticleSentiment) string {
	switch sentiment {
	case types.SentimentOptimistic:
		return "sentiment-optimistic"
	case types.SentimentCautious:
		return "sentiment-cautious"
	case types.SentimentPessimistic:
		return "sentiment-pessimistic"
	case types.SentimentNeutral:
		return "sentiment-neutral"
	case types.SentimentProvocative:
		return "sentiment-provocative"
	default:
		return "sentiment-neutral"
	}
}

func stanceClass(stance types.Stance) string {
	switch stance {
	case types.StanceAgrees:
		return "stance-agrees"
	case types.StanceDisagrees:
		return "stance-disagrees"
	case types.StanceNuanced:
		return "stance-nuanced"
	case types.StanceTangential:
		return "stance-tangential"
	default:
		return "stance-nuanced"
	}
}

func prevalenceClass(prevalence types.Prevalence) string {
	switch prevalence {
	case types.PrevalenceDominant:
		return "prevalence-dominant"
	case types.PrevalenceSignificant:
		return "prevalence-significant"
	case types.PrevalenceMinority:
		return "prevalence-minority"
	default:
		return "prevalence-significant"
	}
}

func platformIcon(platform types.Platform) string {
	switch platform {
	case types.PlatformHackerNews:
		return "HN"
	case types.PlatformReddit:
		return "R"
	default:
		return "?"
	}
}

func seq(start, end int) []int {
	var result []int
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}
