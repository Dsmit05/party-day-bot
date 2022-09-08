package benchmark

import (
	"testing"
	_ "unsafe"

	_ "github.com/Dsmit05/party-day-bot/internal/bot/core"
)

//go:linkname nameOfFunction github.com/Dsmit05/party-day-bot/internal/bot/core.nameOfFunction
func nameOfFunction(f interface{}) string

func Benchmark_nameOfFunction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nameOfFunction(func() {})
	}
}
