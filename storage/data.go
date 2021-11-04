package storage

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"sort"
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

func GetBackups(configurationName string) ([]fs.FileInfo, string, error) {
	location, err := getBackupDataLocation(configurationName)
	if err != nil {
		return []fs.FileInfo{}, "", err
	}

	dir, err := ioutil.ReadDir(location)
	if err != nil {
		return []fs.FileInfo{}, "", err
	}

	sort.Slice(dir, func(i, j int) bool {
		return dir[i].ModTime().After(dir[j].ModTime())
	})

	return dir, location, nil
}
