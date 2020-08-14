package logger

import (
	"io"
	"log"
)

// BuildCommit Git commit
var BuildCommit string

// BuildBranch Git branch
var BuildBranch string

// BuildTag Git tag
var BuildTag string

// BuildDate Build date
var BuildDate string

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// Logger
func LoggerInit(infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	Info = log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime)
	Warning = log.New(warningHandle, "[WARN] ", log.Ldate|log.Ltime)
	Error = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime)
}
