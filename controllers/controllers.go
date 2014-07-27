package controllers

import (
	"net/http"

	"github.com/gorilla/sessions"

	"gum/config"
	"gum/templates"
)

var Router Routes
var sessionStore = sessions.NewCookieStore([]byte(config.SessionKey))

func init() {

	Router.Register("/", indexHandler)
	Router.Register("/error/", flashHandler)
	Router.Register("/success/", flashHandler)

}
func Del() {

}

func getSessionUser(request *http.Request) string {
	session, _ := sessionStore.Get(request, "user")
	if session.Values["User"] == nil {
		return "anonymous"
	}

	return session.Values["User"].(string)
}

type Routes struct {
	keys   []string
	values []func(http.ResponseWriter, *http.Request)
}
type Route struct {
	Path     string
	Function func(http.ResponseWriter, *http.Request)
}

func (this *Routes) Iterator() <-chan Route {
	ret := make(chan Route)

	go func() {
		for i, k := range this.keys {
			var route Route
			route.Path = k
			route.Function = this.values[i]

			ret <- route
		}
		close(ret)
	}()

	return ret
}

func (this *Routes) Register(path string, fn func(http.ResponseWriter, *http.Request)) {
	this.keys = append(this.keys, path)
	this.values = append(this.values, fn)
}

func (this *Routes) Value(path string) func(http.ResponseWriter, *http.Request) {
	for i, key := range this.keys {
		if path == key {
			return this.values[i]
		}
	}
	return nil
}

func (this *Routes) Key(fn *func(http.ResponseWriter, *http.Request)) string {
	for i, val := range this.values {
		if fn == &val {
			return this.keys[i]
		}
	}
	return ""
}

func showError(document http.ResponseWriter, request *http.Request, message string) {
	session, _ := sessionStore.Get(request, "user")

	session.AddFlash(message)
	session.Save(request, document)
	http.Redirect(document, request, "/error/", http.StatusSeeOther)
}

type flashMember struct {
	*templates.DefaultMember
	Messages []interface{}
	Referer  string
}

func flashHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "flash.tmpl"

	session, _ := sessionStore.Get(request, "user")

	var message string
	if request.URL.Path == "/error/" {
		message = "エラー"
	} else {
		message = "成功"
	}

	flashes := session.Flashes()
	session.Save(request, document)

	tmpl.Render(document, flashMember{
		DefaultMember: &templates.DefaultMember{
			Title: message,
			User:  session.Values["User"].(string),
		},
		Messages: flashes,
		Referer:  request.Referer(),
	})
}