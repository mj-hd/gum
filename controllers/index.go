package controllers

import (
	"net/http"

	"../config"
	"../templates"
)

type indexMember struct {
	*templates.DefaultMember
}

func indexHandler(document http.ResponseWriter, request *http.Request) (err error) {

	var tmpl templates.Template

	tmpl.Layout = "default.tmpl"
	tmpl.Template = "index.tmpl"

	return tmpl.Render(document, indexMember{
		DefaultMember: &templates.DefaultMember{
			Title:  config.SiteTitle,
			UserID: getSessionUser(request),
		},
	})
}
