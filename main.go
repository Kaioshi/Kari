// KARI NO Goooooooooooo
package main

import (
	"Kari/config"
	"Kari/irc"
	"Kari/irc/core"
	"Kari/irc/ial"
	"Kari/lib/logger"
	"Kari/plugins/google"
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

	logger.Info(fmt.Sprintf("Took %s to register plugin hooks.", time.Since(start)*time.Microsecond))
	bot.Start()
}
