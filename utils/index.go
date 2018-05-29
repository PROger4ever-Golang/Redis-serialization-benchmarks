package utils

import (
	"fmt"
	"os"
)

func GetExistingFiles(files ...string) (existingFiles []string, err error) {
	existingFiles = make([]string, 0)
	for _, file := range files {
		_, errStat := os.Stat(file)
		if errStat == nil {
			existingFiles = append(existingFiles, file)
		} else if !os.IsNotExist(errStat) {
			err = WrapIfError(errStat, "utils->GetExistingFiles")
			return
		}
	}
	return
}

func WrapIfError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s\n%s", msg, err)
}
