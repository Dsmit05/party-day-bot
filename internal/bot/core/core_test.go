package core

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestBot_checkEvent(t *testing.T) {
	tests := []struct {
		name    string
		arg     tgbotapi.Update
		want    *Event
		wantErr error
	}{
		{
			name:    "Case 1: Empty tgbotapi.Update",
			arg:     tgbotapi.Update{},
			want:    nil,
			wantErr: errors.New("unknown event"),
		},
		{
			name:    "Case 2: Empty Message text from tgbotapi.Update",
			arg:     tgbotapi.Update{Message: &tgbotapi.Message{Text: ""}},
			want:    nil,
			wantErr: errors.New("unknown event"),
		},
		{
			name: "Case 3: Get Photo from tgbotapi.Update",
			arg:  tgbotapi.Update{Message: &tgbotapi.Message{Photo: []tgbotapi.PhotoSize{}}},
			want: &Event{
				Form:    "Photo",
				Command: "",
				Private: false,
			},
			wantErr: nil,
		},
		{
			name: "Case 4: Get Document from tgbotapi.Update",
			arg:  tgbotapi.Update{Message: &tgbotapi.Message{Document: &tgbotapi.Document{}}},
			want: &Event{
				Form:    "Document",
				Command: "",
				Private: false,
			},
			wantErr: nil,
		},
		{
			name: "Case 5: Get text Command from tgbotapi.Update",
			arg:  tgbotapi.Update{Message: &tgbotapi.Message{Text: "/add"}},
			want: &Event{
				Form:    "Command",
				Command: "/add",
				Private: false,
			},
			wantErr: nil,
		},
		{
			name: "Case 6: Get text secret Command from tgbotapi.Update",
			arg:  tgbotapi.Update{Message: &tgbotapi.Message{Text: "/secret"}},
			want: &Event{
				Form:    "Secret",
				Command: "",
				Private: false,
			},
			wantErr: nil,
		},
		{
			name: "Case 7: Get text len > 10 from tgbotapi.Update",
			arg:  tgbotapi.Update{Message: &tgbotapi.Message{Text: "Hi, you cool!!!"}},
			want: &Event{
				Form:    "Text",
				Command: "",
				Private: false,
			},
			wantErr: nil,
		},
		{
			name: "Case 8: Get text not command and len < 10 from tgbotapi.Update",
			arg:  tgbotapi.Update{Message: &tgbotapi.Message{Text: "Hi"}},
			want: &Event{
				Form:    "Unknown",
				Command: "",
				Private: false,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{config: Config{SecretCMD: "/secret"}}

			actual, err := b.checkEvent(tt.arg)

			require.Equal(t, tt.want, actual)

			if err != nil {
				if tt.wantErr == nil {
					t.Errorf("checkEvent() return error: %v, but wantErr is nil", err)
				} else {
					require.EqualError(t, err, tt.wantErr.Error())
				}
			}
		})
	}
}

func Test_nameOfFunction(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want string
	}{
		{
			name: "Case 1: anon func",
			arg:  func() {},
			want: "gitlab.ozon.dev/Dsmit05/party-day-bot/internal/bot/core.Test_nameOfFunction.func1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := nameOfFunction(tt.arg)

			require.Equal(t, tt.want, actual)
		})
	}
}
