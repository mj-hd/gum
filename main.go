package gum

import (
	"net/http"
	"os"

	"gum/config"
	"gum/controllers"
	"gum/models"
	"gum/plugins"
	"gum/templates"
	"gum/utils"

	"github.com/gorilla/context"
)

func init() {

	utils.LogFile = config.LogFile
	utils.DisplayLog = config.DisplayLog
	utils.LogLevel = config.LogLevel

}

func Del() {
  models.Del()
  controllers.Del()
  templates.Del()
  plugins.Del()
}

func Start() {
	for route := range controllers.Router.Iterator() {
		http.HandleFunc(route.Path, route.Function)
		utils.PromulgateDebugStr(os.Stdout, route.Path+"に関数を割当")
	}

	http.Handle("/"+config.StaticPath, http.StripPrefix("/"+config.StaticPath, http.FileServer(http.Dir(config.StaticPath))))
	utils.PromulgateDebugStr(os.Stdout, "/"+config.StaticPath+"に静的コンテンツを割当")

	utils.PromulgateInfoStr(os.Stdout, "ポート"+config.ServerPort+"でサーバを開始...")
	http.ListenAndServe(":"+config.ServerPort, context.ClearHandler(http.DefaultServeMux))
}
