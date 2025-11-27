// Package locale provides localization support for Parsley
package locale

import (
	"fmt"
	"strings"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
)

// RelativeTimeUnit represents a unit of time for relative formatting
type RelativeTimeUnit string

const (
	UnitSecond RelativeTimeUnit = "second"
	UnitMinute RelativeTimeUnit = "minute"
	UnitHour   RelativeTimeUnit = "hour"
	UnitDay    RelativeTimeUnit = "day"
	UnitWeek   RelativeTimeUnit = "week"
	UnitMonth  RelativeTimeUnit = "month"
	UnitYear   RelativeTimeUnit = "year"
)

// RelativeTimeData holds the patterns for a specific unit in a locale
type RelativeTimeData struct {
	// Named relative values (e.g., -1 → "yesterday", 0 → "today", 1 → "tomorrow")
	Relative map[int]string
	// Future patterns by plural form (e.g., "in {0} day", "in {0} days")
	Future map[plural.Form]string
	// Past patterns by plural form (e.g., "{0} day ago", "{0} days ago")
	Past map[plural.Form]string
}

// LocaleRelativeTime holds all relative time data for a locale
type LocaleRelativeTime struct {
	Units map[RelativeTimeUnit]*RelativeTimeData
}

// relativeTimeLocales maps locale codes to their relative time data
var relativeTimeLocales = map[string]*LocaleRelativeTime{
	"en": englishRelativeTime(),
	"de": germanRelativeTime(),
	"fr": frenchRelativeTime(),
	"es": spanishRelativeTime(),
	"it": italianRelativeTime(),
	"pt": portugueseRelativeTime(),
	"nl": dutchRelativeTime(),
	"ru": russianRelativeTime(),
	"ja": japaneseRelativeTime(),
	"zh": chineseRelativeTime(),
	"ko": koreanRelativeTime(),
}

// FormatRelativeTime formats a duration as relative time
// value is the numeric value (negative for past, positive for future)
// unit is the time unit (second, minute, hour, day, week, month, year)
// locale is the BCP 47 locale tag (e.g., "en-US", "de-DE")
func FormatRelativeTime(value int64, unit RelativeTimeUnit, locale string) string {
	// Normalize locale
	normalizedLocale := normalizeLocale(locale)

	// Get locale data (fall back to English)
	localeData := relativeTimeLocales[normalizedLocale]
	if localeData == nil {
		localeData = relativeTimeLocales["en"]
	}

	// Get unit data
	unitData := localeData.Units[unit]
	if unitData == nil {
		// Fall back to English if unit not found
		unitData = relativeTimeLocales["en"].Units[unit]
		if unitData == nil {
			return fmt.Sprintf("%d %s", value, unit)
		}
	}

	// Check for named relative (e.g., "yesterday", "today", "tomorrow")
	if named, ok := unitData.Relative[int(value)]; ok {
		return named
	}

	// Get the appropriate plural form
	tag := language.Make(locale)
	absValue := value
	if absValue < 0 {
		absValue = -absValue
	}
	form := plural.Cardinal.MatchPlural(tag, int(absValue), 0, 0, 0, 0)

	// Get pattern based on past/future
	var pattern string
	if value < 0 {
		pattern = unitData.Past[form]
		if pattern == "" {
			pattern = unitData.Past[plural.Other]
		}
	} else {
		pattern = unitData.Future[form]
		if pattern == "" {
			pattern = unitData.Future[plural.Other]
		}
	}

	if pattern == "" {
		return fmt.Sprintf("%d %s", value, unit)
	}

	// Replace {0} with the absolute value
	return strings.Replace(pattern, "{0}", fmt.Sprintf("%d", absValue), 1)
}

// normalizeLocale extracts the language code from a locale string
func normalizeLocale(locale string) string {
	locale = strings.ToLower(locale)
	locale = strings.ReplaceAll(locale, "-", "_")
	parts := strings.Split(locale, "_")
	return parts[0]
}

