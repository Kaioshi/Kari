package logger

import (
	"Kari/lib"
	"fmt"
	"strings"
)

func log(logtype int, line *string) {
	switch logtype {
	case 0:
		*line = "\u001b[90m[\u001b[94mInfo\u001b[90m]\u001b[0m " + *line
	case 1:
		*line = "\u001b[90m[\u001b[93mWarn\u001b[90m]\u001b[0m " + *line
	case 2:
		*line = "\u001b[90m[\u001b[93mDEBUG\u001b[90m]\u001b[0m " + *line
	case 3:
		*line = "\u001b[90m[\u001b[0mServ\u001b[90m]\u001b[0m " + *line
	case 4:
		*line = "\u001b[90m[\u001b[36mChat\u001b[90m]\u001b[0m " + *line
	case 5:
		*line = "\u001b[90m[\u001b[37mTraf\u001b[90m]\u001b[0m " + *line
	case 6:
		*line = "\u001b[90m[\u001b[91mERROR\u001b[90m]\u001b[0m " + *line
	case 7:
		*line = "\u001b[90m[\u001b[32mSent\u001b[90m]\u001b[0m " + *line
	}
	fmt.Println(lib.Timestamp(*line))
}

func Filter(TYPE *string, line *string) {
	switch *TYPE {
	case "NICK":
		fallthrough
	case "KICK":
		fallthrough
	case "MODE":
		fallthrough
	case "JOIN":
		fallthrough
	case "PART":
		fallthrough
	case "QUIT":
		fallthrough
	case "TOPIC":
		log(5, line)
	case "PRIVMSG":
		log(4, line)
	case "NOTICE":
		ln := *line
		if strings.Index(ln[0:strings.Index(ln, " ")], "@") > -1 {
			log(4, line)
		} else {
			log(3, line)
		}
	default:
		log(3, line)
	}
}

func Sent(line string) {
	log(7, &line)
}

func Error(line string) {
	log(6, &line)
}

func Traf(line string) {
	log(5, &line)
}

func Chat(line string) {
	log(4, &line)
}

func Serv(line string) {
	log(3, &line)
}

func Debug(line string) {
	log(2, &line)
}

func Info(line string) {
	log(0, &line)
}

func Warn(line string) {
	log(1, &line)
}

func Warning(line string) {
	log(1, &line)
}
