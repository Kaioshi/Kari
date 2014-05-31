package alias

import (
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/storage"
	"math/rand"
	"strings"
	"text/template"
	"time"
)

var DB storage.StringListDB

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
		"rand": e.RandomSelect,
	}
}

func (e *Event) RandomSelect(choices string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	selectFrom := strings.Split(choices, " | ")
	return selectFrom[r.Intn(len(selectFrom))]
}

func (e *Event) GetArg(index int) string {
	if len(e.Args) < index {
		return ""
	}
	return e.Args[index-1]
}

func Register() {
	defer logger.Info(lib.TimeTrack(time.Now(), "Loading the Alias Backend plugin"))
	DB = *storage.NewStringListDB("alias.db")
}
