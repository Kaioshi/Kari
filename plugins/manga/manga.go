package manga

import (
	"Kari/irc"
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/timer"
	"Kari/lib/web"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

var watched map[string]MSEntry

type MSEntry struct {
	Manga    string
	Title    string
	Link     string
	Date     int64
	Desc     string
	Announce []string
}

func parseRSS(rss []byte) (map[string]MSEntry, error) {
	ms := strings.Split(string(rss[bytes.Index(rss, []byte("<item>")):len(rss)-20]), "</item>")
	entries := map[string]MSEntry{}
	var title, tmpDate string
	for i, _ := range ms {
		if ms[i] == "" {
			continue
		}
		title = ms[i][strings.Index(ms[i], "<title>")+7 : strings.Index(ms[i], "</title>")]
		tmpDate = ms[i][strings.Index(ms[i], "<pubDate>")+9 : strings.Index(ms[i], "</pubDate>")]
		date, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", tmpDate)
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.New("parseRSS failed to parse time :" + tmpDate)
		}
		entries[strings.ToLower(title)] = MSEntry{
			Title: title,
			Link:  ms[i][strings.Index(ms[i], "<link>")+6 : strings.Index(ms[i], "</link>")],
			Date:  date.Unix(),
			Desc:  ms[i][strings.Index(ms[i], "<description>")+13 : strings.Index(ms[i], "</description>")],
		}
	}
	return entries, nil
}

func addManga(title string, context string, bot *irc.IRC) string {
	if _, ok := watched[strings.ToLower(title)]; ok {
		return title + " is already on the watch list."
	}
	watched[strings.ToLower(title)] = MSEntry{
		Title:    title,
		Announce: []string{context},
	}
	saveWatched()
	go checkUpdates(bot, "")
	return "Added!"
}

func removeManga(title string) string {
	ltitle := strings.ToLower(title)
	if _, ok := watched[ltitle]; ok {
		delete(watched, ltitle)
		saveWatched()
		return "Removed."
	}
	return fmt.Sprintf("%q isn't on the MangaStream watch list.", title)
}

func loadWatched() {
	db, err := ioutil.ReadFile("manga.db")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//var entries map[string]MSEntry
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = json.Unmarshal(db, &watched)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func saveWatched() {
	out, err := json.Marshal(watched)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	ioutil.WriteFile("manga.db", out, 0666)
	logger.Info("Saved Manga watch list.")
}

func checkUpdates(bot *irc.IRC, context string) {
	logger.Info("Checking MangaStream...")
	uri := "http://mangastream.com/rss"
	data, err := web.Get(&uri)
	if err != "" {
		logger.Error(err)
		return
	}
	var entries map[string]MSEntry
	entries, perr := parseRSS(data)
	if perr != nil {
		logger.Error(perr.Error())
		return
	}
	keys := getKeys()
	updates := false
	for title, entry := range entries {
		for _, key := range keys {
			if strings.Index(title, key) > -1 {
				if entry.Date > watched[key].Date {
					// update found
					newEntry := MSEntry{
						Manga:    entry.Title[:len(key)],
						Title:    entry.Title,
						Date:     entry.Date,
						Desc:     entry.Desc,
						Link:     entry.Link,
						Announce: watched[key].Announce,
					}
					delete(watched, key)
					watched[key] = newEntry
					for _, target := range watched[key].Announce {
						bot.Say(target, fmt.Sprintf("%s is out \\o/ ~ %s ~ %s",
							entry.Title, entry.Link, entry.Desc))
					}
					updates = true
				}
			}
		}
	}
	if updates {
		saveWatched()
	} else if context != "" {
		bot.Say(context, "Nothing new. :\\")
	}
}

func getKeys() []string {
	var keys []string
	for key, _ := range watched {
		keys = append(keys, key)
	}
	return keys
}

func Register(bot *irc.IRC) {
	defer logger.Info(lib.TimeTrack(time.Now(), "Loading the MangaStream plugin"))

	watched = map[string]MSEntry{}
	loadWatched()
	timer.AddEvent("Checking MangaStream", 900, func() {
		checkUpdates(bot, "")
	})

	events.CmdListen(&events.CmdListener{
		Commands: []string{"mangastream", "ms"},
		Help:     "Manages the MangaStream release watcher",
		Syntax:   bot.Config.Prefix + "ms <add/remove/list> <manga title> - Example: " + bot.Config.Prefix + "ms add One Piece",
		Callback: func(input *events.Params) {
			if len(input.Args) == 0 {
				bot.Say(input.Context, events.Help("ms", "syntax"))
				return
			}
			switch strings.ToLower(input.Args[0]) {
			case "list":
				if len(watched) == 0 {
					bot.Say(input.Context, "I'm not watching for any MangaStream releases right now. :<")
					return
				}
				var titles string
				for _, entry := range watched {
					titles += entry.Manga + ", "
				}
				bot.Say(input.Context, fmt.Sprintf("I'm currently watching for updates to %s.", titles[:len(titles)-2]))
			case "add":
				if len(input.Args) < 2 {
					bot.Say(input.Context, events.Help("ms", "syntax"))
					return
				}
				bot.Say(input.Context, addManga(strings.Join(input.Args[1:], " "), input.Context, bot))
			case "remove":
				if len(input.Args) < 2 {
					bot.Say(input.Context, events.Help("ms", "syntax"))
					return
				}
				bot.Say(input.Context, removeManga(strings.Join(input.Args[1:], " ")))
			case "check":
				checkUpdates(bot, input.Context)
			}
		}})
}
