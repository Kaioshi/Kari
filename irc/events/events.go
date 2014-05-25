// events
package events

import (
	"Kari/lib/logger"
	"fmt"
)

var eventListeners = make([]EvListener, 0)
var commandListeners = make([]CmdListener, 0)

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

type Params struct {
	Context, Nick, Address, Data, Command string
	Newnick, Kicknick, Message            string
	Args                                  []string
}

func (c CmdListener) String() string {
	return fmt.Sprintf("Command: %s, Help: %s, Syntax: %s", c.Command, c.Help, c.Syntax)
}

func (e EvListener) String() string {
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

func CmdListen(command *CmdListener) {
	commandListeners = append(commandListeners, *command)
}

func EvListen(event *EvListener) {
	eventListeners = append(eventListeners, *event)
}

func Emit(event string, input *Params) {
	if event == "PRIVMSG" || event == "NOTICE" { // this seems wasteful somehow. TODO: make this efficient when you know how~
		for _, listener := range commandListeners {
			if len(listener.Commands) > 0 {
				for _, command := range listener.Commands {
					if command == input.Command {
						go fireCommand(listener, input, command)
					}
				}
			} else if input.Command == listener.Command {
				go fireCommand(listener, input, listener.Command)
			}
		}
	}
	for _, listener := range eventListeners {
		if listener.Event == event {
			go fireEvent(listener, input)
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

func catchPanic(listenType string, handle string) {
	if e := recover(); e != nil {
		logger.Error(fmt.Sprintf("== Error in %s \"%s\": %s", listenType, handle, e))
	}
}
