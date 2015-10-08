package config

import "time"

const (
	TemplatesPath string = "templates/"
	LayoutsPath   string = "templates/layouts/"
	StaticPath    string = "static/"
	CssPath       string = "static/css/"
	ImgPath       string = "static/img/"
	JsPath        string = "static/js/"
)

// DIRTY:
func JST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}
