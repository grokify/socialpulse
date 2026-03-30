// Package textutil provides text processing utilities for social discussion analysis.
package textutil

import (
	"html"
	"regexp"
	"strings"
)

// StripHTML removes HTML tags and decodes HTML entities from text.
// It preserves paragraph breaks as spaces and removes anchor tags while keeping their text.
func StripHTML(s string) string {
	result := s

	// Convert block elements to newlines
	result = strings.ReplaceAll(result, "<p>", "\n")
	result = strings.ReplaceAll(result, "</p>", "")
	result = strings.ReplaceAll(result, "<br>", "\n")
	result = strings.ReplaceAll(result, "<br/>", "\n")
	result = strings.ReplaceAll(result, "<br />", "\n")

	// Remove inline formatting tags
	result = strings.ReplaceAll(result, "<i>", "")
	result = strings.ReplaceAll(result, "</i>", "")
	result = strings.ReplaceAll(result, "<b>", "")
	result = strings.ReplaceAll(result, "</b>", "")
	result = strings.ReplaceAll(result, "<em>", "")
	result = strings.ReplaceAll(result, "</em>", "")
	result = strings.ReplaceAll(result, "<strong>", "")
	result = strings.ReplaceAll(result, "</strong>", "")
	result = strings.ReplaceAll(result, "<code>", "")
	result = strings.ReplaceAll(result, "</code>", "")
	result = strings.ReplaceAll(result, "<pre>", "")
	result = strings.ReplaceAll(result, "</pre>", "")

	// Remove anchor tags but keep text
	for strings.Contains(result, "<a ") {
		start := strings.Index(result, "<a ")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	result = strings.ReplaceAll(result, "</a>", "")

	// Decode HTML entities
	result = html.UnescapeString(result)

	// Normalize whitespace
	result = strings.Join(strings.Fields(result), " ")

	return strings.TrimSpace(result)
}

// RemoveQuotedLines removes lines that start with ">" (quoted parent text in forums).
// This is common in HackerNews and Reddit where users quote parent comments.
func RemoveQuotedLines(s string) string {
	lines := strings.Split(s, "\n")
	var filtered []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip lines that start with ">" (quoted text from parent)
		if strings.HasPrefix(trimmed, ">") {
			continue
		}
		// Skip consecutive empty lines
		if len(trimmed) == 0 && len(filtered) > 0 && filtered[len(filtered)-1] == "" {
			continue
		}
		filtered = append(filtered, line)
	}

	// Join and normalize whitespace
	result := strings.Join(filtered, " ")
	result = strings.Join(strings.Fields(result), " ")
	return strings.TrimSpace(result)
}

// StripHTMLAndQuotes combines StripHTML and RemoveQuotedLines for processing forum comments.
func StripHTMLAndQuotes(s string) string {
	// First strip HTML to get plain text with newlines
	result := s

	// Convert block elements to newlines (before HTML unescape to preserve structure)
	result = strings.ReplaceAll(result, "<p>", "\n")
	result = strings.ReplaceAll(result, "</p>", "")
	result = strings.ReplaceAll(result, "<br>", "\n")
	result = strings.ReplaceAll(result, "<br/>", "\n")
	result = strings.ReplaceAll(result, "<br />", "\n")

	// Remove inline formatting tags
	result = strings.ReplaceAll(result, "<i>", "")
	result = strings.ReplaceAll(result, "</i>", "")
	result = strings.ReplaceAll(result, "<b>", "")
	result = strings.ReplaceAll(result, "</b>", "")
	result = strings.ReplaceAll(result, "<em>", "")
	result = strings.ReplaceAll(result, "</em>", "")
	result = strings.ReplaceAll(result, "<strong>", "")
	result = strings.ReplaceAll(result, "</strong>", "")
	result = strings.ReplaceAll(result, "<code>", "")
	result = strings.ReplaceAll(result, "</code>", "")
	result = strings.ReplaceAll(result, "<pre>", "")
	result = strings.ReplaceAll(result, "</pre>", "")

	// Remove anchor tags but keep text
	for strings.Contains(result, "<a ") {
		start := strings.Index(result, "<a ")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	result = strings.ReplaceAll(result, "</a>", "")

	// Decode HTML entities
	result = html.UnescapeString(result)

	// Remove quoted lines (must happen before whitespace normalization)
	lines := strings.Split(result, "\n")
	var filtered []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, ">") {
			continue
		}
		if trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}
	result = strings.Join(filtered, " ")

	// Normalize whitespace
	result = strings.Join(strings.Fields(result), " ")

	return strings.TrimSpace(result)
}

