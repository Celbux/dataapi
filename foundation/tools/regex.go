package tools

import (
	"regexp"
)

func StripField(data string, regex string) string {
	match, _ := regexp.MatchString(regex, data)
	if !match {
		return ""
	}

	r, _ := regexp.Compile(regex)
	matches := r.FindStringSubmatch(data)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
