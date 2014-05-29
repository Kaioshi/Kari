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

var manga Manga

type Manga struct {
	MangaFox    map[string]MangaEntry
	MangaStream map[string]MangaEntry
}

type MangaEntry struct {
	Manga    string
	Title    string
	Link     string
	Date     int64
	Desc     string
	Announce []string
}

func parseRSS(rss []byte, source string) (map[string]MangaEntry, error) {
	src := strings.Split(lib.Sanitise(string(rss[bytes.Index(rss, []byte("<item>")):])), "</item>")
	src = src[0 : len(src)-1]
	entries := map[string]MangaEntry{}
	var title, tmpDate string
	var link [2]string
	if source == "mangafox" {
		link = [2]string{"<feedburner:origLink>", "</feedburner:origLink>"}
	} else {
		link = [2]string{"<link>", "</link>"}
	}
	for _, line := range src {
		if line == "" {
			continue
		}
		title = line[strings.Index(line, "<title>")+7 : strings.Index(line, "</title>")]
		tmpDate = line[strings.Index(line, "<pubDate>")+9 : strings.Index(line, "</pubDate>")]
		date, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", tmpDate)
		if err != nil {
			logger.Error(err.Error())
			return nil, errors.New("parseRSS failed to parse time : " + tmpDate)
		}
		entries[strings.ToLower(title)] = MangaEntry{
			Title: title,
			Link:  line[strings.Index(line, link[0])+len(link[0]) : strings.Index(line, link[1])],
			Date:  date.Unix(),
			Desc:  line[strings.Index(line, "<description>")+13 : strings.Index(line, "</description>")],
		}
	}
	return entries, nil
}

func addManga(title string, context string, source string, bot *irc.IRC) string {
	ltitle := strings.ToLower(title)
	switch source {
	case "MangaFox":
		if _, ok := manga.MangaFox[ltitle]; ok {
			return fmt.Sprintf("I'm already watching for %q updates on %s.", title, source)
		}
		manga.MangaFox[ltitle] = MangaEntry{
			Title:    title,
			Announce: []string{context},
		}
	case "MangaStream":
		if _, ok := manga.MangaStream[ltitle]; ok {
			return fmt.Sprintf("I'm already watching for %q updates on %s.", title, source)
		}
		manga.MangaStream[ltitle] = MangaEntry{
			Title:    title,
			Announce: []string{context},
		}
	}
	saveWatched()
	go checkUpdates(bot, "mangafox", "")
	go checkUpdates(bot, "mangastream", "")
	return "Added!"
}

func removeManga(title string, source string) string {
	ltitle := strings.ToLower(title)
	switch source {
	case "MangaFox":
		if _, ok := manga.MangaFox[ltitle]; ok {
			delete(manga.MangaFox, ltitle)
			saveWatched()
			return "Removed."
		}
	case "MangaStream":
		if _, ok := manga.MangaStream[ltitle]; ok {
			delete(manga.MangaStream, ltitle)
			saveWatched()
			return "Removed."
		}
	}
	return fmt.Sprintf("%q isn't on the %s watch list.", title, source)
}

