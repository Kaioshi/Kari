package lib

import (
	"regexp"
	"strings"
)

// helper functions

func Sanitise(line string) string {
	if !strings.ContainsAny(line, "\n\t\r") {
		return line
	}
	reg := regexp.MustCompile("\\n|\\t|\\r")
	return reg.ReplaceAllString(line, "")
}

func StripHtml(html string) string {
	reg := regexp.MustCompile("<[^<]+?>|\\n|\\t|\\r")
	return reg.ReplaceAllString(html, "")
}

func HasElementString(arr []string, match string) bool {
	for _, line := range arr {
		if strings.ToLower(line) == strings.ToLower(match) {
			return true
		}
	}
	return false
}
