# Localization Plan for Parsley

## Implementation Status

| Phase | Feature | Status | Notes |
|-------|---------|--------|-------|
| 1 | `formatNumber()` | ✅ Complete | Using golang.org/x/text |
| 1 | `formatCurrency()` | ✅ Complete | Using golang.org/x/text/currency |
| 1 | `formatPercent()` | ✅ Complete | Using golang.org/x/text/number |
| 2 | `formatDate()` | ✅ Complete | Using github.com/goodsign/monday |
| 3 | `format()` (durations) | ✅ Complete | Custom CLDR implementation in pkg/locale |
| 3 | Negative duration literals | ✅ Complete | `@-1d`, `@-2w`, etc. |
| 4 | `format()` (arrays/lists) | ✅ Complete | Custom CLDR implementation in pkg/locale |

## Executive Summary

This document outlines a plan for adding localization (l10n) support to Parsley, covering date/time formatting, number formatting, duration/relative time formatting, and list formatting. The approach is inspired by JavaScript's `Intl` API while leveraging Go's available libraries.

## Current State

Parsley currently has:
- `now()` and `time()` functions for datetime objects
- Duration literals with `@` syntax: `@1d`, `@2h30m`, `@-1d` (negative durations)
- `.format` computed property returning English-only "Month DD, YYYY" format
- `.date` and `.time` computed properties returning ISO formats
- Locale-aware formatting functions: `formatNumber()`, `formatCurrency()`, `formatPercent()`, `formatDate()`, `format()` (for durations and arrays)

## Goals

1. **Locale-aware formatting** for dates, numbers, currencies, durations, and lists
2. **Consistent API** inspired by JavaScript's `Intl` standard
3. **Minimal binary size impact** through selective CLDR data inclusion
4. **Graceful fallbacks** when locale data is unavailable

---

## 1. Available Go Libraries

### golang.org/x/text (Official)

The official Go text processing library provides:

| Package | Purpose | Status |
|---------|---------|--------|
| `language` | BCP 47 language tag parsing & matching | ✅ Excellent |
| `message` | Locale-sensitive string formatting | ✅ Good |
| `number` | Number formatting with grouping/decimals | ✅ Good |
| `currency` | Currency symbols and formatting | ✅ Good |
| `collate` | Language-specific string sorting | ✅ Good |

**Gaps**: No date/time formatting, no relative time, no list formatting.

### github.com/goodsign/monday

Localized date/time formatting for 40+ locales:
- Wraps `time.Format` with locale translation
- Handles month/day name translation
- Supports grammatical cases (e.g., Russian genitive)

```go
monday.Format(t, "2 January 2006", monday.LocaleFrFR)  // "12 avril 2024"
```

### github.com/nicksnyder/go-i18n

Message translation and pluralization:
- CLDR-based pluralization rules
- YAML/JSON message file support
- Template integration

---

## 2. Proposed Parsley API

### Option A: Object-Oriented API (JavaScript Intl Style)

```parsley
// Number formatting
let nf = NumberFormat("de-DE")
nf.format(1234567.89)  // "1.234.567,89"

// Currency formatting
let cf = NumberFormat("en-US", {style: "currency", currency: "USD"})
cf.format(99.99)  // "$99.99"

// Date/time formatting
let dtf = DateTimeFormat("fr-FR", {dateStyle: "long"})
dtf.format(now())  // "27 novembre 2025"

// Relative time
let rtf = RelativeTimeFormat("en", {numeric: "auto"})
rtf.format(-1, "day")  // "yesterday"

// List formatting
let lf = ListFormat("en", {type: "conjunction"})
lf.format(["Alice", "Bob", "Charlie"])  // "Alice, Bob, and Charlie"
```

### Option B: Function-Based API (Simpler)

