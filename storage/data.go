package storage

import (
	"fmt"
	"os"
	"path"
	"time"
)

func getBackupDataLocation(configurationName string) (string, error) {
	fullPath := path.Join(AppDir, "data", configurationName)
	err := os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		return "", nil
	}

	return fullPath, nil
}

func GetNewBackupPath(configurationName string) (string, error) {
	location, err := getBackupDataLocation(configurationName)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().Format("2006-01-02 15-04-05")
	backupPath := path.Join(location, fmt.Sprintf("%v.sql", timestamp))

	return backupPath, nil
}
