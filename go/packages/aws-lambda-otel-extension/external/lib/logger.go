package lib

import (
	"encoding/json"
	"log"
	"os"

	"go.uber.org/zap"
)

type Logger = zap.Logger

func PrettyPrint(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return ""
	}
	return string(data)
}

func NewLogger() (logger *zap.Logger) {
	var err error
	stage, ok := os.LookupEnv("SLS_DEBUG_EXTENSION")
	if !ok || stage == "" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	defer logger.Sync()

	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	return logger
}

var ErrLogger = log.New(os.Stderr, "", 0)
