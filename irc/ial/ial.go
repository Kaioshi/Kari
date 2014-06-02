// Internal Address List
package ial

import (
	"Kari/irc"
	"Kari/irc/events"
	"Kari/irc/globals"
	"Kari/lib"
	"Kari/lib/logger"
	"fmt"
	"strings"
)

func Register(bot *irc.IRC) {
	defer logger.Info(lib.TimeTrack(lib.TimeNow(), "Loading the Internal Address List plugin"))

	events.CmdListen(&events.CmdListener{
		Command: "ial",
		Help:    "Shows ial internal details",
		Syntax:  "TBD",
		Callback: func(input *events.Params) {
			fmt.Println(globals.Channels)
		}})

	events.EvListen(&events.EvListener{ // grab our nick on connect
		Handle: "ial001",
		Event:  "001",
		Callback: func(input *events.Params) {
			bot.Info.Nick = input.Args[1]
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialJoin",
		Event:  "JOIN",
		Callback: func(input *events.Params) {
			if input.Nick == bot.Info.Nick {
				bot.Send("WHO " + input.Context)
				bot.Send("MODE " + input.Context)
				if bot.Info.User == "" {
					bot.Info.User = input.Nick + "!" + input.Address
				}
				bot.Info.Channels.Add(input.Context)
				chdata := &globals.ChannelData{User: make(map[string]*globals.UserData)}
				globals.Channels[strings.ToLower(input.Context)] = chdata
			} else {
				globals.Channels[strings.ToLower(input.Context)].User[strings.ToLower(input.Nick)] = &globals.UserData{
					Nick:     input.Nick,
					User:     input.Address[:strings.Index(input.Address, "@")],
					Address:  input.Address[strings.Index(input.Address, "@")+1:],
					Fulluser: input.Args[0][1:],
				}
			}
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialPart",
		Event:  "PART",
		Callback: func(input *events.Params) {
			if input.Nick == bot.Info.Nick {
				delete(globals.Channels, strings.ToLower(input.Context))
			} else {
				delete(globals.Channels[strings.ToLower(input.Context)].User, strings.ToLower(input.Nick))
			}
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialKick",
		Event:  "KICK",
		Callback: func(input *events.Params) {
			for channel, _ := range globals.Channels {
				for key, user := range globals.Channels[channel].User {
					if user.Nick == input.Kicknick {
						delete(globals.Channels[channel].User, key)
					}
				}
			}
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialQuit",
		Event:  "QUIT",
		Callback: func(input *events.Params) {
			for channel, _ := range globals.Channels {
				for key, user := range globals.Channels[channel].User {
					if user.Nick == input.Nick {
						delete(globals.Channels[channel].User, key)
					}
				}
			}
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialNick",
		Event:  "NICK",
		Callback: func(input *events.Params) {
			for channel, _ := range globals.Channels {
				for key, user := range globals.Channels[channel].User {
					if user.Nick == input.Nick {
						newuser := &globals.UserData{ // assigned in this order because of pure case nick changes
							Nick:     input.Newnick, // ie nick -> NICK
							User:     user.User,
							Address:  user.Address,
							Fulluser: input.Newnick + "!" + user.User + "@" + user.Address,
						}
						delete(globals.Channels[channel].User, key)
						globals.Channels[channel].User[strings.ToLower(input.Newnick)] = newuser
					}
				}
			}
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialWho",
		Event:  "352",
		Callback: func(input *events.Params) {
			globals.Channels[strings.ToLower(input.Context)].User[strings.ToLower(input.Args[6])] = &globals.UserData{
				Nick:     input.Args[6],
				User:     input.Args[3],
				Address:  input.Args[4],
				Fulluser: input.Args[6] + "!" + input.Args[3] + "@" + input.Args[4],
			}
		}})
}
