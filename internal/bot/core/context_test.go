package core

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
)

func TestContextBot_GetCommandArg(t *testing.T) {
	tests := []struct {
		name string
		Req  tgbotapi.Update
		want []string
	}{
		{
			name: "Case 1: Check command arg",
			Req:  tgbotapi.Update{Message: &tgbotapi.Message{Text: "/root id"}},
			want: []string{"id"},
		},
		{
			name: "Case 2: Text empty",
			Req:  tgbotapi.Update{Message: &tgbotapi.Message{Text: ""}},
			want: []string{},
		},
		{
			name: "Case 3: No args",
			Req:  tgbotapi.Update{Message: &tgbotapi.Message{Text: "/add"}},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ContextBot{
				Req: tt.Req,
			}

			actual := r.GetCommandArg()

			require.Equal(t, tt.want, actual)
		})
	}
}

func TestContextBot_GetRoutersInfoNonePrivate(t *testing.T) {
	type fields struct {
		bot         *Bot
		event       Event
		handler     HandlerFunc
		description string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Case 1: Check new router",
			fields: fields{
				bot: &Bot{},
				event: Event{
					Form:    Command,
					Command: "/root",
					Private: false,
				},
				handler:     nil,
				description: "New Router",
			},
			want: "Command /root - New Router\n",
		},
		{
			name: "Case 2: Check router is private",
			fields: fields{
				bot: &Bot{},
				event: Event{
					Form:    Command,
					Command: "/root",
					Private: true,
				},
				handler:     nil,
				description: "New Router",
			},
			want: "",
		},
		{
			name: "Case 3: Check router is secret",
			fields: fields{
				bot: &Bot{},
				event: Event{
					Form:    Secret,
					Command: "/root",
					Private: true,
				},
				handler:     nil,
				description: "New Router",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ContextBot{
				bot: tt.fields.bot,
			}
			r.bot.AddRoute(tt.fields.event, tt.fields.handler, tt.fields.description)

			actual := r.GetRoutersInfoNonePrivate()

			require.Equal(t, tt.want, actual)
		})
	}

	t.Run("Case 4: No route", func(t *testing.T) {
		r := &ContextBot{
			bot: &Bot{},
		}

		actual := r.GetRoutersInfoNonePrivate()

		require.Equal(t, "", actual)
	})
}
