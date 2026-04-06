package razdel_test

import (
	"testing"
	"unicode/utf8"

	"github.com/muonsoft/go-razdel"
	"github.com/muonsoft/go-razdel/internal/testkit"
)

const fuzzMaxInputBytes = 24 * 1024

func FuzzTokenize(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("a b"))
	f.Add([]byte("привет, мир!"))
	f.Add([]byte("что-то и К_тому_же"))
	f.Add([]byte("1,5 и 1/2 :)"))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > fuzzMaxInputBytes {
			data = data[:fuzzMaxInputBytes]
		}
		if !utf8.Valid(data) {
			t.Skip()
		}
		s := string(data)
		toks := razdel.Tokenize(s)
		if err := testkit.ValidateTokenOffsets(s, toks); err != nil {
			t.Fatal(err)
		}
	})
}

func FuzzSentenize(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("   \n\t  "))
	f.Add([]byte("Одно. Два."))
	f.Add([]byte("a. B"))
	f.Add([]byte("Т. е. он пришёл."))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > fuzzMaxInputBytes {
			data = data[:fuzzMaxInputBytes]
		}
		if !utf8.Valid(data) {
			t.Skip()
		}
		s := string(data)
		sents := razdel.Sentenize(s)
		if err := testkit.ValidateSentenceOffsets(s, sents); err != nil {
			t.Fatal(err)
		}
	})
}
