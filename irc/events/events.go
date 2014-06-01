// events
package events

import (
	"Kari/lib/logger"
	"fmt"
	"regexp"
)

var eventListeners map[string][]*EvListener = map[string][]*EvListener{}
var complexListeners map[string][]*ComplexEventListener = map[string][]*ComplexEventListener{}
var commandListeners map[string]*CmdListener = map[string]*CmdListener{}

func CommandList() []string {
	cmdList := make([]string, 0)
	for _, cmd := range commandListeners {
		if len(cmd.Commands) > 0 {
			for _, subcmd := range cmd.Commands {
				cmdList = append(cmdList, subcmd)
			}
		} else {
			cmdList = append(cmdList, cmd.Command)
		}
	}
	return cmdList
}

func Help(command string, helpType string) string {
	if _, ok := commandListeners[command]; ok {
		switch helpType {
		case "syntax":
			return fmt.Sprintf("[Help] %s", commandListeners[command].Syntax)
		case "help":
			return fmt.Sprintf("[Help] %s", commandListeners[command].Help)
		}
	}
	return fmt.Sprintf("No %s found for %q :\\", helpType, command)
}

type EvListener struct {
	Handle   string
	Event    string
	Callback func(input *Params)
}

type CmdListener struct {
	Command  string
	Commands []string
	Help     string
	Syntax   string
	Callback func(input *Params)
}

type ComplexEventListener struct {
	Handle   string
	Event    string
	Regex    string
	r        regexp.Regexp
	Callback func(input *Params)
}

type Params struct {
	Context, Nick, Address, Data, Command string
	Newnick, Kicknick, Message            string
	Match                                 string
	Args                                  []string
}

func (cpe *ComplexEventListener) String() string {
	return fmt.Sprintf("Handle: %s, Event: %s, Regex: %s", cpe.Handle, cpe.Event, cpe.Regex)
}

func (c *CmdListener) String() string {
	return fmt.Sprintf("Command: %s, Help: %s, Syntax: %s", c.Command, c.Help, c.Syntax)
}

func (e *EvListener) String() string {
	return fmt.Sprintf("Handle: %s, Event: %s", e.Handle, e.Event)
}

func (p *Params) String() string {
	args := ""
	if len(p.Args) > 0 {
		for _, value := range p.Args {
			args += value + ", "
		}
		args = args[0 : len(args)-2]
	}
	return fmt.Sprintf("Context: %s, Nick: %s, Address: %s, Data: %s, Command: %s, Args: %s, Newnick: %s, Kicknick: %s, Message: %s",
		p.Context, p.Nick, p.Address, p.Data, p.Command, args, p.Newnick, p.Kicknick, p.Message)
}

func EvListenComplex(event *ComplexEventListener) {
	r, err := regexp.Compile(event.Regex)
	if err != nil {
		logger.Error(fmt.Sprintf("EvListenComplex: Couldn't compile %s's Regex: %s",
			event.Handle, event.Regex))
		return
	}
	event.r = *r
	complexListeners[event.Event] = append(complexListeners[event.Event], event)
}

func CmdListen(command *CmdListener) {
	if len(command.Commands) > 0 {
		commandListeners[command.Commands[0]] = command
	} else {
		commandListeners[command.Command] = command
	}
}

func EvListen(event *EvListener) {
	eventListeners[event.Event] = append(eventListeners[event.Event], event)
}

func Emit(event string, input *Params) {
	if event == "PRIVMSG" || event == "NOTICE" {
		for command, _ := range commandListeners {
			if len(commandListeners[command].Commands) > 0 {
				for _, subcmd := range commandListeners[command].Commands {
					if subcmd == input.Command {
						go fireCommand(*commandListeners[command], input, command)
					}
				}
			} else if input.Command == command {
				go fireCommand(*commandListeners[command], input, command)
			}
		}
	}
	if events, ok := eventListeners[event]; ok {
		for _, event := range events {
			go fireEvent(*event, input)
		}
	}
	if events, ok := complexListeners[event]; ok {
		for _, event := range events {
			go fireComplexEvent(*event, input)
		}
	}
}

func fireCommand(c CmdListener, input *Params, command string) {
	defer catchPanic("command", command)
	c.Callback(input)
}

func fireEvent(e EvListener, input *Params) {
	defer catchPanic("event "+e.Event, e.Handle)
	e.Callback(input)
}

func fireComplexEvent(cpe ComplexEventListener, input *Params) {
	if match := cpe.r.FindString(input.Data); match != "" {
		defer catchPanic(fmt.Sprintf("complex event %s (%s)",
			cpe.Event, cpe.Regex), cpe.Handle)
		input.Match = match
		cpe.Callback(input)
	}
}

func catchPanic(listenType string, handle string) {
	if e := recover(); e != nil {
		logger.Error(fmt.Sprintf("Caught panic in %s \"%s\": %s", listenType, handle, e))
	}
}
