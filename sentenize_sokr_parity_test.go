package razdel_test

import (
	"testing"

	"github.com/muonsoft/go-razdel"
)

// Upstream test_sentenize.py UNIT — SOKR section; expected chunks match razdel.sentenize (Python).
func TestSentenize_upstreamSokrUnitParity(t *testing.T) {
	t.Parallel()
	cases := []struct {
		text string
		want []string
	}{
		{
			"И т. д. и т. п. В общем, вся газета",
			[]string{"И т. д. и т. п.", "В общем, вся газета"},
		},
		{
			"специалистом, к.п.н. И. П. Карташовым.",
			[]string{"специалистом, к.п.н. И. П. Карташовым."},
		},
		{
			"основании п. 2, ст. 5 УПК",
			[]string{"основании п. 2, ст. 5 УПК"},
		},
		{
			"Вблизи оз. Селяха",
			[]string{"Вблизи оз. Селяха"},
		},
		{
			"уменьшить с 20 до 18 проц. (при сохранении",
			[]string{"уменьшить с 20 до 18 проц. (при сохранении"},
		},
		{
			"6 июля 2007 г. \"в связи с совершением",
			[]string{"6 июля 2007 г. \"в связи с совершением"},
		},
		{
			"на 500 тыс. машин",
			[]string{"на 500 тыс. машин"},
		},
		{
			"Влияние взглядов Л. В. Щербы",
			[]string{"Влияние взглядов Л. В. Щербы"},
		},
		{
			"директор фирмы Чарльз Дж. Филлипс",
			[]string{"директор фирмы Чарльз Дж. Филлипс"},
		},
		{
			"Т.е. ОБЯЗАТЕЛЬНО письменно",
			[]string{"Т.е. ОБЯЗАТЕЛЬНО письменно"},
		},
		{
			"была утечка т.н. Таблицы боевых действий",
			[]string{"была утечка т.н. Таблицы боевых действий"},
		},
		{
			"В 1996-1999гг. теффт",
			[]string{"В 1996-1999гг. теффт"},
		},
		{
			"России, т. е. 55 % опрошенных",
			[]string{"России, т. е. 55 % опрошенных"},
		},
		{
			"я ощущал в 1990-е. Славное было время",
			[]string{"я ощущал в 1990-е.", "Славное было время"},
		},
	}
	for _, tc := range cases {
		got := razdel.Sentenize(tc.text)
		if len(got) != len(tc.want) {
			t.Fatalf("%q: len got %d want %d (%#v)", tc.text, len(got), len(tc.want), got)
		}
		for i := range tc.want {
			if got[i].Text != tc.want[i] {
				t.Fatalf("%q: [%d] got %q want %q", tc.text, i, got[i].Text, tc.want[i])
			}
		}
	}
}
