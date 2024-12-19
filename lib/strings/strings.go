package strings

import (
	"strings"
)

const (
	SpaceChar   = " "
	EmptyString = ""
)

func IsEmpty(str string) bool {
	return strings.Trim(str, SpaceChar) == EmptyString
}
