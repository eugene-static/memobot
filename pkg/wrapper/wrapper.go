package wrapper

import (
	"fmt"
	"strings"
)

const (
	sep = ":"
)

func Wrap(s1, s2 string) string {
	return fmt.Sprintf("%s%s%s", s1, sep, s2)
}

func Unwrap(s string) (string, string) {
	data := strings.Split(s, sep)
	if len(data) != 2 {
		return "", ""
	}
	return data[0], data[1]
}
