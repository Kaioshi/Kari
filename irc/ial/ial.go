// Internal Address List
package ial

import (
	"Kari/events"
	"Kari/irc"
	"fmt"
	"strings"
)

type ChannelData struct {
	User map[string]*UserData
	//Active []string  - don't need this yet
}

type UserData struct {
	Nick, User, Address, Fulluser string
}

func (ud *UserData) String() string {
	return "\"" + ud.Fulluser + "\""
}

func (chd *ChannelData) String() string {
	//return fmt.Sprintf("%s -- Active: %s", chd.User, "["+strings.Join(chd.Active, ", ")+"]")
	return fmt.Sprintf("Users: %s", chd.User)
}

func Register(bot *irc.IRC) {
	fmt.Println("Registering Internal Address List hooks")

	Channels := make(map[string]*ChannelData)

	events.CmdListen(&events.CmdListener{
		Command: "ial",
		Help:    "Shows ial internal details",
		Syntax:  "TBD",
		Callback: func(input *events.Params) {
			fmt.Println(Channels)
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
				if bot.Info.Nick == "" {
					bot.Info.Nick = input.Nick
					bot.Info.User = input.Nick + "!" + input.Address
				}
				bot.Info.Channels.Add(input.Context)
				chdata := &ChannelData{User: make(map[string]*UserData)}
				Channels[strings.ToLower(input.Context)] = chdata
			}
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialPart",
		Event:  "PART",
		Callback: func(input *events.Params) {
			if input.Nick == bot.Info.Nick {
				delete(Channels, strings.ToLower(input.Context))
			} else {
				delete(Channels[strings.ToLower(input.Context)].User, strings.ToLower(input.Nick))
			}
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialWho",
		Event:  "352",
		Callback: func(input *events.Params) {
			Channels[strings.ToLower(input.Context)].User[strings.ToLower(input.Args[6])] = &UserData{
				Nick:     input.Args[6],
				User:     input.Args[3],
				Address:  input.Args[4],
				Fulluser: input.Args[6] + "!" + input.Args[3] + "@" + input.Args[4],
			}
		}})
}
