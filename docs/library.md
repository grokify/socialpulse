# Library Usage

SocialPulse can be used as a Go library for building custom tools.

## Installation

```bash
go get github.com/grokify/socialpulse
```

## Packages

| Package | Description |
|---------|-------------|
| `types` | Core data structures |
| `schema` | Embedded JSON schemas |
| `textutil` | Text processing utilities |

## Types Package

The `types` package defines all data structures:

```go
import "github.com/grokify/socialpulse/types"
```

### Creating a Summary

```go
summary := types.Summary{
    SchemaVersion: "1.0.0",
    Meta: types.Meta{
        ArticleURL:             "https://example.com/article",
        DiscussionURL:          "https://news.ycombinator.com/item?id=12345",
        Platform:               types.PlatformHackerNews,
        ArticleDate:            time.Now(),
        SummaryDate:            time.Now(),
        ArticleAuthor:          "Jane Doe",
        DiscussionCommentCount: 150,
        ItemID:                 "12345",
    },
    Article: types.Article{
        Title:     "Article Title",
        Thesis:    "The main argument...",
        KeyArguments: []string{
            "First point",
            "Second point",
        },
        Tags:      []string{"technology", "programming"},
        Sentiment: types.SentimentCautious,
    },
    Discussion: types.Discussion{
        Sentiment: types.DiscussionDivided,
        Personas: []types.Persona{
            {
                Name:        "The Pragmatist",
                Description: "Experienced practitioners",
                Stance:      types.StanceNuanced,
                Prevalence:  types.PrevalenceDominant,
                Percentage:  35,
                CoreArgument: "Good in theory, execution matters",
                Quotes: []types.Quote{
                    {
                        Text:   "In my experience...",
                        Author: "dev_user",
                    },
                },
            },
        },
    },
}
```

### Platform Constants

```go
types.PlatformHackerNews  // "hackernews"
types.PlatformReddit      // "reddit"
```

### Sentiment Constants

```go
// Article sentiment
types.SentimentOptimistic
types.SentimentCautious
types.SentimentPessimistic
types.SentimentNeutral
types.SentimentProvocative

// Discussion sentiment
types.DiscussionSupportive
types.DiscussionDivided
types.DiscussionCritical
types.DiscussionExploratory
```

### Stance and Prevalence

```go
// Stance
types.StanceAgrees
types.StanceDisagrees
types.StanceNuanced
types.StanceTangential

// Prevalence
types.PrevalenceDominant
types.PrevalenceSignificant
types.PrevalenceMinority
```

## Schema Package

The `schema` package provides embedded JSON schemas:

```go
import "github.com/grokify/socialpulse/schema"
```

### Get Schema JSON

```go
// Get raw JSON bytes
summaryJSON := schema.SummarySchemaJSON()
digestJSON := schema.DigestSchemaJSON()
siteJSON := schema.SiteSchemaJSON()

// Parse as map
summaryMap, err := schema.SummarySchemaMap()
if err != nil {
    log.Fatal(err)
}
```

### Validate Content

```go
import (
    "github.com/grokify/socialpulse/schema"
    "github.com/xeipuuv/gojsonschema"
)

func validateSummary(summaryJSON []byte) error {
    schemaLoader := gojsonschema.NewBytesLoader(schema.SummarySchemaJSON())
    documentLoader := gojsonschema.NewBytesLoader(summaryJSON)

    result, err := gojsonschema.Validate(schemaLoader, documentLoader)
    if err != nil {
        return err
    }

    if !result.Valid() {
        for _, err := range result.Errors() {
            fmt.Printf("- %s\n", err)
        }
        return fmt.Errorf("validation failed")
    }

    return nil
}
```

## TextUtil Package

The `textutil` package provides text processing utilities:

```go
import "github.com/grokify/socialpulse/textutil"
```

### HTML Processing

```go
// Strip HTML tags and decode entities
plainText := textutil.StripHTML("<p>Hello &amp; world</p>")
// Result: "Hello & world"

// Remove quoted lines (lines starting with >)
cleaned := textutil.RemoveQuotedLines("> quoted\nmy text")
// Result: "my text"

// Combined: strip HTML and remove quotes
text := textutil.StripHTMLAndQuotes("<p>> quoted</p><p>actual content</p>")
// Result: "actual content"
```

### Text Truncation

```go
// Simple truncation with ellipsis
short := textutil.Truncate("This is a long string", 10)
// Result: "This is..."

// Truncate at sentence boundary
sentence := textutil.TruncateAtSentence("First sentence. Second sentence.", 20)
// Result: "First sentence."
```

### Keyword Extraction

```go
// Extract keywords (min length 3)
keywords := textutil.ExtractKeywords("The quick brown fox jumps", 3)
// Result: ["quick", "brown", "fox", "jumps"]

// With custom stop words
stopWords := map[string]bool{"quick": true}
keywords := textutil.ExtractKeywordsWithStopWords(text, 3, stopWords)
```

### Text Scoring

```go
// Score text against keywords
result := textutil.ScoreText("Testing the code quality", []string{"testing", "quality"})
fmt.Printf("Score: %.2f, Matches: %v\n", result.Score, result.Matches)

// Score comment (with length adjustments)
result := textutil.ScoreComment(commentText, keywords)
// Applies bonuses/penalties based on comment length
```

### Word Similarity

```go
// Compute word overlap similarity (0.0 to 1.0)
similarity := textutil.WordSimilarity("hello world", "hello there")
// Result: 0.5 (50% overlap)
```

## Example: Custom Summary Generator

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/grokify/socialpulse/types"
)

func main() {
    summary := types.Summary{
        SchemaVersion: "1.0.0",
        Meta: types.Meta{
            ArticleURL:             "https://example.com/article",
            DiscussionURL:          "https://news.ycombinator.com/item?id=12345",
            Platform:               types.PlatformHackerNews,
            ArticleDate:            time.Now(),
            SummaryDate:            time.Now(),
            ArticleAuthor:          "Author",
            DiscussionCommentCount: 100,
            ItemID:                 "12345",
        },
        Article: types.Article{
            Title:        "My Article",
            Thesis:       "The main point of the article.",
            KeyArguments: []string{"Point 1", "Point 2"},
            Tags:         []string{"tech"},
            Sentiment:    types.SentimentNeutral,
        },
        Discussion: types.Discussion{
            Sentiment: types.DiscussionExploratory,
            Personas:  []types.Persona{},
        },
    }

    // Marshal to JSON
    data, err := json.MarshalIndent(summary, "", "  ")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println(string(data))
}
```

## Example: Quote Processor

```go
package main

import (
    "fmt"

    "github.com/grokify/socialpulse/textutil"
)

func processComment(htmlComment string, keywords []string) {
    // Clean the comment
    text := textutil.StripHTMLAndQuotes(htmlComment)

    // Score against keywords
    result := textutil.ScoreComment(text, keywords)

    if result.Score > 2.0 {
        fmt.Printf("High relevance comment (score: %.2f)\n", result.Score)
        fmt.Printf("Matched keywords: %v\n", result.Matches)
        fmt.Printf("Text: %s\n", textutil.Truncate(text, 200))
    }
}

func main() {
    comment := `<p>> Previous comment quoted here</p>
<p>I think this is a great point about testing and code quality.
The author makes compelling arguments.</p>`

    keywords := textutil.ExtractKeywords("testing code quality software", 3)
    processComment(comment, keywords)
}
```
