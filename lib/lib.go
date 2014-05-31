package lib

import (
	"fmt"
	"math/rand"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// random helper functions
func GC() {
	runtime.GC()
}

func CommaNum(num string) string {
	l := len(num)
	n := l / 3
	ret := ""
	for n > 0 {
		n--
		ret = num[l-3:l] + "," + ret
		l = l - 3
		num = num[0:l]
	}
	if num != "" {
		ret = num + "," + ret
	}
	return ret[0 : len(ret)-1]
}

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

func RandSelect(choices []string) *string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &choices[r.Intn(len(choices))]
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

func TimeTrack(start time.Time, name string) string {
	elapsed := time.Since(start)
	return fmt.Sprintf("%s took %s", name, elapsed)
}

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
