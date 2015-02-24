package templates

import (
	"html/template"
	"io"

	"gum/config"
	"gum/plugins"
)

type Template struct {
	Layout   string
	Template string
}

type DefaultMember struct {
	Title      string
	IsLoggedIn bool
	User       string
}

func init() {
}
func Del() {
}

func (this *Template) Render(w io.Writer, member Member) error {

	return template.Must(template.New(this.Layout).Funcs(map[string]interface{}{
		"linkCSS":    linkCSS,
		"embedImage": embedImage,
		"linkJS":     linkJS,
		"plugin":     plugin,
	}).ParseFiles(config.LayoutsPath+this.Layout, config.TemplatesPath+this.Template)).Execute(w, member)

	return nil
}

func linkCSS(cssFile string) template.HTML {
	return template.HTML("<link rel='stylesheet' href='/" + config.CssPath + cssFile + "' type='text/css' />")
}
func embedImage(imgFile string, alt string) template.HTML {
	return template.HTML("<img alt='" + alt + "' src='/" + config.ImgPath + imgFile + "' />")
}
func linkJS(jsFile string) template.HTML {
	return template.HTML("<script type='text/javascript' src='/" + config.JsPath + jsFile + "' ></script>")
}
func plugin(name string) template.HTML {
	return plugins.Plugins[name]()
}
