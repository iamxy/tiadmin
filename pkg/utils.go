package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func GetCmdDir() (string, error) {
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Dir(path), nil
}

func GetRootDir() (string, error) {
	cd, err := GetCmdDir()
	if err != nil {
		return "", err
	}
	if filepath.Base(cd) == "bin" {
		return filepath.Dir(cd), nil
	} else {
		return cd, nil
	}
}

func CheckFileExist(filepath string) (string, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return "", err
	}
	if fi.IsDir() {
		return "", errors.New(fmt.Sprintf("filepath: %s, is a directory, not a file", filepath))
	}
	return filepath, nil
}
