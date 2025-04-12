//go:build !prod

package logger

import (
	"io"
	"log"
	"os"
)

// Setup: 初期化
func Setup(_ string) io.Writer {
	log.SetFlags(log.Ldate | log.Ltime)
	return os.Stdout
}
