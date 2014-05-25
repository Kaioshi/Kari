package core

import (
	"Kari/events"
	"Kari/irc"
	"fmt"
	"runtime"
	"strings"
)

func Register(bot *irc.IRC) {
	fmt.Println("Registering Core hooks")

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
			if len(input.Args) < 1 {
				if input.Context[0:1] == "#" {
					bot.Part(input.Context)
				} else {
					bot.Say(input.Context, "Either specify the channel to part or use the command in the channel.")
				}
			} else if input.Args[0][0:1] != "#" {
				bot.Say(input.Context, "That ain't how you part a channel sucka")
				return
			} else {
				bot.Part(input.Args[0])
			}
		}})

	// raw
	events.CmdListen(&events.CmdListener{
		Command: "raw",
		Help:    "Sends raw commands to the server",
		Syntax:  bot.Config.Prefix + "raw <command>",
		Callback: func(input *events.Params) {
			bot.Send(strings.Join(input.Args, " "))
		}})

	// memstats
	events.CmdListen(&events.CmdListener{
		Command: "memstats",
		Help:    "Shows mem stats~",
		Syntax:  bot.Config.Prefix + "memstats",
		Callback: func(input *events.Params) {
			m := &runtime.MemStats{}
			runtime.ReadMemStats(m)
			bot.Say(input.Context, fmt.Sprintf("Sys: %d KiB, Allocated and in use: %d KiB, Total Allocated (including freed): %d KiB, Lookups: %d, Mallocs: %d, Frees: %d",
				(m.Sys/1024.0), (m.Alloc/1024), (m.TotalAlloc/1024), m.Lookups, m.Mallocs, m.Frees))
		}})
}
