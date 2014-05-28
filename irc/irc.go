package irc

import (
	"Kari/config"
	"Kari/irc/events"
	"Kari/lib"
	"Kari/lib/logger"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var sem = make(chan int, 1)

type IRC struct {
	Info   Info
	Config config.Config
	Conn   net.Conn
}

type Info struct {
	Nick     string
	Address  string
	User     string
	Network  string
	Server   string
	Channels lib.StrList
}

func (i *Info) String() string {
	return fmt.Sprintf("Nick: %s, Address: %s, User: %s, Network: %s, Server: %s, Channels: %s",
		i.Nick, i.Address, i.User, i.Network, i.Server, strings.Join(i.Channels.List, ", "))
}

func (irc *IRC) SilentSend(line string) {
	fmt.Fprintf(irc.Conn, lib.Sanitise(line)+"\r\n")
}

func (irc *IRC) Send(line string) {
	line = lib.Sanitise(line)
	logger.Sent(line)
	fmt.Fprint(irc.Conn, line+"\r\n")
}

func (irc *IRC) Connect() bufio.Reader {
	conn, err := net.Dial("tcp", irc.Config.Server+":"+irc.Config.Port)
	irc.Conn = conn
	if err != nil {
		logger.Error("Connection error")
		os.Exit(1)
	}
	irc.Send(fmt.Sprintf("NICK %s", irc.Config.Nicknames[0]))
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

type RatedMessage struct {
	Method  string
	Target  string
	Message string
}

func (irc *IRC) Rated(rm *[]RatedMessage) {
	var n int = 0
	for _, message := range *rm {
		func(method string, target string, message string, delay int) {
			time.AfterFunc(time.Duration(delay)*time.Millisecond, func() {
				switch method {
				case "say":
					irc.Say(target, message)
				case "action":
					irc.Action(target, message)
				case "notice":
					irc.Notice(target, message)
				}
			})
		}(message.Method, message.Target, message.Message, n)
		n += 250
	}
}

func trimLine(magic int, target *string, user *string, line *string) {
	var max int = (magic - len(*target)) - len(*user)
	if len(*line) > max {
		newline := *line
		*line = newline[0:max-3] + " .."
	}
}

func (irc *IRC) Say(target string, line string) {
	trimLine(498, &target, &irc.Info.User, &line)
	irc.Send(fmt.Sprintf("PRIVMSG %s :%s", target, line))
}

func (irc *IRC) Action(target string, line string) {
	trimLine(489, &target, &irc.Info.User, &line)
	irc.Send(fmt.Sprintf("PRIVMSG %s :\001ACTION %s\001", target, line))
}

func (irc *IRC) Notice(target string, line string) {
	trimLine(499, &target, &irc.Info.User, &line)
	irc.Send(fmt.Sprintf("NOTICE %s :%s", target, line))
}

// misc
func (irc *IRC) findParams(params *events.Params, line string, args []string) {
	if strings.Index(args[0], "!") == -1 { // server message?
		params.Context = args[3]
		params.Nick = args[0][1:]
		params.Args = args[1:]
		params.Data = strings.Join(args[1:], " ")
		return
	}
	params.Context = args[2]
	params.Nick = args[0][1:strings.Index(args[0], "!")]
	params.Address = args[0][strings.Index(args[0], "!")+1:]
	switch args[1] {
	case "NOTICE":
		fallthrough
	case "PRIVMSG":
		if args[2][0:1] != "#" { // queries
			params.Context = params.Nick
		}
		if args[3][1:2] == irc.Config.Prefix && args[3][2:3] != "" {
			params.Data = strings.Join(args[4:len(args)], " ")
			params.Command = args[3][2:]
			params.Args = args[4:len(args)]
		} else {
			params.Data = strings.Join(args[3:len(args)], " ")[1:]
		}
	case "NICK":
		params.Newnick = args[2][1:]
	case "JOIN":
		if params.Context[0:1] == ":" {
			params.Context = params.Context[1:]
		}
	case "PART":
		if len(args) > 3 {
			params.Message = strings.Join(args[3:], " ")[1:]
		}
		if params.Context[0:1] == ":" { // Y U NO CONSISTENT
			params.Context = params.Context[1:]
		}
	case "QUIT":
		params.Message = strings.Join(args[2:], " ")[1:]
	case "KICK":
		params.Kicknick = args[3]
		params.Message = strings.Join(args[4:], " ")[1:]
	case "MODE": // can't think why this is needed for now, dump its mojo in message
		params.Message = strings.Join(args[3:], " ")
	case "TOPIC":
		params.Message = strings.Join(args[3:], " ")[1:]
	}
	if params.Args == nil {
		params.Args = args
	}
	//fmt.Println(params)
}

func (irc *IRC) handleData(raw []byte) {
	line := string(raw)
	if line[0:1] != ":" {
		if line[0:4] == "PING" {
			irc.SilentSend("PONG " + line[5:])
		}
		sem <- 1
		return
	}
	args := strings.Fields(line)
	logger.Filter(&args[1], &line)
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
			logger.Error("Connection error: " + err.Error())
			os.Exit(4)
		}
		if prefix == true {
			sem <- 1
			continue
		}
		go irc.handleData(line)
	}
}
