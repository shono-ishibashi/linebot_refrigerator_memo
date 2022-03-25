package config

import (
	"linebot/utils"
	"os"
)

type ConfigList struct {
	DbDriverName   string
	DbName         string
	DbUserName     string
	DbUserPassword string
	DbHost         string
	DbPort         string
	ServerPort     string
	LogFile        string
	ChannelSecret  string
	ChannelToken   string
}

var Config ConfigList

func init() {
	loadConfig()
	utils.LoggingSettings(Config.LogFile)
}

func loadConfig() {
	Config = ConfigList{
		DbDriverName:   "postgres",
		DbName:         os.Getenv("POSTGRES_DB"),
		DbUserName:     os.Getenv("POSTGRES_USER"),
		DbUserPassword: os.Getenv("POSTGRES_PASSWORD"),
		DbHost:         os.Getenv("DB_HOSTNAME"),
		DbPort:         "5432",
		ServerPort:     os.Getenv("PORT"),
		LogFile:        "webapp.log",
		ChannelSecret:  os.Getenv("CHANNEL_SECRET"),
		ChannelToken:   os.Getenv("CHANNEL_TOKEN"),
	}
}
