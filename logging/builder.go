package logging

import (
	"db-backup/utils"
	"log"
	"os"
)

func NewBuiltinLogger() *log.Logger {
	return log.New(os.Stdout, "", 5)
}

func NewMockLogger() *log.Logger {
	return log.New(new(utils.NullWriter), "", 0)
}
