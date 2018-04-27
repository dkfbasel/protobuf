package main

import "strings"

// dbNameSnakeCaseMapper implements a mapping function to map struct
// names using snake casing instead of default lowercase for sqlx named mappings
func dbNameSnakeCaseMapper(in string) string {

	// TODO: improve using strings.Builder

	out := ""

	previousIsUpper := false

	for pos, value := range in {

		if pos == 0 {
			out = string(value)
			continue
		}

		char := string(value)
		isUpperCase := strings.ToUpper(char) == char

		if isUpperCase && !previousIsUpper {
			out = out + "_"
		}

		out = out + char

	}

	return strings.ToLower(out)

}
