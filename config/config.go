package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Server, Port, Nick, Username, Realname, Prefix string
	Autojoin                                       []string
}

func (conf *Config) Parse() {
	file, err := os.Open("bot.conf")
	if err != nil {
		fmt.Println("Couldn't open config:", err.Error())
		os.Exit(1)
	}
	defer file.Close()
	fi, err := file.Stat()
	data := make([]byte, fi.Size())
	count, err := file.Read(data)
	if err != nil {
		fmt.Println("Couldn't read config:", err.Error())
		os.Exit(1)
	}
	var _ = count
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line != "" && line[0:1] != "#" {
			index := strings.Index(line, ": ")
			field := line[0:index]
			entry := line[index+2:]
			switch field {
			case "server":
				conf.Server = entry
			case "nickname":
				conf.Nick = entry[0:strings.Index(entry, ",")]
			case "username":
				conf.Username = entry
			case "realname":
				conf.Realname = entry
			case "port":
				conf.Port = entry
			case "command prefix":
				conf.Prefix = entry
			case "autojoin":
				conf.Autojoin = strings.Split(entry, ", ")
			default:
				// unparsed
				//fmt.Println("Unparsed config line:", line)
			}
		}
	}
}
