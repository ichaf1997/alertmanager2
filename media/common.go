package media

import (
	"text/template"
)

type Medium interface {
	SendMsg() error
	Info() any
}

type Channel struct {
	template *template.Template
}

type Response struct {
	RequestInfo any    `json:"request_info"`
	Status      string `json:"status"`
}

func NewChannel(tmpl *template.Template) *Channel {
	return &Channel{
		template: tmpl,
	}
}
