package google

import (
	"Kari/irc"
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/web"
)

func Register(bot *irc.IRC) {
	logger.Info("Registering Google hooks")

	events.CmdListen(&events.CmdListener{
		Commands: []string{"google", "g"},
		Help:     "Googles stuff~",
		Syntax:   bot.Config.Prefix + "g <search terms>",
		Callback: func(input *events.Params) {
			g := &web.Google(input.Data, 1).Results.Data[0]
			bot.Say(input.Context, g.Title+" ~ "+g.URL+" ~ "+lib.StripHtml(g.Content))
		}})
}
