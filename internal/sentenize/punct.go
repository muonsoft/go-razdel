package sentenize

// Punctuation and delimiter sets mirror third_party/razdel/razdel/segmenters/punct.py
// and sentenize.DELIMITERS (ENDINGS + ';' + GENERIC_QUOTES + CLOSE_QUOTES + CLOSE_BRACKETS).

const (
	// genericQuotes is upstream GENERIC_QUOTES: ASCII ", U+201E „, ASCII apostrophe.
	genericQuotes = "\"\u201e'"
	// closeQuotes is upstream CLOSE_QUOTES: » ” ’
	closeQuotes = "\u00bb\u201d\u2019"
	// delimiters is the character set used by upstream delimiter_right (substring membership).
	delimiters = ".?!…;" + genericQuotes + closeQuotes + ")]}"
)
