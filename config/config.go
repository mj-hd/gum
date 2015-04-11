package config

import "gum/utils/log"

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
	LogLevel   int    = log.Level_Debug
	LogFile    string = "log"
	SessionKey string = "sample"
)
