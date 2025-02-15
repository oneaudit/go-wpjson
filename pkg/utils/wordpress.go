package utils

import "strings"

func ExtractURLPathParameters(wpURL string) string {
	const separator1 = "(?P"
	const separator2 = "\u003c"
	const separator3 = "\u003e"
	var result strings.Builder
	i := 0
	for i < len(wpURL) {
		if i+len(separator1) < len(wpURL) && wpURL[i:i+len(separator1)] == separator1 {
			i += len(separator1)
			i += len(separator2)
			result.WriteRune('<')

			for i+len(separator3) < len(wpURL) && wpURL[i:i+len(separator3)] != separator3 {
				result.WriteByte(wpURL[i])
				i += 1
			}
			i += len(separator3)
			result.WriteRune('>')

			depth := 1
			for i < len(wpURL) && depth != 0 {
				if wpURL[i] == '(' {
					depth++
				} else if wpURL[i] == ')' {
					depth--
				}
				i++
			}
			i--
		} else {
			result.WriteByte(wpURL[i])
		}
		i++
	}
	return result.String()
}
