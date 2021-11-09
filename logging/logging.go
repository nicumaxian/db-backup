package logging

import (
	"db-backup/utils"
	"log"
	"os"
)

var Verbose = false
var (
	empty   = log.New(new(utils.NullWriter), "", 0)
	builtin = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

func GetLoggerByVerbose() *log.Logger {
	if Verbose {
		return builtin
	}
	return empty
}
