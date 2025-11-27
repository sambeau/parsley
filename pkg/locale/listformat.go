// Package locale provides localization support for Parsley
// This file implements CLDR-based list formatting
package locale

import (
	"strings"
)

// ListStyle represents the style of list formatting
type ListStyle string

const (
	ListStyleAnd  ListStyle = "and"  // "A, B, and C"
	ListStyleOr   ListStyle = "or"   // "A, B, or C"
	ListStyleUnit ListStyle = "unit" // "A, B, C" (for units like "5 feet, 6 inches")
)

// ListPatterns holds the formatting patterns for lists
type ListPatterns struct {
	// Start pattern for first two items in a list of 3+
	Start string
	// Middle pattern for middle items in a list of 4+
	Middle string
	// End pattern for last two items in a list of 3+
	End string
	// Two pattern for exactly two items
	Two string
}

// LocaleListPatterns holds list patterns for all styles in a locale
type LocaleListPatterns struct {
	And  *ListPatterns
	Or   *ListPatterns
	Unit *ListPatterns
}

// listPatternLocales maps locale codes to their list patterns
var listPatternLocales = map[string]*LocaleListPatterns{
	"en":    englishListPatterns(),
	"en-US": englishUSListPatterns(),
	"en-GB": englishGBListPatterns(),
	"de":    germanListPatterns(),
	"fr":    frenchListPatterns(),
	"es":    spanishListPatterns(),
	"it":    italianListPatterns(),
	"pt":    portugueseListPatterns(),
	"nl":    dutchListPatterns(),
	"ru":    russianListPatterns(),
	"ja":    japaneseListPatterns(),
	"zh":    chineseListPatterns(),
	"ko":    koreanListPatterns(),
}

// FormatList formats a list of items according to locale and style
// items is the list of strings to format
// style is "and", "or", or "unit" (defaults to "and")
// locale is the BCP 47 locale tag (e.g., "en-US", "de-DE")
func FormatList(items []string, style ListStyle, locale string) string {
	// Handle empty and single-item lists
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}

	// Get locale patterns
	patterns := getListPatterns(locale, style)

	// Handle two-item list
	if len(items) == 2 {
		return applyPattern(patterns.Two, items[0], items[1])
	}

	// Build list from end to start for 3+ items
	// Start with end pattern for last two items
	result := applyPattern(patterns.End, items[len(items)-2], items[len(items)-1])

	// Apply middle pattern for middle items (if 4+ items)
	for i := len(items) - 3; i > 0; i-- {
		result = applyPattern(patterns.Middle, items[i], result)
	}

	// Apply start pattern for first item
	result = applyPattern(patterns.Start, items[0], result)

	return result
}

// getListPatterns returns the patterns for a locale and style
func getListPatterns(locale string, style ListStyle) *ListPatterns {
	// Try exact locale match first
	if patterns := listPatternLocales[locale]; patterns != nil {
		return getStylePatterns(patterns, style)
	}

	// Try base locale
	baseLocale := normalizeLocale(locale)
	if patterns := listPatternLocales[baseLocale]; patterns != nil {
		return getStylePatterns(patterns, style)
	}

	// Fall back to English
	return getStylePatterns(listPatternLocales["en"], style)
}

// getStylePatterns returns the patterns for a specific style
func getStylePatterns(locale *LocaleListPatterns, style ListStyle) *ListPatterns {
	switch style {
	case ListStyleOr:
		return locale.Or
	case ListStyleUnit:
		return locale.Unit
	default:
		return locale.And
	}
}

// applyPattern applies a list pattern with two items
// Pattern uses {0} for first item and {1} for second item
func applyPattern(pattern, first, second string) string {
	result := strings.Replace(pattern, "{0}", first, 1)
	result = strings.Replace(result, "{1}", second, 1)
	return result
}

// ============================================================================
// English (default) - Uses Oxford comma
// ============================================================================

func englishListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, and {1}",
			Two:    "{0} and {1}",
		},
		Or: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, or {1}",
			Two:    "{0} or {1}",
		},
		Unit: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, {1}",
			Two:    "{0}, {1}",
		},
	}
}

// English (US) - Uses Oxford comma (same as default English)
func englishUSListPatterns() *LocaleListPatterns {
	return englishListPatterns()
}

// English (GB) - No Oxford comma
func englishGBListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} and {1}",
			Two:    "{0} and {1}",
		},
		Or: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} or {1}",
			Two:    "{0} or {1}",
		},
		Unit: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, {1}",
			Two:    "{0}, {1}",
		},
	}
}

