package textutil

import (
	"testing"
)

func TestStripHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "plain text",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "paragraph tags",
			input:    "<p>First paragraph</p><p>Second paragraph</p>",
			expected: "First paragraph Second paragraph",
		},
		{
			name:     "br tags",
			input:    "Line one<br>Line two<br/>Line three",
			expected: "Line one Line two Line three",
		},
		{
			name:     "inline formatting",
			input:    "This is <b>bold</b> and <i>italic</i> and <code>code</code>",
			expected: "This is bold and italic and code",
		},
		{
			name:     "anchor tags",
			input:    "Check out <a href=\"https://example.com\">this link</a> here",
			expected: "Check out this link here",
		},
		{
			name:     "HTML entities",
			input:    "Tom &amp; Jerry &lt;3 &gt; rock &quot;forever&quot;",
			expected: "Tom & Jerry <3 > rock \"forever\"",
		},
		{
			name:     "slash entity",
			input:    "CI&#x2F;CD pipeline",
			expected: "CI/CD pipeline",
		},
		{
			name:     "apostrophe entities",
			input:    "It&#x27;s working &#39;fine&#39;",
			expected: "It's working 'fine'",
		},
		{
			name:     "multiple whitespace",
			input:    "Too    many   spaces",
			expected: "Too many spaces",
		},
		{
			name:     "nested tags",
			input:    "<p><b>Bold in <i>paragraph</i></b></p>",
			expected: "Bold in paragraph",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripHTML(tt.input)
			if result != tt.expected {
				t.Errorf("StripHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveQuotedLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no quotes",
			input:    "Just a normal comment",
			expected: "Just a normal comment",
		},
		{
			name:     "single quoted line",
			input:    "> This is quoted\nThis is my response",
			expected: "This is my response",
		},
		{
			name:     "multiple quoted lines",
			input:    "> First quote\n> Second quote\nMy reply here",
			expected: "My reply here",
		},
		{
			name:     "quote in middle",
			input:    "Start here\n> Quoted part\nEnd here",
			expected: "Start here End here",
		},
		{
			name:     "nested quotes",
			input:    ">> Double quoted\n> Single quoted\nActual content",
			expected: "Actual content",
		},
		{
			name:     "greater than in text",
			input:    "5 > 3 is true",
			expected: "5 > 3 is true",
		},
		{
			name:     "whitespace before quote marker",
			input:    "   > This should be removed\nKeep this",
			expected: "Keep this",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveQuotedLines(tt.input)
			if result != tt.expected {
				t.Errorf("RemoveQuotedLines(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStripHTMLAndQuotes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HTML with quotes",
			input:    "<p>> Quoted text</p><p>My response here</p>",
			expected: "My response here",
		},
		{
			name:     "complex HN comment",
			input:    "<p>> The original point was about testing</p><p>I disagree. Testing is &lt;important&gt; for CI&#x2F;CD.</p>",
			expected: "I disagree. Testing is <important> for CI/CD.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripHTMLAndQuotes(tt.input)
			if result != tt.expected {
				t.Errorf("StripHTMLAndQuotes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWordSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected float64
	}{
		{
			name:     "identical strings",
			a:        "hello world",
			b:        "hello world",
			expected: 1.0,
		},
		{
			name:     "no overlap",
			a:        "hello world",
			b:        "foo bar",
			expected: 0.0,
		},
		{
			name:     "partial overlap",
			a:        "hello world",
			b:        "hello there",
			expected: 0.5,
		},
		{
			name:     "empty first",
			a:        "",
			b:        "hello",
			expected: 0.0,
		},
		{
			name:     "empty second",
			a:        "hello",
			b:        "",
			expected: 0.0,
		},
		{
			name:     "case insensitive",
			a:        "Hello World",
			b:        "HELLO WORLD",
			expected: 1.0,
		},
		{
			name:     "subset",
			a:        "one two three",
			b:        "one two",
			expected: 2.0 / 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WordSimilarity(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("WordSimilarity(%q, %q) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

//nolint:dupl // test table structure is intentionally similar
func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "no truncation needed",
			input:    "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "exact length",
			input:    "exact",
			maxLen:   5,
			expected: "exact",
		},
		{
			name:     "truncation with ellipsis",
			input:    "this is a long string",
			maxLen:   10,
			expected: "this is...",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   10,
			expected: "",
		},
		{
			name:     "very short maxLen",
			input:    "hello",
			maxLen:   3,
			expected: "hel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

//nolint:dupl // test table structure is intentionally similar
func TestTruncateAtSentence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "no truncation needed",
			input:    "Short text.",
			maxLen:   50,
			expected: "Short text.",
		},
		{
			name:     "truncate at period",
			input:    "First sentence. Second sentence. Third sentence.",
			maxLen:   30,
			expected: "First sentence.",
		},
		{
			name:     "truncate at question mark",
			input:    "Is this working? I hope so.",
			maxLen:   20,
			expected: "Is this working?",
		},
		{
			name:     "truncate at exclamation",
			input:    "Wow! That is amazing!",
			maxLen:   10,
			expected: "Wow!",
		},
		{
			name:     "no sentence break",
			input:    "This is a long string without any sentence breaks at all",
			maxLen:   30,
			expected: "This is a long string witho...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateAtSentence(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateAtSentence(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		minLength int
		expected  []string
	}{
		{
			name:      "basic extraction",
			input:     "The quick brown fox jumps",
			minLength: 3,
			expected:  []string{"quick", "brown", "fox", "jumps"},
		},
		{
			name:      "filters stop words",
			input:     "this is a test of the system",
			minLength: 3,
			expected:  []string{"test", "system"},
		},
		{
			name:      "respects min length",
			input:     "go is fun to use",
			minLength: 3,
			expected:  []string{"fun", "use"},
		},
		{
			name:      "removes duplicates",
			input:     "test test test different test",
			minLength: 3,
			expected:  []string{"test", "different"},
		},
		{
			name:      "handles empty",
			input:     "",
			minLength: 3,
			expected:  nil,
		},
		{
			name:      "case insensitive",
			input:     "Hello WORLD hello",
			minLength: 3,
			expected:  []string{"hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractKeywords(tt.input, tt.minLength)
			if !stringSliceEqual(result, tt.expected) {
				t.Errorf("ExtractKeywords(%q, %d) = %v, want %v", tt.input, tt.minLength, result, tt.expected)
			}
		})
	}
}

func TestScoreText(t *testing.T) {
	tests := []struct {
		name            string
		text            string
		keywords        []string
		expectedScore   float64
		expectedMatches []string
	}{
		{
			name:            "no matches",
			text:            "nothing here",
			keywords:        []string{"foo", "bar"},
			expectedScore:   0.0,
			expectedMatches: nil,
		},
		{
			name:            "single match",
			text:            "the testing framework",
			keywords:        []string{"testing"},
			expectedScore:   7.0 / 5.0, // len("testing") / 5.0
			expectedMatches: []string{"testing"},
		},
		{
			name:            "multiple matches",
			text:            "testing and debugging code",
			keywords:        []string{"testing", "code"},
			expectedScore:   (7.0 + 4.0) / 5.0,
			expectedMatches: []string{"testing", "code"},
		},
		{
			name:            "case insensitive",
			text:            "TESTING the CODE",
			keywords:        []string{"testing", "code"},
			expectedScore:   (7.0 + 4.0) / 5.0,
			expectedMatches: []string{"testing", "code"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScoreText(tt.text, tt.keywords)
			if result.Score != tt.expectedScore {
				t.Errorf("ScoreText score = %v, want %v", result.Score, tt.expectedScore)
			}
			if !stringSliceEqual(result.Matches, tt.expectedMatches) {
				t.Errorf("ScoreText matches = %v, want %v", result.Matches, tt.expectedMatches)
			}
		})
	}
}

func TestScoreComment(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		keywords      []string
		minScore      float64
		maxScore      float64
		expectMatches bool
	}{
		{
			name:          "short comment penalized",
			text:          "yes agree",
			keywords:      []string{"agree"},
			minScore:      0.0,
			maxScore:      1.0, // Should be penalized for being short
			expectMatches: true,
		},
		{
			name:          "removes quoted lines",
			text:          "> This is quoted\nMy actual response with testing keywords",
			keywords:      []string{"testing", "keywords"},
			minScore:      0.0,
			maxScore:      5.0,
			expectMatches: true,
		},
		{
			name:          "reasonable length bonus",
			text:          "This is a reasonable length comment that contains testing keywords and should get a bonus because it has enough words to be considered substantive content",
			keywords:      []string{"testing", "keywords", "content"},
			minScore:      3.0, // Should get 1.2x bonus
			maxScore:      10.0,
			expectMatches: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScoreComment(tt.text, tt.keywords)
			if result.Score < tt.minScore || result.Score > tt.maxScore {
				t.Errorf("ScoreComment score = %v, want between %v and %v", result.Score, tt.minScore, tt.maxScore)
			}
			if tt.expectMatches && len(result.Matches) == 0 {
				t.Error("ScoreComment expected matches but got none")
			}
		})
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
