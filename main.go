// KARI NO Goooooooooooo
package main

import (
	"Kari/config"
	"Kari/irc"
	"Kari/irc/core"
	"Kari/irc/ial"
	"Kari/plugins/google"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	conf := &config.Config{}
	conf.Parse()
	info := &irc.Info{}
	bot := &irc.IRC{Config: *conf, Conn: nil, Info: *info}

	core.Register(bot)
	ial.Register(bot)
	google.Register(bot)

	fmt.Println("Took", time.Since(start)*time.Microsecond, "to register plugin hooks.")
	bot.Start()
}
