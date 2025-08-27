package logs

import (
	"encoding/json"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *log.Logger
var childLogger = logrus.WithFields(logrus.Fields{
	"service": "interview-tracker-service",
})

func LoggerConfig() {

	logger := log.New(
		os.Stderr,
		"|"+os.Getenv("GLOBAL_ENDPOINT")+"| ",
		log.Flags(),
	)
	Logger = logger
}

func LogJson(info interface{}) {
	jsonBytes, _ := json.Marshal(info)
	Logger.Println("json| ", string(jsonBytes))
}

func LoggerInfo(info string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	childLogger.Info(info)
}
