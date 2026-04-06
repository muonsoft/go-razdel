package sentenize

// Punctuation and delimiter sets mirror third_party/razdel/razdel/segmenters/punct.py
// and sentenize.DELIMITERS (ENDINGS + ';' + GENERIC_QUOTES + CLOSE_QUOTES + CLOSE_BRACKETS).

const (
	// Endings is upstream ENDINGS.
	Endings = ".?!\u2026"

	// OpenQuotes is upstream OPEN_QUOTES: « “ ‘
	OpenQuotes = "\u00ab\u201c\u2018"
	// CloseQuotes is upstream CLOSE_QUOTES: » ” ’
	CloseQuotes = "\u00bb\u201d\u2019"
	// GenericQuotes is upstream GENERIC_QUOTES: ASCII ", U+201E „, ASCII apostrophe.
	GenericQuotes = "\"\u201e'"
	// Quotes is upstream QUOTES (OPEN + CLOSE + GENERIC).
	Quotes = OpenQuotes + CloseQuotes + GenericQuotes

	// OpenBrackets is upstream OPEN_BRACKETS.
	OpenBrackets = "([{"
	// CloseBrackets is upstream CLOSE_BRACKETS.
	CloseBrackets = ")]}"

	// Delimiters is upstream sentenize.DELIMITERS (splitter / delimiter_right set).
	Delimiters = Endings + ";" + GenericQuotes + CloseQuotes + CloseBrackets

	// Dashes is upstream punct.DASHES (sentenize.dash_right).
	Dashes = "\u2011\u2013\u2014\u2212-"
)
