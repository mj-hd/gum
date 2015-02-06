package config

import "gum/utils"

var (
	RootPath   string = "/path/to/root/of/this/project"
	SiteTitle  string = "SiteTitle"
	ServerPort string = "8080"
	DBHost     string = "localhost"
	DBPort     string = "3306"
	DBName     string = "dbname"
	DBUser     string = "dbuser"
	DBPass     string = "example_passwd"
	DisplayLog bool   = true
	LogLevel   int    = utils.Level_Debug
	LogFile    string = "log"
	SessionKey string = "sample"
)