// DurationToRelativeTime converts months and seconds to a relative time string
// It automatically selects the most appropriate unit
func DurationToRelativeTime(months, seconds int64, locale string) string {
	// Determine sign (negative = past, positive = future)
	negative := months < 0 || seconds < 0
	if months < 0 {
		months = -months
	}
	if seconds < 0 {
		seconds = -seconds
	}

	// Choose the best unit based on magnitude
	var value int64
	var unit RelativeTimeUnit

	if months > 0 {
		if months >= 12 {
			value = months / 12
			unit = UnitYear
		} else {
			value = months
			unit = UnitMonth
		}
	} else if seconds >= 7*24*60*60 {
		value = seconds / (7 * 24 * 60 * 60)
		unit = UnitWeek
	} else if seconds >= 24*60*60 {
		value = seconds / (24 * 60 * 60)
		unit = UnitDay
	} else if seconds >= 60*60 {
		value = seconds / (60 * 60)
		unit = UnitHour
	} else if seconds >= 60 {
		value = seconds / 60
		unit = UnitMinute
	} else {
		value = seconds
		unit = UnitSecond
	}

	// Apply sign
	if negative {
		value = -value
	}

	return FormatRelativeTime(value, unit, locale)
}

// ========================================
// Locale-specific relative time data
// Generated from CLDR data
// ========================================

func englishRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "now"},
				Future:   map[plural.Form]string{plural.One: "in {0} second", plural.Other: "in {0} seconds"},
				Past:     map[plural.Form]string{plural.One: "{0} second ago", plural.Other: "{0} seconds ago"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "in {0} minute", plural.Other: "in {0} minutes"},
				Past:     map[plural.Form]string{plural.One: "{0} minute ago", plural.Other: "{0} minutes ago"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "in {0} hour", plural.Other: "in {0} hours"},
				Past:     map[plural.Form]string{plural.One: "{0} hour ago", plural.Other: "{0} hours ago"},
			},
			UnitDay: {
				Relative: map[int]string{-1: "yesterday", 0: "today", 1: "tomorrow"},
				Future:   map[plural.Form]string{plural.One: "in {0} day", plural.Other: "in {0} days"},
				Past:     map[plural.Form]string{plural.One: "{0} day ago", plural.Other: "{0} days ago"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "last week", 0: "this week", 1: "next week"},
				Future:   map[plural.Form]string{plural.One: "in {0} week", plural.Other: "in {0} weeks"},
				Past:     map[plural.Form]string{plural.One: "{0} week ago", plural.Other: "{0} weeks ago"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "last month", 0: "this month", 1: "next month"},
				Future:   map[plural.Form]string{plural.One: "in {0} month", plural.Other: "in {0} months"},
				Past:     map[plural.Form]string{plural.One: "{0} month ago", plural.Other: "{0} months ago"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "last year", 0: "this year", 1: "next year"},
				Future:   map[plural.Form]string{plural.One: "in {0} year", plural.Other: "in {0} years"},
				Past:     map[plural.Form]string{plural.One: "{0} year ago", plural.Other: "{0} years ago"},
			},
		},
	}
}

func germanRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "jetzt"},
				Future:   map[plural.Form]string{plural.One: "in {0} Sekunde", plural.Other: "in {0} Sekunden"},
				Past:     map[plural.Form]string{plural.One: "vor {0} Sekunde", plural.Other: "vor {0} Sekunden"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "in {0} Minute", plural.Other: "in {0} Minuten"},
				Past:     map[plural.Form]string{plural.One: "vor {0} Minute", plural.Other: "vor {0} Minuten"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "in {0} Stunde", plural.Other: "in {0} Stunden"},
				Past:     map[plural.Form]string{plural.One: "vor {0} Stunde", plural.Other: "vor {0} Stunden"},
			},
			UnitDay: {
				Relative: map[int]string{-2: "vorgestern", -1: "gestern", 0: "heute", 1: "morgen", 2: "übermorgen"},
				Future:   map[plural.Form]string{plural.One: "in {0} Tag", plural.Other: "in {0} Tagen"},
				Past:     map[plural.Form]string{plural.One: "vor {0} Tag", plural.Other: "vor {0} Tagen"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "letzte Woche", 0: "diese Woche", 1: "nächste Woche"},
				Future:   map[plural.Form]string{plural.One: "in {0} Woche", plural.Other: "in {0} Wochen"},
				Past:     map[plural.Form]string{plural.One: "vor {0} Woche", plural.Other: "vor {0} Wochen"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "letzten Monat", 0: "diesen Monat", 1: "nächsten Monat"},
				Future:   map[plural.Form]string{plural.One: "in {0} Monat", plural.Other: "in {0} Monaten"},
				Past:     map[plural.Form]string{plural.One: "vor {0} Monat", plural.Other: "vor {0} Monaten"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "letztes Jahr", 0: "dieses Jahr", 1: "nächstes Jahr"},
				Future:   map[plural.Form]string{plural.One: "in {0} Jahr", plural.Other: "in {0} Jahren"},
				Past:     map[plural.Form]string{plural.One: "vor {0} Jahr", plural.Other: "vor {0} Jahren"},
			},
		},
	}
}

func frenchRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "maintenant"},
				Future:   map[plural.Form]string{plural.One: "dans {0} seconde", plural.Other: "dans {0} secondes"},
				Past:     map[plural.Form]string{plural.One: "il y a {0} seconde", plural.Other: "il y a {0} secondes"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "dans {0} minute", plural.Other: "dans {0} minutes"},
				Past:     map[plural.Form]string{plural.One: "il y a {0} minute", plural.Other: "il y a {0} minutes"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "dans {0} heure", plural.Other: "dans {0} heures"},
				Past:     map[plural.Form]string{plural.One: "il y a {0} heure", plural.Other: "il y a {0} heures"},
			},
			UnitDay: {
				Relative: map[int]string{-2: "avant-hier", -1: "hier", 0: "aujourd'hui", 1: "demain", 2: "après-demain"},
				Future:   map[plural.Form]string{plural.One: "dans {0} jour", plural.Other: "dans {0} jours"},
				Past:     map[plural.Form]string{plural.One: "il y a {0} jour", plural.Other: "il y a {0} jours"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "la semaine dernière", 0: "cette semaine", 1: "la semaine prochaine"},
				Future:   map[plural.Form]string{plural.One: "dans {0} semaine", plural.Other: "dans {0} semaines"},
				Past:     map[plural.Form]string{plural.One: "il y a {0} semaine", plural.Other: "il y a {0} semaines"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "le mois dernier", 0: "ce mois-ci", 1: "le mois prochain"},
				Future:   map[plural.Form]string{plural.One: "dans {0} mois", plural.Other: "dans {0} mois"},
				Past:     map[plural.Form]string{plural.One: "il y a {0} mois", plural.Other: "il y a {0} mois"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "l'année dernière", 0: "cette année", 1: "l'année prochaine"},
				Future:   map[plural.Form]string{plural.One: "dans {0} an", plural.Other: "dans {0} ans"},
				Past:     map[plural.Form]string{plural.One: "il y a {0} an", plural.Other: "il y a {0} ans"},
			},
		},
	}
}

func spanishRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "ahora"},
				Future:   map[plural.Form]string{plural.One: "dentro de {0} segundo", plural.Other: "dentro de {0} segundos"},
				Past:     map[plural.Form]string{plural.One: "hace {0} segundo", plural.Other: "hace {0} segundos"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "dentro de {0} minuto", plural.Other: "dentro de {0} minutos"},
				Past:     map[plural.Form]string{plural.One: "hace {0} minuto", plural.Other: "hace {0} minutos"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "dentro de {0} hora", plural.Other: "dentro de {0} horas"},
				Past:     map[plural.Form]string{plural.One: "hace {0} hora", plural.Other: "hace {0} horas"},
			},
			UnitDay: {
				Relative: map[int]string{-2: "anteayer", -1: "ayer", 0: "hoy", 1: "mañana", 2: "pasado mañana"},
				Future:   map[plural.Form]string{plural.One: "dentro de {0} día", plural.Other: "dentro de {0} días"},
				Past:     map[plural.Form]string{plural.One: "hace {0} día", plural.Other: "hace {0} días"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "la semana pasada", 0: "esta semana", 1: "la semana que viene"},
				Future:   map[plural.Form]string{plural.One: "dentro de {0} semana", plural.Other: "dentro de {0} semanas"},
				Past:     map[plural.Form]string{plural.One: "hace {0} semana", plural.Other: "hace {0} semanas"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "el mes pasado", 0: "este mes", 1: "el mes que viene"},
				Future:   map[plural.Form]string{plural.One: "dentro de {0} mes", plural.Other: "dentro de {0} meses"},
				Past:     map[plural.Form]string{plural.One: "hace {0} mes", plural.Other: "hace {0} meses"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "el año pasado", 0: "este año", 1: "el año que viene"},
				Future:   map[plural.Form]string{plural.One: "dentro de {0} año", plural.Other: "dentro de {0} años"},
				Past:     map[plural.Form]string{plural.One: "hace {0} año", plural.Other: "hace {0} años"},
			},
		},
	}
}

func italianRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "ora"},
				Future:   map[plural.Form]string{plural.One: "tra {0} secondo", plural.Other: "tra {0} secondi"},
				Past:     map[plural.Form]string{plural.One: "{0} secondo fa", plural.Other: "{0} secondi fa"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "tra {0} minuto", plural.Other: "tra {0} minuti"},
				Past:     map[plural.Form]string{plural.One: "{0} minuto fa", plural.Other: "{0} minuti fa"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "tra {0} ora", plural.Other: "tra {0} ore"},
				Past:     map[plural.Form]string{plural.One: "{0} ora fa", plural.Other: "{0} ore fa"},
			},
			UnitDay: {
				Relative: map[int]string{-2: "l'altro ieri", -1: "ieri", 0: "oggi", 1: "domani", 2: "dopodomani"},
				Future:   map[plural.Form]string{plural.One: "tra {0} giorno", plural.Other: "tra {0} giorni"},
				Past:     map[plural.Form]string{plural.One: "{0} giorno fa", plural.Other: "{0} giorni fa"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "settimana scorsa", 0: "questa settimana", 1: "settimana prossima"},
				Future:   map[plural.Form]string{plural.One: "tra {0} settimana", plural.Other: "tra {0} settimane"},
				Past:     map[plural.Form]string{plural.One: "{0} settimana fa", plural.Other: "{0} settimane fa"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "mese scorso", 0: "questo mese", 1: "mese prossimo"},
				Future:   map[plural.Form]string{plural.One: "tra {0} mese", plural.Other: "tra {0} mesi"},
				Past:     map[plural.Form]string{plural.One: "{0} mese fa", plural.Other: "{0} mesi fa"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "anno scorso", 0: "quest'anno", 1: "anno prossimo"},
				Future:   map[plural.Form]string{plural.One: "tra {0} anno", plural.Other: "tra {0} anni"},
				Past:     map[plural.Form]string{plural.One: "{0} anno fa", plural.Other: "{0} anni fa"},
			},
		},
	}
}

func portugueseRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "agora"},
				Future:   map[plural.Form]string{plural.One: "em {0} segundo", plural.Other: "em {0} segundos"},
				Past:     map[plural.Form]string{plural.One: "há {0} segundo", plural.Other: "há {0} segundos"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "em {0} minuto", plural.Other: "em {0} minutos"},
				Past:     map[plural.Form]string{plural.One: "há {0} minuto", plural.Other: "há {0} minutos"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "em {0} hora", plural.Other: "em {0} horas"},
				Past:     map[plural.Form]string{plural.One: "há {0} hora", plural.Other: "há {0} horas"},
			},
			UnitDay: {
				Relative: map[int]string{-2: "anteontem", -1: "ontem", 0: "hoje", 1: "amanhã", 2: "depois de amanhã"},
				Future:   map[plural.Form]string{plural.One: "em {0} dia", plural.Other: "em {0} dias"},
				Past:     map[plural.Form]string{plural.One: "há {0} dia", plural.Other: "há {0} dias"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "semana passada", 0: "esta semana", 1: "próxima semana"},
				Future:   map[plural.Form]string{plural.One: "em {0} semana", plural.Other: "em {0} semanas"},
				Past:     map[plural.Form]string{plural.One: "há {0} semana", plural.Other: "há {0} semanas"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "mês passado", 0: "este mês", 1: "próximo mês"},
				Future:   map[plural.Form]string{plural.One: "em {0} mês", plural.Other: "em {0} meses"},
				Past:     map[plural.Form]string{plural.One: "há {0} mês", plural.Other: "há {0} meses"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "ano passado", 0: "este ano", 1: "próximo ano"},
				Future:   map[plural.Form]string{plural.One: "em {0} ano", plural.Other: "em {0} anos"},
				Past:     map[plural.Form]string{plural.One: "há {0} ano", plural.Other: "há {0} anos"},
			},
		},
	}
}

func dutchRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "nu"},
				Future:   map[plural.Form]string{plural.One: "over {0} seconde", plural.Other: "over {0} seconden"},
				Past:     map[plural.Form]string{plural.One: "{0} seconde geleden", plural.Other: "{0} seconden geleden"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "over {0} minuut", plural.Other: "over {0} minuten"},
				Past:     map[plural.Form]string{plural.One: "{0} minuut geleden", plural.Other: "{0} minuten geleden"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "over {0} uur", plural.Other: "over {0} uur"},
				Past:     map[plural.Form]string{plural.One: "{0} uur geleden", plural.Other: "{0} uur geleden"},
			},
			UnitDay: {
				Relative: map[int]string{-2: "eergisteren", -1: "gisteren", 0: "vandaag", 1: "morgen", 2: "overmorgen"},
				Future:   map[plural.Form]string{plural.One: "over {0} dag", plural.Other: "over {0} dagen"},
				Past:     map[plural.Form]string{plural.One: "{0} dag geleden", plural.Other: "{0} dagen geleden"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "vorige week", 0: "deze week", 1: "volgende week"},
				Future:   map[plural.Form]string{plural.One: "over {0} week", plural.Other: "over {0} weken"},
				Past:     map[plural.Form]string{plural.One: "{0} week geleden", plural.Other: "{0} weken geleden"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "vorige maand", 0: "deze maand", 1: "volgende maand"},
				Future:   map[plural.Form]string{plural.One: "over {0} maand", plural.Other: "over {0} maanden"},
				Past:     map[plural.Form]string{plural.One: "{0} maand geleden", plural.Other: "{0} maanden geleden"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "vorig jaar", 0: "dit jaar", 1: "volgend jaar"},
				Future:   map[plural.Form]string{plural.One: "over {0} jaar", plural.Other: "over {0} jaar"},
				Past:     map[plural.Form]string{plural.One: "{0} jaar geleden", plural.Other: "{0} jaar geleden"},
			},
		},
	}
}

func russianRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "сейчас"},
				Future:   map[plural.Form]string{plural.One: "через {0} секунду", plural.Few: "через {0} секунды", plural.Many: "через {0} секунд", plural.Other: "через {0} секунды"},
				Past:     map[plural.Form]string{plural.One: "{0} секунду назад", plural.Few: "{0} секунды назад", plural.Many: "{0} секунд назад", plural.Other: "{0} секунды назад"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "через {0} минуту", plural.Few: "через {0} минуты", plural.Many: "через {0} минут", plural.Other: "через {0} минуты"},
				Past:     map[plural.Form]string{plural.One: "{0} минуту назад", plural.Few: "{0} минуты назад", plural.Many: "{0} минут назад", plural.Other: "{0} минуты назад"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.One: "через {0} час", plural.Few: "через {0} часа", plural.Many: "через {0} часов", plural.Other: "через {0} часа"},
				Past:     map[plural.Form]string{plural.One: "{0} час назад", plural.Few: "{0} часа назад", plural.Many: "{0} часов назад", plural.Other: "{0} часа назад"},
			},
			UnitDay: {
				Relative: map[int]string{-2: "позавчера", -1: "вчера", 0: "сегодня", 1: "завтра", 2: "послезавтра"},
				Future:   map[plural.Form]string{plural.One: "через {0} день", plural.Few: "через {0} дня", plural.Many: "через {0} дней", plural.Other: "через {0} дня"},
				Past:     map[plural.Form]string{plural.One: "{0} день назад", plural.Few: "{0} дня назад", plural.Many: "{0} дней назад", plural.Other: "{0} дня назад"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "на прошлой неделе", 0: "на этой неделе", 1: "на следующей неделе"},
				Future:   map[plural.Form]string{plural.One: "через {0} неделю", plural.Few: "через {0} недели", plural.Many: "через {0} недель", plural.Other: "через {0} недели"},
				Past:     map[plural.Form]string{plural.One: "{0} неделю назад", plural.Few: "{0} недели назад", plural.Many: "{0} недель назад", plural.Other: "{0} недели назад"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "в прошлом месяце", 0: "в этом месяце", 1: "в следующем месяце"},
				Future:   map[plural.Form]string{plural.One: "через {0} месяц", plural.Few: "через {0} месяца", plural.Many: "через {0} месяцев", plural.Other: "через {0} месяца"},
				Past:     map[plural.Form]string{plural.One: "{0} месяц назад", plural.Few: "{0} месяца назад", plural.Many: "{0} месяцев назад", plural.Other: "{0} месяца назад"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "в прошлом году", 0: "в этом году", 1: "в следующем году"},
				Future:   map[plural.Form]string{plural.One: "через {0} год", plural.Few: "через {0} года", plural.Many: "через {0} лет", plural.Other: "через {0} года"},
				Past:     map[plural.Form]string{plural.One: "{0} год назад", plural.Few: "{0} года назад", plural.Many: "{0} лет назад", plural.Other: "{0} года назад"},
			},
		},
	}
}

func japaneseRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "今"},
				Future:   map[plural.Form]string{plural.Other: "{0} 秒後"},
				Past:     map[plural.Form]string{plural.Other: "{0} 秒前"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.Other: "{0} 分後"},
				Past:     map[plural.Form]string{plural.Other: "{0} 分前"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.Other: "{0} 時間後"},
				Past:     map[plural.Form]string{plural.Other: "{0} 時間前"},
			},
			UnitDay: {
				Relative: map[int]string{-2: "一昨日", -1: "昨日", 0: "今日", 1: "明日", 2: "明後日"},
				Future:   map[plural.Form]string{plural.Other: "{0} 日後"},
				Past:     map[plural.Form]string{plural.Other: "{0} 日前"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "先週", 0: "今週", 1: "来週"},
				Future:   map[plural.Form]string{plural.Other: "{0} 週間後"},
				Past:     map[plural.Form]string{plural.Other: "{0} 週間前"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "先月", 0: "今月", 1: "来月"},
				Future:   map[plural.Form]string{plural.Other: "{0} か月後"},
				Past:     map[plural.Form]string{plural.Other: "{0} か月前"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "昨年", 0: "今年", 1: "来年"},
				Future:   map[plural.Form]string{plural.Other: "{0} 年後"},
				Past:     map[plural.Form]string{plural.Other: "{0} 年前"},
			},
		},
	}
}

func chineseRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "现在"},
				Future:   map[plural.Form]string{plural.Other: "{0}秒钟后"},
				Past:     map[plural.Form]string{plural.Other: "{0}秒钟前"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.Other: "{0}分钟后"},
				Past:     map[plural.Form]string{plural.Other: "{0}分钟前"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.Other: "{0}小时后"},
				Past:     map[plural.Form]string{plural.Other: "{0}小时前"},
			},
			UnitDay: {
				Relative: map[int]string{-2: "前天", -1: "昨天", 0: "今天", 1: "明天", 2: "后天"},
				Future:   map[plural.Form]string{plural.Other: "{0}天后"},
				Past:     map[plural.Form]string{plural.Other: "{0}天前"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "上周", 0: "本周", 1: "下周"},
				Future:   map[plural.Form]string{plural.Other: "{0}周后"},
				Past:     map[plural.Form]string{plural.Other: "{0}周前"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "上个月", 0: "本月", 1: "下个月"},
				Future:   map[plural.Form]string{plural.Other: "{0}个月后"},
				Past:     map[plural.Form]string{plural.Other: "{0}个月前"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "去年", 0: "今年", 1: "明年"},
				Future:   map[plural.Form]string{plural.Other: "{0}年后"},
				Past:     map[plural.Form]string{plural.Other: "{0}年前"},
			},
		},
	}
}

func koreanRelativeTime() *LocaleRelativeTime {
	return &LocaleRelativeTime{
		Units: map[RelativeTimeUnit]*RelativeTimeData{
			UnitSecond: {
				Relative: map[int]string{0: "지금"},
				Future:   map[plural.Form]string{plural.Other: "{0}초 후"},
				Past:     map[plural.Form]string{plural.Other: "{0}초 전"},
			},
			UnitMinute: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.Other: "{0}분 후"},
				Past:     map[plural.Form]string{plural.Other: "{0}분 전"},
			},
			UnitHour: {
				Relative: map[int]string{},
				Future:   map[plural.Form]string{plural.Other: "{0}시간 후"},
				Past:     map[plural.Form]string{plural.Other: "{0}시간 전"},
			},
			UnitDay: {
				Relative: map[int]string{-2: "그저께", -1: "어제", 0: "오늘", 1: "내일", 2: "모레"},
				Future:   map[plural.Form]string{plural.Other: "{0}일 후"},
				Past:     map[plural.Form]string{plural.Other: "{0}일 전"},
			},
			UnitWeek: {
				Relative: map[int]string{-1: "지난주", 0: "이번 주", 1: "다음 주"},
				Future:   map[plural.Form]string{plural.Other: "{0}주 후"},
				Past:     map[plural.Form]string{plural.Other: "{0}주 전"},
			},
			UnitMonth: {
				Relative: map[int]string{-1: "지난달", 0: "이번 달", 1: "다음 달"},
				Future:   map[plural.Form]string{plural.Other: "{0}개월 후"},
				Past:     map[plural.Form]string{plural.Other: "{0}개월 전"},
			},
			UnitYear: {
				Relative: map[int]string{-1: "작년", 0: "올해", 1: "내년"},
				Future:   map[plural.Form]string{plural.Other: "{0}년 후"},
				Past:     map[plural.Form]string{plural.Other: "{0}년 전"},
			},
		},
	}
}
