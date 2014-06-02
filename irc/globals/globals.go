// apparently this is bad. can't think of another way around cyclic imports right now :\
package globals

import (
	"Kari/lib"
	"fmt"
	"strings"
)

var Channels map[string]*ChannelData = map[string]*ChannelData{}

type ChannelData struct {
	User map[string]*UserData
}

func (chd *ChannelData) String() string {
	return fmt.Sprintf("Users: %s", chd.User)
}

func (chd *ChannelData) RandNick() string {
	nicks := make([]string, len(chd.User))
	i := 0
	for nick, _ := range chd.User {
		nicks[i] = chd.User[nick].Nick
		i++
	}
	return *lib.RandSelect(nicks)
}

type UserData struct {
	Nick, User, Address, Fulluser string
}

func (ud *UserData) String() string {
	return "\"" + ud.Fulluser + "\""
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