// ============================================================================
// German - No Oxford comma, uses "und"/"oder"
// ============================================================================

func germanListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} und {1}",
			Two:    "{0} und {1}",
		},
		Or: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} oder {1}",
			Two:    "{0} oder {1}",
		},
		Unit: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, {1}",
			Two:    "{0}, {1}",
		},
	}
}

// ============================================================================
// French - No Oxford comma, uses "et"/"ou"
// ============================================================================

func frenchListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} et {1}",
			Two:    "{0} et {1}",
		},
		Or: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} ou {1}",
			Two:    "{0} ou {1}",
		},
		Unit: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, {1}",
			Two:    "{0}, {1}",
		},
	}
}

// ============================================================================
// Spanish - No Oxford comma, uses "y"/"o"
// ============================================================================

func spanishListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} y {1}",
			Two:    "{0} y {1}",
		},
		Or: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} o {1}",
			Two:    "{0} o {1}",
		},
		Unit: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, {1}",
			Two:    "{0}, {1}",
		},
	}
}

// ============================================================================
// Italian - No Oxford comma, uses "e"/"o"
// ============================================================================

func italianListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} e {1}",
			Two:    "{0} e {1}",
		},
		Or: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} o {1}",
			Two:    "{0} o {1}",
		},
		Unit: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, {1}",
			Two:    "{0}, {1}",
		},
	}
}

// ============================================================================
// Portuguese - No Oxford comma, uses "e"/"ou"
// ============================================================================

func portugueseListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} e {1}",
			Two:    "{0} e {1}",
		},
		Or: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} ou {1}",
			Two:    "{0} ou {1}",
		},
		Unit: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, {1}",
			Two:    "{0}, {1}",
		},
	}
}

// ============================================================================
// Dutch - No Oxford comma, uses "en"/"of"
// ============================================================================

func dutchListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} en {1}",
			Two:    "{0} en {1}",
		},
		Or: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} of {1}",
			Two:    "{0} of {1}",
		},
		Unit: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, {1}",
			Two:    "{0}, {1}",
		},
	}
}

// ============================================================================
// Russian - No Oxford comma, uses "и"/"или"
// ============================================================================

func russianListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} и {1}",
			Two:    "{0} и {1}",
		},
		Or: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} или {1}",
			Two:    "{0} или {1}",
		},
		Unit: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0}, {1}",
			Two:    "{0}, {1}",
		},
	}
}

// ============================================================================
// Japanese - Uses "、" as separator, "と"/"か" for conjunction
// ============================================================================

func japaneseListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}、{1}",
			Middle: "{0}、{1}",
			End:    "{0}、{1}",
			Two:    "{0}と{1}",
		},
		Or: &ListPatterns{
			Start:  "{0}、{1}",
			Middle: "{0}、{1}",
			End:    "{0}、または{1}",
			Two:    "{0}または{1}",
		},
		Unit: &ListPatterns{
			Start:  "{0} {1}",
			Middle: "{0} {1}",
			End:    "{0} {1}",
			Two:    "{0} {1}",
		},
	}
}

// ============================================================================
// Chinese - Uses "、" as separator, "和"/"或" for conjunction
// ============================================================================

func chineseListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}、{1}",
			Middle: "{0}、{1}",
			End:    "{0}和{1}",
			Two:    "{0}和{1}",
		},
		Or: &ListPatterns{
			Start:  "{0}、{1}",
			Middle: "{0}、{1}",
			End:    "{0}或{1}",
			Two:    "{0}或{1}",
		},
		Unit: &ListPatterns{
			Start:  "{0} {1}",
			Middle: "{0} {1}",
			End:    "{0} {1}",
			Two:    "{0} {1}",
		},
	}
}

// ============================================================================
// Korean - Uses ", " as separator, uses particles
// ============================================================================

func koreanListPatterns() *LocaleListPatterns {
	return &LocaleListPatterns{
		And: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} 및 {1}",
			Two:    "{0} 및 {1}",
		},
		Or: &ListPatterns{
			Start:  "{0}, {1}",
			Middle: "{0}, {1}",
			End:    "{0} 또는 {1}",
			Two:    "{0} 또는 {1}",
		},
		Unit: &ListPatterns{
			Start:  "{0} {1}",
			Middle: "{0} {1}",
			End:    "{0} {1}",
			Two:    "{0} {1}",
		},
	}
}
