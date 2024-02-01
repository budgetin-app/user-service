package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

var isProduction bool

func InitLogger() {
	// Check app environment variables
	isProduction = os.Getenv("APP_ENV") == "production"

	// Configure the log format
	log.SetLevel(getLevel())
	log.SetReportCaller(true)
	log.SetFormatter(&nested.Formatter{
		HideKeys:              true,
		ShowFullLevel:         false,
		TrimMessages:          true,
		NoColors:              isProduction,
		TimestampFormat:       time.StampMilli,
		FieldsOrder:           []string{"component", "category"},
		CallerFirst:           true,
		CustomCallerFormatter: getCustomCallerFormatter,
	})

	// Redirect logs to a file in production
	if isProduction {
		// Generate log file name and path
		logFilePath := "./logs/" + generateLogFileName()
		logDir := filepath.Dir(logFilePath)

		// Create the directory if it doesn't exist
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Error("Failed to create directory: ", err)
			return
		}

		// Now you can open the log file
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Error("Failed to open log file: ", err)
		}
	}
}

// getLevel get log level according to the specified environment variables
func getLevel() log.Level {
	// When it's debuggable, then just show all the log level. Otherwise, show
	// according to the environment variable LOG_LEVEL
	if os.Getenv("APP_DEBUGABLE") == "true" {
		// Show all level
		return log.InfoLevel
	} else {
		switch os.Getenv("LOG_LEVEL") {
		case "DEBUG":
			return log.DebugLevel
		case "TRACE":
			return log.TraceLevel
		case "INFO":
		default:
			return log.InfoLevel
		}
	}
	return log.InfoLevel
}

// getCustomCallerFormatter format the caller function file path and line
func getCustomCallerFormatter(caller *runtime.Frame) string {
	// On production APP_ENV the caller file will present,
	// otherwise APP_ENV local for example, the caller detail will not returned
	if isProduction {
		return fmt.Sprintf(" (%s:%d %s)", caller.File, caller.Line, caller.Function)
	} else {
		return ""
	}
}

// generateLogFileName generate log filename
func generateLogFileName() string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return fmt.Sprintf("app_logs_%s.log", timestamp)
}
