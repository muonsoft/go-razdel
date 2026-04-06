package razdel_test

import (
	"fmt"
	"unicode/utf8"

	"github.com/muonsoft/go-razdel"
)

func ExampleTokenize() {
	text := "Привет, мир!"
	tokens := razdel.Tokenize(text)

	for _, tok := range tokens {
		fmt.Printf("%q [%d:%d]\n", tok.Text, tok.Start, tok.End)
	}

	// Output:
	// "Привет" [0:12]
	// "," [12:13]
	// "мир" [14:20]
	// "!" [20:21]
}

func ExampleSentenize() {
	text := "Привет, мир! Это тест."
	sentences := razdel.Sentenize(text)

	for _, sent := range sentences {
		fmt.Printf("%q [%d:%d]\n", sent.Text, sent.Start, sent.End)
	}

	// Output:
	// "Привет, мир!" [0:21]
	// "Это тест." [22:38]
}

func ExampleTokenize_utf8ByteOffsets() {
	text := "a ж"
	fmt.Printf("bytes=%d runes=%d\n", len(text), utf8.RuneCountInString(text))

	for _, tok := range razdel.Tokenize(text) {
		fmt.Printf("%q [%d:%d]\n", tok.Text, tok.Start, tok.End)
	}

	// Output:
	// bytes=4 runes=3
	// "a" [0:1]
	// "ж" [2:4]
}

func ExampleSentenize_trimmedSpans() {
	text := "  Привет. Пока.  "

	for _, sent := range razdel.Sentenize(text) {
		fmt.Printf("%q [%d:%d]\n", sent.Text, sent.Start, sent.End)
	}

	// Output:
	// "Привет." [2:15]
	// "Пока." [16:25]
}
