package auth

import (
	"math/rand"
	"time"
)

func GetRandomEntry[K comparable, V any](m map[K]V) (K, V, bool) {
	// 要素が空なら false を返す
	if len(m) == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	// キーをスライスに格納
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// ランダムにインデックスを選択
	rand.Seed(time.Now().UnixNano())
	randomKey := keys[rand.Intn(len(keys))]

	return randomKey, m[randomKey], true
}

func MapactionsFunc() map[string]string {
	var Mapactions = map[string]string{}
	return Mapactions
}
