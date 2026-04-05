// Code generated from third_party/razdel/razdel/segmenters/sokr.py — keep in sync with upstream.
package sentenize

var (
	// sokrWordsRaw is upstream SOKRS (TAIL ∪ HEAD ∪ OTHER), whitespace-separated.
	sokrWordsRaw = `
		al bd chap co corp dr ed inc mr mrs ms no ps upd
		vs а абз акад англ араб арх ауд барр букв в вв венг внутр
		г га гг гл гор гос гр греч д дес дифф дол долл домовлад
		доп др евр зав зам изд им искл исп итал к каб кат кв
		кит км кн ком комн коп корп корр куб лат латв лит мин млн
		млрд напр нач нем о обл обр оз ок откр отм оф п пер
		пл пол пом пос пп пр прим просп проф проц пс пт р ред
		рис рп руб рус русск с сб св сек слав словацк см сокр ср
		ст стр т тел тов трад тыс укр ул ум устар физ фр х
		хорв ч час чл чч ш шутл эт юр яз яп
	`
	// headSokrWordsRaw is upstream HEAD_SOKRS.
	headSokrWordsRaw = `
		bd chap dr mr mrs ms no ps upd vs а абз акад англ
		араб арх ауд букв венг внутр г гл гор гос гр греч д дифф
		домовлад доп евр зав зам им исп итал к каб кат кит кн ком
		комн корп корр лат латв лит напр нач нем о обл обр оз ок
		откр отм оф п пер пл пол пом пос пп пр просп проф пс
		пт р ред рп рус русск с сб св слав словацк см ср ст
		стр т тел тов трад укр ул ум физ фр х хорв ч чл
		чч ш эт юр яп
	`
	// pairSokrPairs is upstream PAIR_SOKRS (2-tuples only; upstream also has an accidental 3-tuple unused by rules).
	pairSokrPairs = [][2]string{
		{"a", "m"},
		{"p", "m"},
		{"ед", "ч"},
		{"з", "д"},
		{"и", "о"},
		{"к", "н"},
		{"к", "п"},
		{"к", "т"},
		{"л", "д"},
		{"л", "с"},
		{"мн", "ч"},
		{"н", "э"},
		{"п", "н"},
		{"повел", "накл"},
		{"р", "х"},
		{"с", "г"},
		{"с", "ш"},
		{"т", "д"},
		{"т", "е"},
		{"т", "к"},
		{"т", "н"},
		{"т", "п"},
		{"у", "е"},
		{"ч", "т"},
	}
	// headPairSokrPairs is upstream HEAD_PAIR_SOKRS.
	headPairSokrPairs = [][2]string{
		{"и", "о"},
		{"к", "н"},
		{"к", "п"},
		{"к", "т"},
		{"л", "д"},
		{"п", "н"},
		{"т", "е"},
		{"т", "к"},
		{"т", "н"},
	}
)
