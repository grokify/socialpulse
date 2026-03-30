# Writing Personas

Personas are the heart of discussion analysis. They categorize commenters into archetypal viewpoints, making complex discussions digestible.

## What Makes a Good Persona?

A well-crafted persona:

1. **Represents a distinct viewpoint** - Not just "people who agree"
2. **Has internal consistency** - All members share core beliefs
3. **Provides insight** - Reveals something about the discussion
4. **Includes real voices** - Quotes from actual commenters

## Persona Structure

```yaml
personas:
  - name: "The Pragmatic Senior"
    description: "Experienced developers who've seen hype cycles before"
    stance: nuanced
    prevalence: dominant
    percentage: 35
    core_argument: "AI tools are useful when applied deliberately"
    quotes:
      - text: "I've been coding for 20 years. Every new tool promises revolution."
        author: "veteran_dev"
      - text: "The key is knowing when AI helps and when it hurts."
        author: "tech_lead_42"
```

## Naming Personas

Good persona names are:

- **Memorable** - Easy to recall and reference
- **Descriptive** - Suggest the viewpoint
- **Neutral** - Don't prejudge the position

### Examples

| Good | Avoid |
|------|-------|
| The Pragmatic Senior | The Old Guard |
| The Skeptic | The Naysayer |
| The Enthusiast | The Fanboy |
| The Newcomer | The Clueless Beginner |

### Common Archetypes

| Name | Description |
|------|-------------|
| The Pragmatist | Focuses on practical applicability |
| The Skeptic | Demands evidence, questions claims |
| The Enthusiast | Excited about possibilities |
| The Contrarian | Challenges consensus |
| The Historian | References past experiences |
| The Theorist | Focuses on principles and frameworks |
| The Builder | Shares implementation experience |

## Writing Descriptions

Descriptions should be 50-100 characters and capture the essence:

**Good:**
> Experienced devs who've seen hype cycles before and approach claims skeptically

**Too vague:**
> People who have opinions about the topic

**Too specific:**
> Senior engineers at FAANG companies with 15+ years experience

## Stance Values

| Stance | When to Use |
|--------|-------------|
| `agrees` | Supports the article's thesis |
| `disagrees` | Opposes the article's thesis |
| `nuanced` | Sees validity in multiple positions |
| `tangential` | Focuses on related but different topics |

### Nuanced Stance

Most valuable discussions have `nuanced` personas. They:

- Acknowledge merit in the article
- Raise valid concerns or limitations
- Add context the author missed
- Bridge different viewpoints

## Prevalence

Prevalence indicates how common this viewpoint is:

| Value | Percentage Range | Use When |
|-------|------------------|----------|
| `dominant` | 30-50% | Most common viewpoint |
| `significant` | 15-30% | Substantial presence |
| `minority` | 5-15% | Notable but uncommon |

### Percentage Field

The `percentage` field should:

- Sum to ~100% across all personas
- Match the `prevalence` category
- Be an estimate, not exact count

Example distribution:

```yaml
personas:
  - name: "The Pragmatist"
    prevalence: dominant
    percentage: 35
  - name: "The Skeptic"
    prevalence: significant
    percentage: 25
  - name: "The Enthusiast"
    prevalence: significant
    percentage: 20
  - name: "The Theorist"
    prevalence: minority
    percentage: 12
  - name: "The Contrarian"
    prevalence: minority
    percentage: 8
```

## Writing Core Arguments

The `core_argument` distills the persona's position to 100-150 characters:

**Good:**
> AI tools are powerful when used deliberately, but dangerous when used as a crutch for understanding

**Too long:**
> While AI coding tools certainly have their place in modern development workflows, the key insight is that they work best when developers maintain a clear understanding of what the AI is producing and why...

**Too short:**
> AI is good sometimes

## Selecting Quotes

Quotes should:

1. **Be verbatim** - Use `fetch-quotes` to verify
2. **Represent the persona** - Clearly express the viewpoint
3. **Be substantive** - Not just "I agree" or "This"
4. **Be readable** - Avoid excessive jargon or context-dependent text

### Quote Selection Process

1. Identify comments that match the persona
2. Look for articulate expressions of the viewpoint
3. Choose 1-3 quotes per persona
4. Verify with `socialpulse fetch-quotes --dry-run`

### Cleaning Quotes

Remove quoted parent text (lines starting with `>`):

**Original:**
```
> The author claims AI is dangerous

I disagree. In my experience, AI tools have been incredibly helpful
for routine tasks while freeing me to focus on architecture.
```

**Cleaned:**
```
I disagree. In my experience, AI tools have been incredibly helpful
for routine tasks while freeing me to focus on architecture.
```

## How Many Personas?

| Discussion Size | Recommended Personas |
|-----------------|---------------------|
| <50 comments | 2-3 |
| 50-200 comments | 3-5 |
| 200+ comments | 4-6 |

More personas aren't always better. Each should add distinct insight.

## Common Mistakes

### 1. Overlapping Personas

**Problem:** Two personas that are essentially the same viewpoint

```yaml
# Bad - these are the same
- name: "The Skeptic"
  core_argument: "The evidence is weak"
- name: "The Critic"
  core_argument: "The claims aren't supported"
```

**Solution:** Merge into one persona or find distinct angles

### 2. Missing the Main Debate

**Problem:** Personas don't capture the core disagreement

**Solution:** Start by identifying what people are arguing about, then build personas around those positions

### 3. Strawman Personas

**Problem:** Unfairly characterizing a position

**Solution:** Use actual quotes, let commenters speak for themselves

### 4. Too Many Small Personas

**Problem:** Five personas each with 10% prevalence

**Solution:** Focus on dominant viewpoints, merge similar minorities
