package media

import (
	"text/template"
)

type Medium interface {
	SendMsg() error
}

type Handler struct {
	template *template.Template
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewHandler(tmpl *template.Template) *Handler {
	return &Handler{
		template: tmpl,
	}
}
