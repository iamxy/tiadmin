package pkg

import (
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
