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
	Help     string
	Syntax   string
	Callback func(input *Params)
}

type Params struct {
	Context, Nick, Address, Data, Command string
	Args                                  []string
}

func CmdListen(command *CmdListener) {
	commandListeners = append(commandListeners, *command)
	fmt.Println("Added command listener", command.Command)
}

func EvListen(event *EvListener) {
	eventListeners = append(eventListeners, *event)
	fmt.Println("Added event listener", event.Handle)
}

func Emit(event string, input *Params) {
	if event == "PRIVMSG" {
		for _, listener := range commandListeners {
			if input.Command == listener.Command {
				listener.Callback(input)
			}
		}
	}
	for _, listener := range eventListeners {
		if listener.Event == event {
			listener.Callback(input)
		}
	}
}
