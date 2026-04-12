package tokenize

import "unicode/utf8"

// isEmojiLikeRune approximates Unicode Extended_Pictographic plus common emoji
// joiners (ZWJ, skin tones, VS, regional indicators) so we can split emoji-shaped
// glyphs from letter atoms without pulling in every OTHER symbol (IGO-47).
func isEmojiLikeRune(r rune) bool {
	switch {
	case r == 0xFE0F || r == 0xFE0E:
		return true
	case r == 0x200D || r == 0x200C:
		return true
	case r == 0x20E3:
		return true
	case r >= 0x1F1E6 && r <= 0x1F1FF:
		return true
	case r >= 0x1F3FB && r <= 0x1F3FF:
		return true
	case r >= 0x1F000 && r <= 0x1FFFF:
		return true
	case r >= 0x2600 && r <= 0x27BF:
		return true
	case r >= 0x2300 && r <= 0x23FF:
		return true
	case r >= 0x2B00 && r <= 0x2BFF:
		return true
	case r == 0x24C2:
		return true
	}
	return false
}

func firstRune(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

// otherShouldSplitEmojiFromWord reports whether OTHER join must not run so that
// emoji-like OTHER atoms stay separate from adjacent RU/LAT letters.
func otherShouldSplitEmojiFromWord(left, right Atom) bool {
	if left.Type == OTHER && (right.Type == RU || right.Type == LAT) {
		return isEmojiLikeRune(firstRune(left.Text))
	}
	if right.Type == OTHER && (left.Type == RU || left.Type == LAT) {
		return isEmojiLikeRune(firstRune(right.Text))
	}
	return false
}