```parsley
// Number formatting
formatNumber(1234567.89, "de-DE")  // "1.234.567,89"

// Currency formatting
formatCurrency(99.99, "USD", "en-US")  // "$99.99"

// Date/time formatting
formatDate(now(), "long", "fr-FR")  // "27 novembre 2025"
formatDateTime(now(), "full", "short", "en-US")  // "Thursday, November 27, 2025 at 3:45 PM"

// Relative time
formatRelativeTime(-1, "day", "en")  // "yesterday"
formatRelativeTime(3, "hour", "es")  // "dentro de 3 horas"

// List formatting
formatList(["Alice", "Bob", "Charlie"], "and", "en")  // "Alice, Bob, and Charlie"
formatList(["red", "blue"], "or", "en")  // "red or blue"
```

### Option C: Method Extensions on Existing Types (Recommended)

Extend datetime and number handling with locale-aware methods:

```parsley
// DateTime computed properties with locale
let d = time({year: 2024, month: 12, day: 25})
d.format                     // "December 25, 2024" (current behavior, English default)
d.format("long", "fr-FR")    // "25 décembre 2024"
d.format("short", "de-DE")   // "25.12.24"
d.format("full", "ja-JP")    // "2024年12月25日水曜日"

// Relative time
d.relative("en")             // "in 28 days"
d.relative("es")             // "dentro de 28 días"

// Number formatting via format() function
format(1234567.89, "number", "de-DE")      // "1.234.567,89"
format(99.99, "currency:USD", "en-US")     // "$99.99"
format(0.25, "percent", "en-US")           // "25%"

// List formatting
format(["a", "b", "c"], "list:and", "en")  // "a, b, and c"
format(["a", "b", "c"], "list:or", "de")   // "a, b oder c"
```

---

## 3. Locale Identifier Format

Use BCP 47 language tags (same as JavaScript):

```
language[-script][-region]
```

Examples:
- `"en"` - English
- `"en-US"` - English (United States)
- `"en-GB"` - English (United Kingdom)  
- `"zh-Hans"` - Chinese (Simplified script)
- `"zh-Hant-TW"` - Chinese (Traditional script, Taiwan)
- `"pt-BR"` - Portuguese (Brazil)

### Locale Fallback Chain

```parsley
// Request: "en-AU" (Australian English)
// Fallback: "en-AU" → "en" → default

// Request: "sr-Latn-RS" (Serbian in Latin script, Serbia)
// Fallback: "sr-Latn-RS" → "sr-Latn" → "sr" → default
```

---

## 4. Implementation Plan

### Phase 1: Number & Currency Formatting

**Effort**: Low (use golang.org/x/text)

```go
// In evaluator.go
import (
    "golang.org/x/text/language"
    "golang.org/x/text/message"
    "golang.org/x/text/currency"
)

func formatNumber(value float64, locale string) string {
    tag := language.MustParse(locale)
    p := message.NewPrinter(tag)
    return p.Sprintf("%v", value)
}

func formatCurrency(value float64, curr, locale string) string {
    tag := language.MustParse(locale)
    cur := currency.MustParseISO(curr)
    p := message.NewPrinter(tag)
    return p.Sprintf("%v", currency.Symbol(cur.Amount(value)))
}
```

**Parsley API**:
```parsley
formatNumber(1234.56, "de-DE")        // "1.234,56"
formatCurrency(1234.56, "EUR", "de")  // "1.234,56 €"
formatPercent(0.1234, "en-US")        // "12.34%"
```

### Phase 2: Date/Time Formatting

**Effort**: Low (use github.com/goodsign/monday)

```go
import "github.com/goodsign/monday"

func formatDate(t time.Time, style, locale string) string {
    loc := parseMonthLocale(locale)  // Map "en-US" → monday.LocaleEnUS
    format := getFormatForStyle(style, loc)  // "long" → "January 2, 2006"
    return monday.Format(t, format, loc)
}
```

**Supported Styles**:
- `"short"` - Numeric: "12/25/24" or "25.12.24"
- `"medium"` - Abbreviated: "Dec 25, 2024" or "25 déc. 2024"
- `"long"` - Full month: "December 25, 2024" or "25 décembre 2024"
- `"full"` - With weekday: "Wednesday, December 25, 2024"

**Parsley API**:
```parsley
let d = time({year: 2024, month: 12, day: 25})
d.toLocaleString("en-US")                    // "12/25/2024"
d.toLocaleString("de-DE", {dateStyle: "long"})  // "25. Dezember 2024"

// Or simpler function API
formatDate(d, "long", "fr-FR")               // "25 décembre 2024"
```

