package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"

	"../config"
	"../templates"
	"../utils/log"
)

var Router Routes
var sessionStore = sessions.NewCookieStore([]byte(config.SessionKey))

func Init() {

	apiInit()

	Router.RegisterPage("/", indexHandler)
	Router.Register("/error/", flashHandler)
	Router.Register("/success/", flashHandler)
	Router.RegisterApi("/api/", apiHandler)

}
func Del() {
	apiDel()
}

func getSession(request *http.Request) (*sessions.Session, error) {
	session, err := sessionStore.Get(request, config.SessionName)

	if err == nil {
		// 一週間
		session.Options.MaxAge = 86400 * 7
	}

	return session, err
}

func getSessionUser(request *http.Request) int {
	session, _ := getSession(request)
	if session.Values["User"] == nil {
		return 0
	}

	return session.Values["User"].(int)
}

func removeSession(document http.ResponseWriter, request *http.Request) {
	session, err := getSession(request)
	if err != nil {
		return
	}

	session.Options = &sessions.Options{MaxAge: -1, Path: "/"}
	session.Save(request, document)
}

func writeStruct(document http.ResponseWriter, s interface{}, httpStatus int) {

	var err error

	document.Header().Set("Content-Type", "application/json")
	jso, err := json.Marshal(s)

	if err != nil {

		log.Fatal(err)

		document.WriteHeader(500)
		document.Write([]byte("{ \"Status\" : \"error\", \"Message\" : \"不明のエラーです。\" }"))

		return
	}

	document.WriteHeader(httpStatus)
	document.Write(jso)
}

type apiMember struct {
	Status  string
	Message string
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

func (this *Routes) RegisterPage(path string, fn func(http.ResponseWriter, *http.Request) error) {
	this.keys = append(this.keys, path)
	this.values = append(this.values, func(document http.ResponseWriter, request *http.Request) {
		err := fn(document, request)
		if err != nil {
			log.FatalStr("ページの表示に失敗:")
			log.Fatal(err)

			showError(document, request, "ページの表示中にエラーが発生しました。管理人へ報告してください。")
		}
	})
}

func (this *Routes) RegisterApi(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	this.keys = append(this.keys, path)
	this.values = append(this.values, func(document http.ResponseWriter, request *http.Request) {
		status, err := fn(document, request)
		if err != nil {
			log.FatalStr("APIの実行に失敗:")
			log.Fatal(err)

			writeStruct(document, apiMember{
				Status:  "error",
				Message: err.Error(),
			}, status)

			return
		}
	})
}

func (this *Routes) RegisterPostApi(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	wrapper := func(document http.ResponseWriter, request *http.Request) (int, error) {
		if request.Method != "POST" {
			return http.StatusBadRequest, errors.New("POST以外のメソッド")
		}

		return fn(document, request)
	}

	this.RegisterApi(path, wrapper)
}

func (this *Routes) RegisterGetApi(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	wrapper := func(document http.ResponseWriter, request *http.Request) (int, error) {
		if request.Method != "GET" {
			return http.StatusBadRequest, errors.New("GET以外のメソッド")
		}

		return fn(document, request)
	}

	this.RegisterApi(path, wrapper)
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
	session, _ := getSession(request)

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

	session, _ := getSession(request)

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
			Title:  message,
			UserID: getSessionUser(request),
		},
		Messages: flashes,
		Referer:  request.Referer(),
	})
}
