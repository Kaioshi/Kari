package lib

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// STRING LISTS ---- why isn't this built in or less pointlessly complex
type StrList struct {
	List []string
}

func (sl *StrList) Pop() {
	if len(sl.List) == 0 {
		return
	}
	sl.List = sl.List[:len(sl.List)-1]
}

func (sl *StrList) HasElement(match string, ignoreCase bool) bool {
	if ignoreCase {
		match = strings.ToLower(match)
		for _, element := range sl.List {
			if strings.ToLower(element) == match {
				return true
			}
		}
	} else {
		for _, element := range sl.List {
			if element == match {
				return true
			}
		}
	}
	return false
}

func (sl *StrList) Add(s ...string) {
	sl.List = append(sl.List, s...)
}

func (sl *StrList) RemoveByMatch(match string, ignoreCase bool) {
	if ignoreCase {
		for i, value := range sl.List {
			if strings.ToLower(value) == strings.ToLower(match) {
				sl.RemoveByIndex(i)
				break
			}
		}
	} else {
		for i, value := range sl.List {
			if value == match {
				sl.RemoveByIndex(i)
				break
			}
		}
	}
}

func (sl *StrList) RemoveByIndex(i int) {
	s := sl.List
	s = append(s[:i], s[i+1:]...)
	sl.List = s
}

func (sl *StrList) String() string {
	return fmt.Sprintf("%#v", sl.List)
}

// random helper functions
func Sanitise(line string) string {
	if !strings.ContainsAny(line, "\n\t\r") {
		return line
	}
	reg := regexp.MustCompile("\\n|\\t|\\r")
	return reg.ReplaceAllString(line, "")
}

func SingleSpace(line string) string {
	if strings.Index(line, "  ") > -1 {
		return strings.Join(strings.Fields(line), " ")
	}
	return line
}

func Timestamp(line string) string {
	if line == "" {
		return time.Now().Format("[" + time.Stamp + "]")
	}
	return time.Now().Format("["+time.Stamp+"]") + " " + line
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
