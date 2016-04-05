package pkg

import (
	"fmt"
	"strings"
)

type StringSlice []string

func (f *StringSlice) Set(value string) error {
	var s StringSlice
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimLeft(item, " [\"")
		item = strings.TrimRight(item, " \"]")
		s = append(s, item)
	}
	*f = s
	return nil
}

func (f *StringSlice) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *StringSlice) Value() []string {
	return *f
}

func (f *StringSlice) Get() interface{} {
	return *f
}