func loadWatched() {
	db, err := ioutil.ReadFile("manga.db")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = json.Unmarshal(db, &manga)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func saveWatched() {
	out, err := json.Marshal(manga)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	ioutil.WriteFile("manga.db", out, 0666)
	logger.Info("Saved Manga watch list.")
}

func checkUpdates(bot *irc.IRC, source string, context string) {
	logger.Info("Checking Manga Sources for updates...")
	var uri, message string
	var watched map[string]MangaEntry
	switch source {
	case "mangafox":
		uri = "http://feeds.feedburner.com/mangafox/latest_manga_chapters?format=xml"
		watched = manga.MangaFox
	case "mangastream":
		uri = "http://mangastream.com/rss"
		watched = manga.MangaStream
	}
	data, err := web.Get(&uri)
	if err != "" {
		logger.Error(err)
		return
	}
	var entries map[string]MangaEntry
	entries, perr := parseRSS(data, source)
	if perr != nil {
		logger.Error(perr.Error())
		return
	}
	keys := getKeys(source)
	updates := make([]irc.RatedMessage, 0)
	var method string
	for title, entry := range entries {
		for _, key := range keys {
			if strings.Index(title, key) > -1 {
				if entry.Date > watched[key].Date {
					// update found
					newEntry := MangaEntry{
						Manga:    entry.Title[:len(key)],
						Title:    entry.Title,
						Date:     entry.Date,
						Desc:     entry.Desc,
						Link:     entry.Link,
						Announce: watched[key].Announce,
					}
					delete(watched, key)
					watched[key] = newEntry
					switch source {
					case "mangafox": // too much junk in mangafaucks's description
						message = fmt.Sprintf("%s is out \\o/ ~ %s", entry.Title, entry.Link)
					case "mangastream":
						message = fmt.Sprintf("%s is out \\o/ ~ %s ~ %q",
							entry.Title, entry.Link, entry.Desc)
					}
					if context != "" && !lib.HasElementString(watched[key].Announce, context) {
						if context[0:1] == "#" {
							method = "say"
						} else {
							method = "notice"
						}
						updates = append(updates, irc.RatedMessage{
							Method:  method,
							Target:  context,
							Message: message,
						})
					}
					for _, target := range watched[key].Announce {
						if target[0:1] == "#" {
							method = "say"
						} else {
							method = "notice"
						}
						updates = append(updates, irc.RatedMessage{
							Method:  method,
							Target:  target,
							Message: message,
						})
					}
				}
			}
		}
	}
	if len(updates) > 0 {
		bot.Rated(&updates)
		saveWatched()
	} else if context != "" {
		bot.Say(context, "Nothing new. :\\")
	}
}

func getKeys(source string) []string {
	var keys []string
	var watched *map[string]MangaEntry
	switch source {
	case "mangafox":
		watched = &manga.MangaFox
	case "mangastream":
		watched = &manga.MangaStream
	}
	for key, _ := range *watched {
		keys = append(keys, key)
	}
	return keys
}

func Register(bot *irc.IRC) {
	defer logger.Info(lib.TimeTrack(time.Now(), "Loading the MangaStream plugin"))

	manga = Manga{map[string]MangaEntry{}, map[string]MangaEntry{}}
	loadWatched()
	timer.AddEvent("Checking Manga Sources", 900, func() {
		go checkUpdates(bot, "mangafox", "")
		go checkUpdates(bot, "mangastream", "")
	})

	events.CmdListen(&events.CmdListener{
		Commands: []string{"mangastream", "ms"},
		Help:     "Manages the MangaStream release watcher",
		Syntax: fmt.Sprintf("%sms <add/remove/list> <manga title> - Example: %sms add One Piece",
			bot.Config.Prefix, bot.Config.Prefix),
		Callback: func(input *events.Params) {
			if len(input.Args) == 0 {
				bot.Say(input.Context, events.Help("ms", "syntax"))
				return
			}
			switch strings.ToLower(input.Args[0]) {
			case "list":
				if len(manga.MangaStream) == 0 {
					bot.Say(input.Context, "I'm not watching for any MangaStream releases right now. :<")
					return
				}
				var titles string
				for _, entry := range manga.MangaStream {
					if entry.Manga == "" {
						titles += entry.Title + ", "
					} else {
						titles += entry.Manga + ", "
					}
				}
				bot.Say(input.Context, fmt.Sprintf("I'm currently watching for %s updates to %s.",
					"MangaStream", titles[:len(titles)-2]))
			case "add":
				if len(input.Args) < 2 {
					bot.Say(input.Context, events.Help("ms", "syntax"))
					return
				}
				bot.Say(input.Context, addManga(strings.Join(input.Args[1:], " "), input.Context, "MangaStream", bot))
			case "remove":
				if len(input.Args) < 2 {
					bot.Say(input.Context, events.Help("ms", "syntax"))
					return
				}
				bot.Say(input.Context, removeManga(strings.Join(input.Args[1:], " "), "MangaStream"))
			case "check":
				checkUpdates(bot, "mangastream", input.Context)
			}
		}})

	events.CmdListen(&events.CmdListener{ // not sure how to make this neat yet ^ v
		Commands: []string{"mangafox", "mf"},
		Help:     "Manages the MangaFox release watcher",
		Syntax: fmt.Sprintf("%smf <add/remove/list> <manga title> - Example: %smf add One Piece",
			bot.Config.Prefix, bot.Config.Prefix),
		Callback: func(input *events.Params) {
			if len(input.Args) == 0 {
				bot.Say(input.Context, events.Help("mf", "syntax"))
				return
			}
			switch strings.ToLower(input.Args[0]) {
			case "list":
				if len(manga.MangaFox) == 0 {
					bot.Say(input.Context, "I'm not watching for any MangaFox releases right now. :<")
					return
				}
				var titles string
				for _, entry := range manga.MangaFox {
					if entry.Manga == "" {
						titles += entry.Title + ", "
					} else {
						titles += entry.Manga + ", "
					}
				}
				bot.Say(input.Context, fmt.Sprintf("I'm currently watching for %s updates to %s.",
					"MangaFox", titles[:len(titles)-2]))
			case "add":
				if len(input.Args) < 2 {
					bot.Say(input.Context, events.Help("mf", "syntax"))
					return
				}
				bot.Say(input.Context, addManga(strings.Join(input.Args[1:], " "), input.Context, "MangaFox", bot))
			case "remove":
				if len(input.Args) < 2 {
					bot.Say(input.Context, events.Help("mf", "syntax"))
					return
				}
				bot.Say(input.Context, removeManga(strings.Join(input.Args[1:], " "), "MangaFox"))
			case "check":
				checkUpdates(bot, "mangafox", input.Context)
			}
		}})
}
