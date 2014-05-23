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
			var resp web.GoogleResult
			resp = web.Google(strings.Join(input.Args, " "), 1)
			bot.Say(input.Context, resp.Results.Data[0].Title+" ~ "+resp.Results.Data[0].URL+" ~ "+lib.StripHtml(resp.Results.Data[0].Content))
		}})
}
