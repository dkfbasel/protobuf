package plugin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// tagsMatch is used to find tags in trailing inline comments
var tagsMatch = regexp.MustCompile("^(`{1}.*`{1})")

// tagInfo is used for parsed comment information
type tagInfo struct {
	Input  string
	Output string
	Parts  struct {
		Comment string
		Tags    string
	}
	Tags map[string]string
}

// String will print the tags
func (t tagInfo) String() string {
	tmp := []string{}
	for key, value := range t.Tags {
		tmp = append(tmp, fmt.Sprintf(`%s:::%s`, key, value))
	}
	return strings.Join(tmp, " ")
}

// check if the struct should be embedded
func (t tagInfo) IsEmbedded() bool {
	return t.Tags["compose"] == "embed"
}

// parseTags will parse any inline comments and check for custom tags
func parseTags(comment string) tagInfo {

	dta := tagInfo{}

	output := strings.Replace(comment, "\n", "", -1)
	output = strings.TrimSpace(output)

	dta.Input = output
	dta.Output = output

	// find all items enclosed in backticks
	match := tagsMatch.FindStringIndex(output)

	if match == nil {
		dta.Output = fmt.Sprintf(" // %s", dta.Output)

	} else {
		dta.Parts.Tags = output[match[0]:match[1]]
		dta.Parts.Comment = strings.TrimSpace(output[match[1]:])
		dta.Output = fmt.Sprintf(" %s", dta.Parts.Tags)

		if dta.Parts.Comment != "" {
			dta.Output = fmt.Sprintf("%s // %s", dta.Output, dta.Parts.Comment)
		}
	}

	// parse tags
	dta.Tags = scanTags(dta.Parts.Tags)
	return dta
}

const keyPart = "key"
const valuePart = "value"

// scanTags will scan the given string to find tags
func scanTags(tags string) map[string]string {

	// initialize the output map
	out := make(map[string]string)

	// trim the tags string
	tags = strings.TrimSpace(tags)
	tags = strings.Trim(tags, "`")

	// return empty map
	if tags == "" {
		return out
	}

	// initialize variables to handle scanning
	var key string
	var value string

	var status = keyPart
	var isEscaped bool

	// scan tags as runes
	for _, runeValue := range tags {

		if status == keyPart {

			if runeValue == ':' {
				status = valuePart
				continue
			}

			if runeValue == ' ' {
				continue
			}

			key += string(runeValue)
			continue

		}

		if status == valuePart {

			// add escaped character
			if isEscaped == true {
				value += "\\" + string(runeValue)
				isEscaped = false
				continue
			}

			// check if escape character is used
			if runeValue == '\\' {
				isEscaped = true
				continue
			}

			value += string(runeValue)

			if runeValue == '"' && len(value) > 1 {
				// store key and value
				out[key], _ = strconv.Unquote(value)

				// reset key and value
				key = ""
				value = ""

				// start again with the key part
				status = keyPart
			}

		}

	}

	return out

}
