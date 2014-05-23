// events
package events

import "fmt"

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
		//args = args[0:-2]
	}
	return fmt.Sprintf("Context: %s, Nick: %s, Address: %s, Data: %s, Command: %s, Args: %s", p.Context, p.Nick, p.Address, p.Data, p.Command, args)
}

func CmdListen(command *CmdListener) {
	commandListeners = append(commandListeners, *command)
}

func EvListen(event *EvListener) {
	eventListeners = append(eventListeners, *event)
}

func Emit(event string, input *Params) {
	if event == "PRIVMSG" {
		for _, listener := range commandListeners {
			if len(listener.Commands) > 0 {
				for _, command := range listener.Commands {
					if command == input.Command {
						listener.Callback(input)
					}
				}
			} else if input.Command == listener.Command {
				//fmt.Printf("Listener: %s\n", listener)
				listener.Callback(input)
			}
		}
	}
	for _, listener := range eventListeners {
		if listener.Event == event {
			//fmt.Printf("Listener: %s\n", listener)
			listener.Callback(input)
		}
	}
}