### Phase 3: Relative Time Formatting

**Effort**: Medium (custom implementation required)

No good Go library exists. Must implement using CLDR data.

**CLDR Relative Time Data Structure**:
```json
{
  "en": {
    "day": {
      "-1": "yesterday",
      "0": "today", 
      "1": "tomorrow",
      "past": {
        "one": "{0} day ago",
        "other": "{0} days ago"
      },
      "future": {
        "one": "in {0} day",
        "other": "in {0} days"
      }
    }
  },
  "es": {
    "day": {
      "-2": "anteayer",
      "-1": "ayer",
      "0": "hoy",
      "1": "mañana",
      "2": "pasado mañana"
    }
  }
}
```

**Implementation Approach**:
1. Extract relative time data from CLDR JSON
2. Generate Go code at build time
3. Implement pluralization using `golang.org/x/text/feature/plural`

**Parsley API**:
```parsley
let d = time({year: 2024, month: 12, day: 25})
let now = now()

// Relative to now
d.relative("en")                // "in 28 days"
d.relative("es")                // "dentro de 28 días"

// Or standalone function
formatRelativeTime(-1, "day", "en")     // "yesterday"
formatRelativeTime(-1, "day", "de")     // "gestern"
formatRelativeTime(2, "hour", "fr")     // "dans 2 heures"
```

### Phase 4: List Formatting

**Effort**: Medium (custom implementation required)

**CLDR List Pattern Data**:
```json
{
  "en": {
    "conjunction": {
      "start": "{0}, {1}",
      "middle": "{0}, {1}",
      "end": "{0}, and {1}",
      "pair": "{0} and {1}"
    },
    "disjunction": {
      "end": "{0}, or {1}",
      "pair": "{0} or {1}"
    }
  },
  "de": {
    "conjunction": {
      "end": "{0} und {1}",
      "pair": "{0} und {1}"
    },
    "disjunction": {
      "end": "{0} oder {1}",
      "pair": "{0} oder {1}"
    }
  }
}
```

**Implementation**:
```go
func formatList(items []string, listType, locale string) string {
    if len(items) == 0 { return "" }
    if len(items) == 1 { return items[0] }
    if len(items) == 2 {
        return applyPattern(patterns[locale][listType]["pair"], items[0], items[1])
    }
    // Build from start → middle → end patterns
}
```

**Parsley API**:
```parsley
formatList(["Alice", "Bob", "Charlie"], "and", "en")    // "Alice, Bob, and Charlie"
formatList(["Alice", "Bob", "Charlie"], "and", "en-GB") // "Alice, Bob and Charlie"
formatList(["red", "blue", "green"], "or", "de")        // "rot, blau oder grün"
```

---

## 5. Common Pitfalls & Mitigations

### Pitfall 1: Time Zone vs Locale Confusion

**Problem**: Locale doesn't determine timezone. "de-DE" means German formatting, not Berlin time.

**Solution**: Keep timezone and locale separate:
```parsley
let t = now()                           // Current time in system timezone
let berlin = t.in("Europe/Berlin")      // Convert to Berlin time
formatDateTime(berlin, "full", "de-DE") // Format in German
```

### Pitfall 2: Date Format Ambiguity

**Problem**: "01/02/03" means January 2nd in US, February 1st in EU.

**Solution**: Always use unambiguous ISO format for storage/exchange, locale format only for display:
```parsley
let d = time("2024-01-02")              // Unambiguous ISO format
d.format("short", "en-US")              // "1/2/2024" for US users
d.format("short", "en-GB")              // "02/01/2024" for UK users
```

### Pitfall 3: Pluralization Complexity

**Problem**: Languages have different plural rules:
- English: 2 forms (1 apple, 2 apples)
- Russian: 4 forms (1 яблоко, 2 яблока, 5 яблок, 21 яблоко)
- Arabic: 6 forms

**Solution**: Use CLDR plural rules, never hardcode:
```parsley
// DON'T
if count == 1 { "1 item" } else { count + " items" }

// DO
formatPlural(count, "item", "en")  // Uses CLDR rules
```

