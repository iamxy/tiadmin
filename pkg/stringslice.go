package pkg

import (
	"strings"
)

type StringSlice []string

func NewStringSlice(value string) StringSlice {
	var s = []string{}
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimLeft(item, " [\"")
		item = strings.TrimRight(item, " \"]")
		s = append(s, item)
	}
	return StringSlice(s)
}

func (f StringSlice) String() string {
	return strings.Join(f.Value(), ",")
}

func (f StringSlice) Value() []string {
	return []string(f)
}
