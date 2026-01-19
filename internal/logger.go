package internal

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

var (
	loggerInstance *log.Logger
	once           sync.Once
)

func getLogger() *log.Logger {
	once.Do(func() {
		f, err := os.OpenFile("mcp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		loggerInstance = log.New(f, "", log.LstdFlags)
	})
	return loggerInstance
}

func logJSON(msgType, tool string, v interface{}) {
	logger := getLogger()
	bytes, err := json.Marshal(v)
	if err != nil {
		logger.Printf("[%s][%s] Could not marshal value: %v", msgType, tool, err)
		return
	}
	logger.Printf("[%s][%s][%s] %s", msgType, tool, time.Now().UTC().Format(time.RFC3339), string(bytes))
}

func LogRequest(tool string, req interface{}) {
	logJSON("REQUEST", tool, req)
}

func LogResponse(tool string, res interface{}) {
	logJSON("RESPONSE", tool, res)
}

func LogError(tool string, err error) {
	logger := getLogger()
	logger.Printf("[ERROR][%s][%s] %v", tool, time.Now().UTC().Format(time.RFC3339), err)
}
