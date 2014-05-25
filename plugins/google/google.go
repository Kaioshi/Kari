package google

import (
	"Kari/events"
	"Kari/irc"
	"Kari/lib"
	"Kari/lib/web"
	"fmt"
	"strings"
)

func Register(bot *irc.IRC) {
	fmt.Println("Registering Google hooks")

	events.CmdListen(&events.CmdListener{
		Commands: []string{"google", "g"},
		Help:     "Googles stuff~",
		Syntax:   bot.Config.Prefix + "g <search terms>",
		Callback: func(input *events.Params) {
			g := &web.Google(strings.Join(input.Args, " "), 1).Results.Data[0]
			bot.Say(input.Context, g.Title+" ~ "+g.URL+" ~ "+lib.StripHtml(g.Content))
		}})
}
