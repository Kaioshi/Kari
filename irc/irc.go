package irc

import (
	"Kari/config"
	"Kari/events"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var sem = make(chan int, 1)

type IRC struct {
	Config config.Config
	Conn   net.Conn
}

func (irc *IRC) Send(line string) {
	fmt.Println(time.Now().Format("["+time.Stamp+"]"), "->", line)
	fmt.Fprintf(irc.Conn, line+"\r\n")
}

func (irc *IRC) Connect() bufio.Reader {
	conn, err := net.Dial("tcp", irc.Config.Server+":"+irc.Config.Port)
	irc.Conn = conn
	if err != nil {
		fmt.Println("Connection error")
		os.Exit(1)
	}
	irc.Send(fmt.Sprintf("NICK %s", irc.Config.Nick))
	irc.Send(fmt.Sprintf("USER %s localhost * :%s", irc.Config.Username, irc.Config.Realname))

	out := bufio.NewReader(irc.Conn)
	return *out
}

// basic IRC commands
func (irc *IRC) Join(channel string) {
	if channel != "" && channel[0:1] == "#" {
		irc.Send("JOIN " + channel)
	}
}

func (irc *IRC) Part(channel string) {
	if channel != "" && channel[0:1] == "#" {
		irc.Send("PART " + channel)
	}
}

func (irc *IRC) Say(target string, line string) {
	irc.Send("PRIVMSG " + target + " :" + line)
}

func (irc *IRC) Action(target string, line string) {
	irc.Send("PRIVMSG " + target + " :\001ACTION " + line + "\001")
}

// misc
func (irc *IRC) findParams(params *events.Params, line string, args []string) {
	if args[1] == "PRIVMSG" {
		params.Nick = args[0][1:strings.Index(args[0], "!")]
		params.Address = args[0][strings.Index(args[0], "!")+1:]
		params.Context = args[2]
		params.Data = strings.Join(args[3:len(args)], " ")[1:]
		if params.Data[0:1] == irc.Config.Prefix && params.Data[1:2] != "" {
			params.Command = args[3][2:]
			params.Args = strings.Fields(params.Data[len(params.Command)+1:])
		}
	}
}

func (irc *IRC) handleData(raw []byte) {
	line := string(raw)
	fmt.Println(time.Now().Format("["+time.Stamp+"]"), "<-", line)
	if line[0:1] != ":" {
		if line[0:4] == "PING" {
			irc.Send("PONG " + line[5:])
		}
		line = ""
		sem <- 1
		return
	}
	args := strings.Fields(line)
	params := &events.Params{}
	irc.findParams(params, line, args)
	events.Emit(args[1], params)
	sem <- 1
}

func (irc *IRC) Start() {
	out := irc.Connect()
	sem <- 1
	for {
		<-sem
		line, prefix, err := out.ReadLine()
		if err != nil {
			fmt.Println("ReadLine err:", err.Error())
			os.Exit(4)
		}
		if prefix == true {
			sem <- 1
			continue
		}
		go irc.handleData(line)
	}
}
