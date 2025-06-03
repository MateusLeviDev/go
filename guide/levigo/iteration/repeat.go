package iteration

import "strings"

const repeatCount = 5

func Repeat(character string) string {
	var repeated strings.Builder
	for i := 0; i < repeatCount; i++ {
		repeated.WriteString(character)
	}
	return repeated.String()
}

func Repeat2(character string) string {
	return strings.Repeat(character, 5)
}
