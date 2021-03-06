// KARI NO Goooooooooooo
package main

import (
	"Kari/config"
	"Kari/irc"
	"Kari/irc/core"
	"Kari/irc/globals"
	"Kari/irc/ial"
	"Kari/lib"
	"Kari/lib/alias"
	"Kari/lib/timer"
	"Kari/plugins/aliasfrontend"
	"Kari/plugins/google"
	"Kari/plugins/manga"
	"Kari/plugins/titlesnarfer"
	"Kari/plugins/urbandictionary"
	"Kari/plugins/youtube"
)

func main() {
	conf := &config.Config{}
	conf.Parse()
	info := &globals.Info{}
	bot := &irc.IRC{Config: *conf, Info: *info}

	// required
	core.Register(bot)
	ial.Register(bot)
	alias.Register(conf)
	aliasfrontend.Register(bot)

	// optional - comment out if you don't want 'em
	google.Register(bot)
	youtube.Register(bot)
	urbandictionary.Register(bot)
	manga.Register(bot)
	titlesnarfer.Register(bot)

	timer.AddEvent("Garbage Collect", 15, lib.GC)

	bot.Start()
}
