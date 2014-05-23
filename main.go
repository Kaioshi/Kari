// Garionette project main.go
package main

import (
	"Kari/config"
	"Kari/irc"
	"Kari/plugins/core"
	"Kari/plugins/google"
)

func main() {
	conf := &config.Config{}
	conf.Parse()
	bot := &irc.IRC{*conf, nil}

	core.Register(bot)
	google.Register(bot)

	bot.Start()
}
