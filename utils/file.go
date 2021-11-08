package utils

import "os"

func GetFileByName(files []os.FileInfo, search string) os.FileInfo {
	for _, el := range files {
		if el.Name() == search {
			return el
		}
	}

	return nil
}