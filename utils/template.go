package utils

import (
	"text/template"

	alertTemplate "github.com/prometheus/alertmanager/template"
	"github.com/sirupsen/logrus"
)

type Notification struct {
	Receiver          string           `json:"receiver"`
	Status            string           `json:"status"`
	Alerts            Alerts           `json:"alerts"`
	GroupLabels       alertTemplate.KV `json:"groupLabels"`
	CommonLabels      alertTemplate.KV `json:"commonLabels"`
	CommonAnnotations alertTemplate.KV `json:"commonAnnotations"`
	ExternalURL       string           `json:"externalURL"`
	Version           string           `json:"version"`
	GroupKey          string           `json:"groupKey"`
}

type Alerts alertTemplate.Alerts

func ParseLocalTemplate(tmplDir string) (*template.Template, error) {
	t, err := template.New("local").Funcs(template.FuncMap(alertTemplate.DefaultFuncs)).ParseGlob(tmplDir)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func GetLocalTemplate(tmplDir string) *template.Template {
	t, err := ParseLocalTemplate(tmplDir)
	if err != nil {
		logrus.Errorf("Failed to parse template Dir %s with error %v", tmplDir, err)
		return nil
	}
	return t
}
