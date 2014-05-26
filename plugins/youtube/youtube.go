package youtube

import (
	"Kari/irc"
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/web"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type YouTubeResults struct {
	Feed YouTubeFeed `json:"feed"`
}

type YouTubeFeed struct {
	Entry []YouTubeEntry `json:"entry"`
}

type YouTubeEntry struct {
	Title  map[string]string `json:"title"`
	Info   InfoEntry         `json:"media$group"`
	Rating RatingEntry       `json:"gd$rating"`
	Stats  map[string]string `json:"yt$statistics"`
}

type InfoEntry struct {
	Description map[string]string `json:"media$description"`
	Duration    map[string]string `json:"yt$duration"`
	Date        map[string]string `json:"yt$uploaded"`
	ID          map[string]string `json:"yt$videoid"`
}

type RatingEntry struct {
	Average float32 `json:"average"`
	Max     int     `json:"max"`
	Min     int     `json:"min"`
}

func Register(bot *irc.IRC) {
	logger.Info("Registering YouTube hooks")

	events.CmdListen(&events.CmdListener{
		Commands: []string{"youtube", "yt"},
		Help:     "YouTubes stuff.",
		Syntax:   bot.Config.Prefix + "yt <search terms> - Example: " + bot.Config.Prefix + "yt we like big booty bitches",
		Callback: func(input *events.Params) {
			ytr := &YouTubeResults{}
			uri := fmt.Sprintf("https://gdata.youtube.com/feeds/api/videos?q=%s&max-results=1&v=2&alt=json",
				url.QueryEscape(strings.Join(input.Args, " ")))
			web.Get(uri, func(ERROR string, body []byte) {
				err := json.Unmarshal(body, &ytr)
				if err != nil {
					logger.Error("Couldn't parse youtube's JSON:" + err.Error())
					return
				}
				yt := &ytr.Feed.Entry[0]
				duration, err := time.ParseDuration(yt.Info.Duration["seconds"] + "s")
				resp := fmt.Sprintf("%s ~ [%s] %s - %s views ~ http://youtu.be/%s ~ %s",
					yt.Title["$t"], duration, yt.Info.Date["$t"][0:10],
					lib.CommaNum(yt.Stats["viewCount"]), yt.Info.ID["$t"],
					yt.Info.Description["$t"])
				bot.Say(input.Context, resp)
			})
		}})
}
