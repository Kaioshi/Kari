// KARI NO Goooooooooooo
package main

import (
	"Kari/config"
	"Kari/irc"
	"Kari/irc/core"
	"Kari/irc/ial"
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/timer"
	"Kari/plugins/google"
	"Kari/plugins/urbandictionary"
	"Kari/plugins/youtube"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	conf := &config.Config{}
	conf.Parse()
	info := &irc.Info{}
	bot := &irc.IRC{Config: *conf, Info: *info}

	// required
	core.Register(bot)
	ial.Register(bot)

	// optional - comment out if you don't want 'em
	google.Register(bot)
	youtube.Register(bot)
	urbandictionary.Register(bot)

	logger.Info(fmt.Sprintf("Took %s to register plugin hooks.", time.Since(start)*time.Microsecond))

	timer.AddEvent("Garbage Collect", 15, lib.GC)

	bot.Start()
}
