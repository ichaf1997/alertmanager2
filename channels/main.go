package channels

import (
	"text/template"
)

type Medium interface {
	SendMsg(s string) error
}

type Channel struct {
	template *template.Template
}

func NewChannel(tmpl *template.Template) *Channel {
	return &Channel{
		template: tmpl,
	}
}
