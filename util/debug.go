package util

import (
	"log"
	"fmt"
	"os"
)

var Dbg DebugLog
func init () {
	Dbg.logger = log.New(os.Stdout, "", log.Lshortfile)
}

type DebugLog struct {
	Enabled bool
	logger *log.Logger
}

func (dl *DebugLog) Printf (format string, v ...interface{}) {
	if dl.Enabled {
		dl.logger.Output(2, fmt.Sprintf(format, v...))
	}
}
