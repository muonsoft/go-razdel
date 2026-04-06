package sentenize_test

// upstreamSentsTxtPartitionDrift — partition-строки из
// third_party/razdel/razdel/tests/data/sents.txt, у которых эталон (несущие куски между «|»)
// не совпадает с razdel.segmenters.sentenize на том же восстановленном тексте в закреплённом
// submodule. Go совпадает с Python; расхождение только между строкой корпуса и кодом razdel.
//
// Полный корпусный тест пропускает эти строки (с записью в лог), чтобы остальные кейсы
// оставались строгими. После обновления submodule пересоберите список (и уберите/допишите ключи).
var upstreamSentsTxtPartitionDrift = map[string]struct{}{}
