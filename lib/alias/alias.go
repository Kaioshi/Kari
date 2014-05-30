package alias

import (
	"Kari/lib"
	"Kari/lib/logger"
	"Kari/lib/storage"
	"time"
)

var DB storage.StringDB

func Register() {
	defer logger.Info(lib.TimeTrack(time.Now(), "Loading the Alias Backend plugin"))
	DB = *storage.NewStringDB("alias.db")
}