// WordSimilarity computes word overlap similarity between two strings.
// Returns a value between 0.0 (no overlap) and 1.0 (complete overlap).
// Uses a Jaccard-like coefficient based on word matching.
func WordSimilarity(a, b string) float64 {
	wordsA := strings.Fields(strings.ToLower(a))
	wordsB := strings.Fields(strings.ToLower(b))

	if len(wordsA) == 0 || len(wordsB) == 0 {
		return 0
	}

	// Build set of words from a
	wordSet := make(map[string]bool)
	for _, w := range wordsA {
		wordSet[w] = true
	}

	// Count matches
	matches := 0
	for _, w := range wordsB {
		if wordSet[w] {
			matches++
		}
	}

	// Return ratio of matches to words in a
	return float64(matches) / float64(len(wordsA))
}

// Truncate shortens a string to maxLen characters, adding "..." if truncated.
func Truncate(s string, maxLen int) string {
	if maxLen <= 3 {
		return s[:min(len(s), maxLen)]
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// TruncateAtSentence truncates text at a sentence boundary near maxLen.
// It looks for sentence-ending punctuation (.!?) within a range before maxLen.
func TruncateAtSentence(s string, maxLen int) string {
	if len(s) <= maxLen {
		return strings.TrimSpace(s)
	}

	// Look for sentence break in the last 50 chars before maxLen
	searchStart := maxLen - 50
	if searchStart < 0 {
		searchStart = 0
	}

	for i := maxLen - 1; i >= searchStart; i-- {
		if s[i] == '.' || s[i] == '!' || s[i] == '?' {
			return strings.TrimSpace(s[:i+1])
		}
	}

	// No sentence break found, truncate with ellipsis
	return strings.TrimSpace(s[:maxLen-3]) + "..."
}

// DefaultStopWords returns a set of common English stop words to filter from keyword extraction.
var DefaultStopWords = map[string]bool{
	"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
	"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
	"with": true, "by": true, "from": true, "as": true, "is": true, "was": true,
	"are": true, "were": true, "been": true, "be": true, "have": true, "has": true,
	"had": true, "do": true, "does": true, "did": true, "will": true, "would": true,
	"could": true, "should": true, "may": true, "might": true, "must": true,
	"that": true, "which": true, "who": true, "whom": true, "this": true,
	"these": true, "those": true, "it": true, "its": true, "they": true,
	"their": true, "them": true, "we": true, "our": true, "us": true,
	"i": true, "me": true, "my": true, "you": true, "your": true,
	"not": true, "no": true, "yes": true, "so": true, "if": true, "then": true,
	"than": true, "too": true, "very": true, "just": true, "only": true,
	"also": true, "about": true, "into": true, "over": true, "after": true,
}

var wordRegex = regexp.MustCompile(`[a-zA-Z]+`)

// ExtractKeywords extracts meaningful keywords from text, filtering stop words.
// Words shorter than minLength are excluded. Duplicates are removed.
func ExtractKeywords(text string, minLength int) []string {
	return ExtractKeywordsWithStopWords(text, minLength, DefaultStopWords)
}

// ExtractKeywordsWithStopWords extracts keywords using a custom stop word set.
func ExtractKeywordsWithStopWords(text string, minLength int, stopWords map[string]bool) []string {
	text = strings.ToLower(text)
	words := wordRegex.FindAllString(text, -1)

	var keywords []string
	seen := make(map[string]bool)

	for _, word := range words {
		if len(word) < minLength {
			continue
		}
		if stopWords[word] {
			continue
		}
		if seen[word] {
			continue
		}
		seen[word] = true
		keywords = append(keywords, word)
	}

	return keywords
}

// ScoreTextResult holds the result of scoring text against keywords.
type ScoreTextResult struct {
	Score   float64
	Matches []string
}

// ScoreText scores text against a list of keywords.
// Longer keywords contribute more to the score.
// Returns the score and list of matched keywords.
func ScoreText(text string, keywords []string) ScoreTextResult {
	textLower := strings.ToLower(text)

	var matches []string
	score := 0.0

	for _, kw := range keywords {
		if strings.Contains(textLower, kw) {
			matches = append(matches, kw)
			// Longer keywords are more specific, worth more
			score += float64(len(kw)) / 5.0
		}
	}

	return ScoreTextResult{
		Score:   score,
		Matches: matches,
	}
}

// ScoreComment scores a comment with length-based adjustments.
// Reasonable length comments (20-200 words) get a bonus.
// Very short (<10 words) or very long (>300 words) comments are penalized.
func ScoreComment(text string, keywords []string) ScoreTextResult {
	// Remove quoted text before scoring
	text = RemoveQuotedLines(text)

	result := ScoreText(text, keywords)

	// Adjust score based on comment length
	wordCount := len(strings.Fields(text))
	if wordCount >= 20 && wordCount <= 200 {
		result.Score *= 1.2 // Bonus for reasonable length
	} else if wordCount < 10 {
		result.Score *= 0.5 // Penalty for very short
	}
	if wordCount > 300 {
		result.Score *= 0.7 // Penalty for very long (likely code dumps)
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
