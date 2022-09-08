package benchmark

import (
	"testing"
	_ "unsafe"

	"github.com/Dsmit05/party-day-bot/internal/bot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//go:linkname checkEvent github.com/Dsmit05/party-day-bot/internal/bot/core.(*Bot).checkEvent
func checkEvent(b *core.Bot, update tgbotapi.Update) (*core.Event, error)

func Benchmark_checkEvent(b *testing.B) {
	bot := &core.Bot{}
	update := tgbotapi.Update{Message: &tgbotapi.Message{Text: "Hi"}}
	for i := 0; i < b.N; i++ {
		checkEvent(bot, update)
	}
}