### Pitfall 4: Currency Symbol Placement

**Problem**: Currency symbols appear in different positions:
- US: $100.00 (before, no space)
- Germany: 100,00 € (after, with space)
- Switzerland: CHF 100.00 (before, with space)

**Solution**: Use `currency` formatting, never manual concatenation:
```parsley
// DON'T
"$" + amount

// DO
formatCurrency(amount, "USD", "en-US")  // Handles placement automatically
```

### Pitfall 5: RTL Languages

**Problem**: Hebrew, Arabic require right-to-left text direction.

**Solution**: 
1. Return Unicode directional markers when needed
2. Document that HTML output should use `dir="auto"` or explicit `dir="rtl"`

---

## 6. Binary Size Considerations

### CLDR Data Size

| Data Type | Approximate Size |
|-----------|------------------|
| Number formats (all locales) | ~500KB |
| Date formats (all locales) | ~2MB |
| Relative time (all locales) | ~1MB |
| List patterns (all locales) | ~200KB |
| Pluralization rules | ~100KB |

### Strategies

1. **Embed Common Locales Only**: Include ~20 most common locales by default
   - en, en-GB, de, fr, es, it, pt, pt-BR, ja, zh-Hans, zh-Hant, ko, ru, ar, hi, nl, pl, tr, vi, th

2. **Lazy Loading** (Future): Load additional locale data from external files
   ```parsley
   loadLocale("sw")  // Load Swahili data at runtime
   ```

3. **Tree Shaking**: Only include formatters actually used in code (complex to implement)

---

## 7. Recommended Implementation Order

| Phase | Feature | Effort | Dependencies |
|-------|---------|--------|--------------|
| 1 | Number formatting | Low | golang.org/x/text |
| 2 | Currency formatting | Low | golang.org/x/text |
| 3 | Percentage formatting | Low | golang.org/x/text |
| 4 | Date/time formatting | Low | github.com/goodsign/monday |
| 5 | Relative time formatting | Medium | Custom + CLDR data |
| 6 | List formatting | Medium | Custom + CLDR data |
| 7 | Duration formatting | Medium | Custom + pluralization |

---

## 8. Example Usage After Implementation

```parsley
// Set default locale (optional, defaults to "en")
setLocale("fr-FR")

// Numbers
log(formatNumber(1234567.89))                    // "1 234 567,89"
log(formatCurrency(99.99, "EUR"))                // "99,99 €"
log(formatPercent(0.1234))                       // "12,34 %"

// Dates
let christmas = time({year: 2024, month: 12, day: 25})
log(christmas.format("long"))                    // "25 décembre 2024"
log(christmas.format("full"))                    // "mercredi 25 décembre 2024"
log(christmas.relative())                        // "dans 28 jours"

// Lists
let items = ["pommes", "oranges", "bananes"]
log(formatList(items, "and"))                    // "pommes, oranges et bananes"

// Override locale for specific call
log(christmas.format("long", "de-DE"))           // "25. Dezember 2024"
log(formatNumber(1234567.89, "en-IN"))           // "12,34,567.89" (Indian grouping)
```

---

## 9. Open Questions

1. **Global locale state?** Should there be a `setLocale()` that affects all formatting, or always require explicit locale parameter?

2. **Timezone handling?** Should locale imply a default timezone, or keep them completely separate?

3. **Custom format patterns?** Should users be able to specify custom patterns like `"dd/MM/yyyy"`, or only use predefined styles?

4. **Which locales to include by default?** All CLDR locales (~200) or a subset (~20-40)?

5. **Error handling?** What happens with invalid locale codes? Silent fallback to English, or error?

---

## 10. References

- [CLDR - Unicode Common Locale Data Repository](https://cldr.unicode.org/)
- [JavaScript Intl API - MDN](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl)
- [golang.org/x/text](https://pkg.go.dev/golang.org/x/text)
- [github.com/goodsign/monday](https://github.com/goodsign/monday)
- [github.com/nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n)
- [BCP 47 Language Tags](https://www.rfc-editor.org/info/bcp47)
