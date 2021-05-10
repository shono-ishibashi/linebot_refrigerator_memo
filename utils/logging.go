package utils

import (
	"io"
	"log"
	"os"
)

var Logger *log.Logger

func LoggingSettings(logFile string) {
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	//logの書き込み先をファイルと　標準出力に指定
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	Logger = log.New(multiLogFile, "[LINE BOT]", log.Ldate|log.Ltime|log.Lshortfile)
}
