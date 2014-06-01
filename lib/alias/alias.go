package alias

import (
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/storage"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var DB storage.StringListDB
var Vars storage.StringListDB

type Event struct {
	From    string
	Context string
	Data    string
	Args    []string
}

func (e *Event) Populate(params *events.Params, args []string) {
	e.Data = strings.Join(args, " ")
	e.Args = args
	e.From = params.Nick
	e.Context = params.Context
}

func (e *Event) TmplFuncs() template.FuncMap {
	return template.FuncMap{
		"args": e.GetArg,
		"rand": randomSelect,
	}
}

func randomSelect(choices string) string {
	return *lib.RandSelect(strings.Split(choices, " | "))
}

func (e *Event) GetArg(index int) string {
	if len(e.Args) < index {
		return ""
	}
	return e.Args[index-1]
}

func ReplaceVars(aliasStr string) string {
	r, err := regexp.Compile("\\$\\{([a-z0-9]+)\\}")
	if err != nil {
		logger.Error(fmt.Sprintf("Couldn't compile VarReg regexp: %s", err.Error()))
		return aliasStr
	}
	for i := 0; i < 10; i++ {
		match := r.FindString(aliasStr)
		if match != "" {
			variable := Vars.GetOne(match[2 : len(match)-1])
			if variable == match {
				aliasStr = fmt.Sprintf("%s -> Error: %q variable refers to itself.", aliasStr, variable)
				return aliasStr
			}
			if variable != "" {
				aliasStr = strings.Replace(aliasStr, match, variable, -1)
			}
		} else {
			return aliasStr
		}
	}
	return aliasStr
}

func Register() {
	defer logger.Info(lib.TimeTrack(time.Now(), "Loading the Alias Backend plugin"))
	DB = *storage.NewStringListDB("alias.db")
	Vars = *storage.NewStringListDB("variables.db")
}
