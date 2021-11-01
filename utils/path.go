package utils

import (
	"fmt"
	"os"
	"path"
	"time"
)

func getBackupDataLocation(name string) (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	join := path.Join(dir, ".db-backup", "data", name)
	err = os.MkdirAll(join, os.ModePerm)
	if err != nil {
		return "", nil
	}

	return join, nil
}

func GetBackupPath(name string) (string, error) {
	location, err := getBackupDataLocation(name)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().Unix()
	backupPath := path.Join(location, fmt.Sprintf("%v.sql", timestamp))

	return backupPath, nil
}
