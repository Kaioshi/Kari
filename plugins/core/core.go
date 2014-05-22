package core

import (
	"Kari/events"
	"Kari/irc"
	"fmt"
	"strings"
)

func Register(bot *irc.IRC) {
	// autojoin
	events.EvListen(&events.EvListener{
		Handle: "autojoin",
		Event:  "376",
		Callback: func(input *events.Params) {
			bot.Join(strings.Join(bot.Config.Autojoin, ","))
		}})

	// say
	events.CmdListen(&events.CmdListener{
		Command: "say",
		Help:    "Says stuff!",
		Syntax:  bot.Config.Prefix + "say <thing you want to say>",
		Callback: func(input *events.Params) {
			bot.Say(input.Context, strings.Join(input.Args, " "))
		}})

	// action
	events.CmdListen(&events.CmdListener{
		Command: "action",
		Help:    "/me's stuff!",
		Syntax:  bot.Config.Prefix + "action <thing you want the bot to emote>",
		Callback: func(input *events.Params) {
			bot.Action(input.Context, strings.Join(input.Args, " "))
		}})

	// join
	events.CmdListen(&events.CmdListener{
		Command: "join",
		Help:    "Joins channels!",
		Syntax:  bot.Config.Prefix + "join #channel",
		Callback: func(input *events.Params) {
			fmt.Println(input)
			fmt.Printf("input.Args length: %d\n", len(input.Args))
			if len(input.Args) < 1 || input.Args[0][0:1] != "#" {
				bot.Say(input.Context, "That ain't how you join a channel sucka")
				return
			}
			bot.Join(input.Args[0])
		}})

	// part
	events.CmdListen(&events.CmdListener{
		Command: "part",
		Help:    "Parts channels!",
		Syntax:  bot.Config.Prefix + "part #channel",
		Callback: func(input *events.Params) {
			if len(input.Args) < 1 || input.Args[0][0:1] != "#" {
				bot.Say(input.Context, "That ain't how you part a channel sucka")
				return
			}
			bot.Part(input.Args[0])
		}})
}
