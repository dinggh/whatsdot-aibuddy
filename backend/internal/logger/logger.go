package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var defaultDir = "logs"

// Init creates the log directory and configures logging to write to logs/app.log.
// Also writes to stdout for local dev visibility.
func Init(logDir string) (io.Closer, error) {
	if logDir == "" {
		logDir = defaultDir
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	fpath := filepath.Join(logDir, "app.log")
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	w := io.MultiWriter(os.Stdout, f)
	log.SetOutput(w)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return f, nil
}
