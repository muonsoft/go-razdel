package razdel_test

import (
	"strings"
	"testing"

	"github.com/muonsoft/go-razdel"
)

// BenchCorpus holds representative inputs for docs/contracts.md performance notes (T018).
var benchShortASCII = "The quick brown fox jumps over the lazy dog."

var benchCyrillicParagraph = strings.TrimSpace(`
Федеральный министр транспорта объявил, что правительство выделяет средства на исследования
в области технологии водородного двигателя. Это заявление было сделано на открытии ярмарки.
Среди тем — безопасность, экономика и экология; т. е. вопросы комплексные.
`)

var benchPunctHeavy = strings.Repeat(`("a", 'b' [c] {d}; e: f! g? h—i/j\\k) `, 80)

var benchUnicodeHeavy = strings.Repeat("привет мир 🙂 — «цитата» ", 40)

func BenchmarkTokenize(b *testing.B) {
	benches := []struct {
		name string
		text string
	}{
		{"short_ASCII", benchShortASCII},
		{"long_Cyrillic", benchCyrillicParagraph},
		{"punct_heavy", benchPunctHeavy},
		{"unicode_mixed", benchUnicodeHeavy},
	}
	for _, bc := range benches {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(bc.text)))
			for b.Loop() {
				_ = razdel.Tokenize(bc.text)
			}
		})
	}
}

func BenchmarkSentenize(b *testing.B) {
	benches := []struct {
		name string
		text string
	}{
		{"short_ASCII", "First. Second. Third?"},
		{"long_Cyrillic", benchCyrillicParagraph},
		{"punct_heavy", benchPunctHeavy},
		{"unicode_mixed", benchUnicodeHeavy},
	}
	for _, bc := range benches {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(bc.text)))
			for b.Loop() {
				_ = razdel.Sentenize(bc.text)
			}
		})
	}
}
