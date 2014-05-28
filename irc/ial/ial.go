// Internal Address List
package ial

import (
	"Kari/irc"
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/logger"
	"fmt"
	"strings"
	"time"
)

var Channels map[string]*ChannelData = make(map[string]*ChannelData)

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
	return fmt.Sprintf("Users: %s", chd.User)
}

func Ison(channel string, nick string) bool {
	lchan := strings.ToLower(channel)
	if _, ok := Channels[lchan]; !ok {
		return false
	}
	lnick := strings.ToLower(nick)
	for key, _ := range Channels[lchan].User {
		if key == lnick {
			return true
		}
	}
	return false
}

func Register(bot *irc.IRC) {
	defer logger.Info(lib.TimeTrack(time.Now(), "Loading the Internal Address List plugin"))

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
				if bot.Info.User == "" {
					bot.Info.User = input.Nick + "!" + input.Address
				}
				bot.Info.Channels.Add(input.Context)
				chdata := &ChannelData{User: make(map[string]*UserData)}
				Channels[strings.ToLower(input.Context)] = chdata
			} else {
				Channels[strings.ToLower(input.Context)].User[strings.ToLower(input.Nick)] = &UserData{
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
				delete(Channels, strings.ToLower(input.Context))
			} else {
				delete(Channels[strings.ToLower(input.Context)].User, strings.ToLower(input.Nick))
			}
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialKick",
		Event:  "KICK",
		Callback: func(input *events.Params) {
			for channel, _ := range Channels {
				for key, user := range Channels[channel].User {
					if user.Nick == input.Kicknick {
						delete(Channels[channel].User, key)
					}
				}
			}
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialQuit",
		Event:  "QUIT",
		Callback: func(input *events.Params) {
			for channel, _ := range Channels {
				for key, user := range Channels[channel].User {
					if user.Nick == input.Nick {
						delete(Channels[channel].User, key)
					}
				}
			}
		}})

	events.EvListen(&events.EvListener{
		Handle: "ialNick",
		Event:  "NICK",
		Callback: func(input *events.Params) {
			for channel, _ := range Channels {
				for key, user := range Channels[channel].User {
					if user.Nick == input.Nick {
						newuser := &UserData{ // assigned in this order because of pure case nick changes
							Nick:     input.Newnick, // ie nick -> NICK
							User:     user.User,
							Address:  user.Address,
							Fulluser: input.Newnick + "!" + user.User + "@" + user.Address,
						}
						delete(Channels[channel].User, key)
						Channels[channel].User[strings.ToLower(input.Newnick)] = newuser
					}
				}
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
