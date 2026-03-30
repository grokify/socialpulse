# Themes

SocialPulse uses a template-based theme system for generating HTML.

## Default Theme

The default theme is a data dashboard optimized for information density:

- **Sidebar navigation** - Quick access to sections
- **Metrics bar** - Key statistics at a glance
- **Card layouts** - Persona cards, article cards
- **ECharts integration** - Donut charts for persona distribution
- **Responsive design** - Works on mobile and desktop

## Theme Structure

The default theme includes:

```
internal/theme/
├── theme.go           # Theme embedding
├── assets/
│   ├── dashboard.css  # Main stylesheet
│   └── charts.js      # ECharts initialization
└── templates/
    ├── base.html      # Layout wrapper
    ├── index.html     # Dashboard home
    ├── article.html   # Article detail
    ├── digest.html    # Digest page
    └── tag.html       # Tag listing
```

## Page Types

### Index Page

The dashboard home page shows:

- Site title and description
- Metrics bar (article count, comment count, platforms, date range)
- Recent articles grid
- Tag cloud
- Sentiment distribution
- Digest links

### Article Page

Individual article summaries display:

- Article metadata (date, author, platform, comments)
- Thesis in highlighted box
- Key arguments list
- Examples section
- Persona cards with:
  - Name and percentage
  - Stance and prevalence badges
  - Core argument
  - Expandable quotes
- Tangents
- Consensus points
- Open questions
- ECharts donut for persona distribution

### Tag Page

Tag pages list all articles with that tag:

- Tag name and count
- Article cards sorted by date

### Digest Page

Periodic aggregations show:

- Date range
- Aggregated metrics
- Article list
- Combined persona analysis

## Styling

### CSS Classes

Key CSS classes for customization:

| Class | Element |
|-------|---------|
| `.sidebar` | Navigation sidebar |
| `.content` | Main content area |
| `.metrics-bar` | Statistics bar |
| `.article-card` | Article preview card |
| `.persona-card` | Persona display card |
| `.badge` | Status badges |
| `.tag` | Tag pills |
| `.thesis-box` | Highlighted thesis |

### Sentiment Colors

| Sentiment | Color Variable |
|-----------|---------------|
| Optimistic | `--sentiment-optimistic` |
| Cautious | `--sentiment-cautious` |
| Pessimistic | `--sentiment-pessimistic` |
| Neutral | `--sentiment-neutral` |
| Provocative | `--sentiment-provocative` |

### Stance Colors

| Stance | Color Variable |
|--------|---------------|
| Agrees | `--stance-agrees` |
| Disagrees | `--stance-disagrees` |
| Nuanced | `--stance-nuanced` |
| Tangential | `--stance-tangential` |

## Custom Themes

!!! warning "Coming Soon"
    Custom themes are not yet supported. This section describes planned functionality.

### Theme Override Directory

Place custom templates in your site's `themes/` directory:

```
my-site/
├── socialpulse.yaml
├── themes/
│   └── custom/
│       ├── assets/
│       │   └── custom.css
│       └── templates/
│           └── article.html  # Override just this template
└── content/
```

### Configuration

```yaml
# socialpulse.yaml
theme:
  name: default
  custom: ./themes/custom
```

Custom templates override default templates by name.

## Template Functions

Templates have access to these functions:

| Function | Description |
|----------|-------------|
| `truncate` | Truncate text to length |
| `formatDate` | Format date for display |
| `slugify` | Convert to URL-safe slug |
| `lower` | Lowercase string |
| `upper` | Uppercase string |

## ECharts Integration

The default theme includes ECharts for visualizations:

### Persona Donut Chart

```javascript
// Rendered from persona-data JSON
const chart = echarts.init(document.getElementById('persona-donut-chart'));
chart.setOption({
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    data: personaData
  }]
});
```

### Customizing Charts

Chart options are embedded in template `<script>` tags. To customize:

1. Override the template
2. Modify the ECharts options
3. Rebuild the site

## Responsive Behavior

The default theme is responsive:

| Breakpoint | Behavior |
|------------|----------|
| >1200px | Full sidebar, multi-column grid |
| 768-1200px | Collapsed sidebar, 2-column grid |
| <768px | Hidden sidebar, single column |

## Accessibility

The default theme includes:

- Semantic HTML elements
- ARIA labels where appropriate
- Keyboard navigation support
- Sufficient color contrast
- Focus indicators
