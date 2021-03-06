// shows url titles
package titlesnarfer

import (
	"Kari/irc"
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/web"
)

func Register(bot *irc.IRC) {
	defer logger.Info(lib.TimeTrack(lib.TimeNow(), "Loading the TitleSnarfer plugin"))

	events.EvListenComplex(&events.ComplexEventListener{
		Handle: "titlesnarfer",
		Event:  "PRIVMSG",
		Regex:  "(https?:\\/\\/[^\001 ]+)",
		Callback: func(input *events.Params) {
			if input.Command == "" {
				bot.Say(input.Context, web.GetTitle(input.Match))
			}
		}})
}
