package urbandictionary

import (
	"Kari/irc"
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/web"
	"encoding/json"
	"fmt"
	"net/url"
)

type UDResponse struct {
	List []ListEntry `json:"list"`
}

type ListEntry struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

func Register(bot *irc.IRC) {
	logger.Info("Registering Urban Dictionary hooks")

	events.CmdListen(&events.CmdListener{
		Commands: []string{"urbandictionary", "ud"},
		Help:     "Looks up Urban Dictionary entries. NSFW",
		Syntax:   bot.Config.Prefix + "ud <term> - Example: " + bot.Config.Prefix + "ud scrobble",
		Callback: func(input *events.Params) {
			uri := fmt.Sprintf("http://api.urbandictionary.com/v0/define?term=%s", url.QueryEscape(input.Data))
			body, err := web.Get(&uri)
			if err != "" {
				bot.Say(input.Context, err)
				return
			}
			ud := &UDResponse{}
			jserr := json.Unmarshal(body, &ud)
			if jserr != nil {
				logger.Error("Couldn't parse UD's JSON: " + jserr.Error())
				return
			}
			var resp string = ""
			var max int = 3
			if len(ud.List) < max {
				max = len(ud.List)
			}
			for i := 0; i < max; i++ {
				resp += fmt.Sprintf("%d) %s, ", i+1, ud.List[i].Definition)
			}
			fmt.Println(fmt.Sprintf(":%s PRIVMSG #pyoshi :%s", bot.Info.User, ud.List[0].Word+" ~ "+lib.SingleSpace(resp[0:len(resp)-2])))
			bot.Say(input.Context, ud.List[0].Word+" ~ "+lib.SingleSpace(resp[0:len(resp)-2]))
		}})
}
