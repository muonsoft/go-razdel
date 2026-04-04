package tokenize

// Punctuation character sets mirror third_party/razdel/razdel/segmenters/punct.py.

const (
	endings   = ".?!…"
	dashes    = "\u2011\u2013\u2014\u2212-" // ‑ – — − -
	openQuote = "«“‘"
	closeQuote = "»”’"
	genericQuote = "\"„'"
	openBracket  = "([{"
	closeBracket = ")]}"
)

// Quotes is OPEN + CLOSE + GENERIC (upstream QUOTES).
const Quotes = openQuote + closeQuote + genericQuote

// Brackets is OPEN + CLOSE (upstream BRACKETS).
const Brackets = openBracket + closeBracket

// Dashes is upstream DASHES.
const Dashes = dashes

// Endings is upstream ENDINGS.
const Endings = endings

// PunctSet is the character class used in upstream ATOM's PUNCT alternative
// (PUNCTS in tokenize.py: ascii punct + № … + dashes + quotes + brackets).
const PunctSet = "\\/.!#$%&*+,.:;<=>?@^_`|~№…" + Dashes + Quotes + Brackets
