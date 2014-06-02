package irc

import (
	"Kari/config"
	"Kari/irc/events"
	"Kari/irc/globals"
	"Kari/lib"
	"Kari/lib/alias"
	"Kari/lib/logger"
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"text/template"
	"time"
)

var sem = make(chan int, 1)

type IRC struct {
	Info   globals.Info
	Config config.Config
	Conn   net.Conn
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
	var magic int = (498 - len(target)) - len(irc.Info.User)
	var i, index int
	for i = 0; i < 3; i++ {
		if len(line) > magic {
			index = strings.LastIndex(line[0:magic], " ")
			irc.Send(fmt.Sprintf("PRIVMSG %s :%s", target, line[0:index]))
			line = strings.TrimLeft(line[index:], " ")
		} else {
			irc.Send(fmt.Sprintf("PRIVMSG %s :%s", target, line))
			break
		}
	}
}

func (irc *IRC) Action(target string, line string) {
	var magic int = (489 - len(target)) - len(irc.Info.User)
	var i, index int
	for i = 0; i < 3; i++ {
		if len(line) > magic {
			index = strings.LastIndex(line[0:magic], " ")
			irc.Send(fmt.Sprintf("PRIVMSG %s :\001ACTION %s\001", target, line[0:index]))
			line = strings.TrimLeft(line[index:], " ")
		} else {
			irc.Send(fmt.Sprintf("PRIVMSG %s :\001ACTION %s\001", target, line))
			break
		}
	}
}

func (irc *IRC) Notice(target string, line string) {
	var magic int = (499 - len(target)) - len(irc.Info.User)
	var i, index int
	for i = 0; i < 3; i++ {
		if len(line) > magic {
			index = strings.LastIndex(line[0:magic], " ")
			irc.Send(fmt.Sprintf("NOTICE %s :%s", target, line[0:index]))
			line = strings.TrimLeft(line[index:], " ")
		} else {
			irc.Send(fmt.Sprintf("NOTICE %s :%s", target, line))
			break
		}
	}
}

// misc
func (irc *IRC) findParams(params *events.Params, line string, args []string) {
	var command string
	var index int
	if strings.Index(args[0], "!") == -1 { // server message?
		params.Context = args[3]
		params.Nick = args[0][1:]
		params.Args = args[1:]
		params.Data = line
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
			command = args[3][2:]
			if alias.DB.HasKey(command) { // temporary, need alias/args parser
				aliasEntry := alias.DB.GetOne(command)
				index = strings.Index(aliasEntry, " ")
				params.Command = aliasEntry[:index]
				if strings.Index(aliasEntry, "${") > -1 {
					aliasEntry = alias.ReplaceVars(aliasEntry)
				}
				if strings.Index(aliasEntry, "{{") > -1 {
					var aEvent alias.Event
					aEvent.Populate(params, args[4:len(args)], &aliasEntry)
					var out bytes.Buffer
					t, err := template.New(command).Funcs(aEvent.TmplFuncs()).Parse(aliasEntry[index+1:])
					if err != nil {
						params.Data = fmt.Sprintf("Couldn't parse %q template: %s", aliasEntry, err.Error())
						logger.Error(params.Data)
						return
					}
					err = t.Execute(&out, aEvent)
					if err != nil {
						params.Data = fmt.Sprintf("Couldn't execute %q template: %s", aliasEntry, err.Error())
						logger.Error(params.Data)
						return
					}
					params.Data = out.String()
					params.Args = strings.Fields(params.Data)
				} else {
					params.Data = aliasEntry[index+1:]
					params.Args = strings.Fields(params.Data)
				}
			} else {
				params.Command = command
				params.Args = args[4:len(args)]
				params.Data = strings.Join(params.Args, " ")
			}
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
	if params.Data == "" {
		params.Data = line
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
