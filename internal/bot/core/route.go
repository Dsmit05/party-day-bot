package core

import (
	"bytes"
	"fmt"

	"github.com/Dsmit05/party-day-bot/internal/logger"
)

type HandlerFunc func(request ContextBot) error

type Form string

type Event struct {
	Form    Form
	Command string
	Private bool
}

func (e *Event) IsEqualRoute(new *Event) bool {
	return e.Form == new.Form && e.Command == new.Command
}

const (
	Photo    Form = "Photo"
	Document Form = "Document"
	Command  Form = "Command"
	Text     Form = "Text"
	Unknown  Form = "Unknown"
	Secret   Form = "Secret"
)

type Route struct {
	Event       Event
	HandlerName string
	Description string
	HandlerFunc HandlerFunc
}

type Routes []Route

func (r Routes) String() string {
	var buf bytes.Buffer
	for _, v := range r {
		buf.WriteString("Command: ")
		buf.WriteString(string(v.Event.Form))
		buf.WriteString(v.Event.Command)
		buf.WriteString(" Description: ")
		buf.WriteString(v.Description)
		buf.WriteString("\n")
	}

	return buf.String()
}

func (r Routes) Log() {
	for _, v := range r {
		info := fmt.Sprintf("%v -> %v - %v", v.HandlerName, v.Event.Form, v.Event.Command)
		logger.Info("Router init", info)
	}
}
