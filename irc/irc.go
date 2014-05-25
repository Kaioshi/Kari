package irc

import (
	"Kari/config"
	"Kari/events"
	"Kari/lib"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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
	fmt.Println(lib.Timestamp("-> " + line))
	fmt.Fprint(irc.Conn, lib.Sanitise(line)+"\r\n")
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
		params.Data = strings.Join(args[3:len(args)], " ")[1:]
		if params.Data[0:1] == irc.Config.Prefix && params.Data[1:2] != "" {
			params.Command = args[3][2:]
			params.Args = strings.Fields(params.Data[len(params.Command)+1:])
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
	//fmt.Println(params)
}

func (irc *IRC) handleData(raw []byte) {
	line := string(raw)
	if line[0:1] != ":" {
		if line[0:4] == "PING" {
			irc.SilentSend("PONG " + line[5:])
		} else {
			fmt.Println(lib.Timestamp("<- " + line))
		}
		sem <- 1
		return
	}
	fmt.Println(lib.Timestamp("<- " + line))
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
			fmt.Println("Connection error:", err.Error())
			os.Exit(4)
		}
		if prefix == true {
			sem <- 1
			continue
		}
		go irc.handleData(line)
	}
}
