//go:build prod

package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
)

// Setup: 初期化
func Setup(appName string) io.Writer {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return os.Stdout
	}

	fileName := filepath.Join(configDir, appName, "logs", "app.log")

	fileLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     14,
		Compress:   true,
	}

	log.SetFlags(log.Ldate | log.Ltime)

	return fileLogger
}
