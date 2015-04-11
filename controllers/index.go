package controllers

import (
	"net/http"
	"os"

	"gum/config"
	"gum/templates"
	"gum/utils/log"
)

func indexHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template

	tmpl.Layout = "default.tmpl"
	tmpl.Template = "index.tmpl"

	err := tmpl.Render(document, &templates.DefaultMember{
		Title: config.SiteTitle,
		User:  getSessionUser(request),
	})
	if err != nil {
		log.Fatal(os.Stdout, err)
		showError(document, request, "ページの表示に失敗しました。")
		return
	}
}
