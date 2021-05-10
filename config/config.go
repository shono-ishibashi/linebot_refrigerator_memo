package config

import (
	"gopkg.in/go-ini/ini.v1"
	"linebot/utils"
	"log"
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
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	Config = ConfigList{
		DbDriverName:   "postgres",
		DbName:         os.Getenv("POSTGRES_DB"),
		DbUserName:     os.Getenv("POSTGRES_USER"),
		DbUserPassword: os.Getenv("POSTGRES_PASSWORD"),
		DbHost:         os.Getenv("DB_HOSTNAME"),
		DbPort:         "5432",
		ServerPort:     os.Getenv("PORT"),
		LogFile:        "webapp.log",
		ChannelSecret:  cfg.Section("line_config").Key("channel_secret").String(),
		ChannelToken:   cfg.Section("line_config").Key("channel_token").String(),
	}
}
