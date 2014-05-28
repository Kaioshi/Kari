package google

import (
	"Kari/irc"
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/web"
	"fmt"
	"time"
)

func Register(bot *irc.IRC) {
	defer logger.Info(lib.TimeTrack(time.Now(), "Loading the Google plugin"))

	events.CmdListen(&events.CmdListener{
		Commands: []string{"google", "g"},
		Help:     "Googles stuff~",
		Syntax:   bot.Config.Prefix + "g <search terms>",
		Callback: func(input *events.Params) {
			g := web.Google(input.Data, 1)
			if g.Error != "" {
				bot.Say(input.Context, g.Error)
				return
			}
			bot.Say(input.Context, fmt.Sprintf("%s ~ %s ~ %s",
				g.Results.Data[0].Title, g.Results.Data[0].URL,
				lib.StripHtml(g.Results.Data[0].Content)))
		}})
}
