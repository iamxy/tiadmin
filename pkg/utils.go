package pkg

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

const (
	KC_RAND_KIND_NUM   = 0
	KC_RAND_KIND_LOWER = 1
	KC_RAND_KIND_UPPER = 2
	KC_RAND_KIND_ALL   = 3
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

func KRand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
