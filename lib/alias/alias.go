package alias

import (
	"Kari/config"
	"Kari/irc/events"
	"Kari/irc/globals"
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/storage"
	"Kari/lib/words"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var DB storage.StringListDB
var Vars storage.StringListDB
var Config *config.Config

type Event struct {
	From        string
	Context     string
	Data        string
	WhippingBoy string
	RandNick    string
	Args        []string
}

func (e *Event) Populate(params *events.Params, args []string, aliasEntry *string) {
	e.Data = strings.Join(args, " ")
	e.Args = args
	e.From = params.Nick
	e.Context = params.Context
	if strings.Index(*aliasEntry, ".RandNick") > -1 {
		e.RandNick = globals.Channels[strings.ToLower(params.Context)].RandNick()
	}
	if strings.Index(*aliasEntry, ".WhippingBoy") > -1 {
		e.WhippingBoy = *lib.RandSelect(Config.WhippingBoys)
	}
}

func (e *Event) TmplFuncs() template.FuncMap {
	return template.FuncMap{
		"args":        e.GetArg,
		"rand":        randomSelect,
		"first":       firstNotEmpty,
		"verb":        words.RandVerb,
		"verbs":       words.RandVerbs,
		"verbed":      words.RandVerbed,
		"verbing":     words.RandVerbing,
		"noun":        words.RandNoun,
		"adjective":   words.RandAdjective,
		"adverb":      words.RandAdverb,
		"pronoun":     words.RandPronoun,
		"preposition": words.RandPreposition,
	}
}

func randomSelect(choices ...string) string {
	var choice string
	if len(choices) > 1 {
		choice = *lib.RandSelect(choices)
	} else {
		choice = choices[0]
	}
	if strings.Index(choice, " | ") > -1 {
		return *lib.RandSelect(strings.Split(choice, " | "))
	}
	return choice

}

func firstNotEmpty(args ...string) string {
	for _, arg := range args {
		if len(arg) > 0 {
			return arg
		}
	}
	return ""
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

func Register(conf *config.Config) {
	defer logger.Info(lib.TimeTrack(time.Now(), "Loading the Alias Backend plugin"))
	DB = *storage.NewStringListDB("alias.db")
	Vars = *storage.NewStringListDB("variables.db")
	Config = conf
}
