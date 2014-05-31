package aliasfrontend

import (
	"Kari/irc"
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/alias"
	"Kari/lib/logger"
	"fmt"
	"strings"
	"time"
)

func Register(bot *irc.IRC) {
	defer logger.Info(lib.TimeTrack(time.Now(), "Loading the Alias Frontend plugin"))

	events.CmdListen(&events.CmdListener{
		Command: "alias",
		Help:    "Aliases things!",
		Syntax: fmt.Sprintf("%salias <add/remove/list> <command> <args> - Example: %salias add whip action whips {args*}'s buttocks!",
			bot.Config.Prefix, bot.Config.Prefix),
		Callback: func(input *events.Params) {
			var argLen int = len(input.Args)
			if argLen == 0 {
				bot.Say(input.Context, events.Help("alias", "syntax"))
				return
			}
			switch strings.ToLower(input.Args[0]) {
			case "list":
				aliases := alias.DB.GetKeys()
				if len(aliases) == 0 {
					bot.Say(input.Context, "There are no aliases defined.")
				} else {
					bot.Say(input.Context, fmt.Sprintf("Aliases: %s", strings.Join(aliases, ", ")))
				}
			case "info":
				if argLen < 2 {
					bot.Say(input.Context, events.Help("alias", "syntax"))
					return
				}
				alias := alias.DB.GetOne(input.Args[1])
				if alias == "" {
					bot.Say(input.Context, fmt.Sprintf("There is no %q alias.", input.Args[1]))
				} else {
					bot.Say(input.Context, fmt.Sprintf("Alias %q contains: %s", input.Args[1], alias))
				}
			case "add":
				if argLen < 3 {
					bot.Say(input.Context, events.Help("alias", "syntax"))
					return
				}
				alias.DB.SaveOne(input.Args[1], strings.Join(input.Args[2:], " "))
				bot.Say(input.Context, "Added!")
			case "remove":
				if argLen < 2 {
					bot.Say(input.Context, events.Help("alias", "syntax"))
					return
				}
				if ok := alias.DB.RemoveOne(input.Args[1]); ok {
					bot.Say(input.Context, "Removed!")
				} else {
					bot.Say(input.Context, fmt.Sprintf("There is no %q alias.", input.Args[1]))
				}
			default:
				bot.Say(input.Context, events.Help("alias", "syntax"))
			}
		}})
}
