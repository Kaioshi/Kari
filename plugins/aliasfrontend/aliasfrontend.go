package aliasfrontend

import (
	"Kari/irc"
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/alias"
	"Kari/lib/logger"
	"fmt"
	"strings"
)

func Register(bot *irc.IRC) {
	defer logger.Info(lib.TimeTrack(lib.TimeNow(), "Loading the Alias Frontend plugin"))

	events.CmdListen(&events.CmdListener{
		Command: "var",
		Help:    "Allows you to set variables.",
		Syntax: fmt.Sprintf("%svar <add/remove/list/info> <varname> [<data>] - Example: %svar add butts are nice",
			bot.Config.Prefix, bot.Config.Prefix),
		Callback: func(input *events.Params) {
			var argLen int = len(input.Args)
			if argLen == 0 {
				bot.Say(input.Context, events.Help("var", "syntax"))
				return
			}
			switch strings.ToLower(input.Args[0]) {
			case "add":
				if argLen < 3 {
					bot.Say(input.Context, events.Help("var", "syntax"))
					return
				}
				alias.Vars.SaveOne(input.Args[1], strings.Join(input.Args[2:], " "))
				bot.Say(input.Context, "Added!")
			case "remove":
				if argLen < 2 {
					bot.Say(input.Context, events.Help("var", "syntax"))
					return
				}
				if ok := alias.Vars.RemoveOne(input.Args[1]); ok {
					bot.Say(input.Context, "Removed!")
				} else {
					bot.Say(input.Context, fmt.Sprintf("There is no %q variable.", input.Args[1]))
				}
			case "info":
				if argLen < 2 {
					bot.Say(input.Context, events.Help("var", "syntax"))
					return
				}
				varInfo := alias.Vars.GetOne(input.Args[1])
				if varInfo == "" {
					bot.Say(input.Context, fmt.Sprintf("There is no %q variable.", input.Args[1]))
				} else {
					bot.Say(input.Context, fmt.Sprintf("Variable %q contains: %s", input.Args[1], varInfo))
				}
			case "list":
				variables := alias.Vars.GetKeys()
				if len(variables) == 0 {
					bot.Say(input.Context, "There are no variables defined.")
				} else {
					bot.Say(input.Context, fmt.Sprintf("Variables: %s", strings.Join(variables, ", ")))
				}
			default:
				bot.Say(input.Context, events.Help("var", "syntax"))
			}
		}})

	events.CmdListen(&events.CmdListener{
		Command: "alias",
		Help:    "Aliases things!",
		Syntax: fmt.Sprintf("%salias <add/remove/list/info> <command> <args> - Example: %salias add whip action whips {args*}'s buttocks!",
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
				aliasInfo := alias.DB.GetOne(input.Args[1])
				if aliasInfo == "" {
					bot.Say(input.Context, fmt.Sprintf("There is no %q alias.", input.Args[1]))
				} else {
					bot.Say(input.Context, fmt.Sprintf("Alias %q contains: %s", input.Args[1], aliasInfo))
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
