// KARI NO Goooooooooooo
package main

import (
	"Kari/config"
	"Kari/irc"
	"Kari/irc/core"
	"Kari/irc/ial"
	"Kari/lib"
	"Kari/lib/alias"
	"Kari/lib/timer"
	"Kari/plugins/aliasfrontend"
	"Kari/plugins/google"
	"Kari/plugins/manga"
	"Kari/plugins/urbandictionary"
	"Kari/plugins/youtube"
)

func main() {
	conf := &config.Config{}
	conf.Parse()
	info := &irc.Info{}
	bot := &irc.IRC{Config: *conf, Info: *info}

	// required
	core.Register(bot)
	ial.Register(bot)
	alias.Register()
	aliasfrontend.Register(bot)

	// optional - comment out if you don't want 'em
	google.Register(bot)
	youtube.Register(bot)
	urbandictionary.Register(bot)
	manga.Register(bot)

	timer.AddEvent("Garbage Collect", 60, lib.GC)

	bot.Start()
}
