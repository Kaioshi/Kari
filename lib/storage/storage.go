// basic JSON storage
package storage

import (
	"Kari/lib/logger"
	"encoding/json"
	"fmt"
	//	"io/ioutil"
	//"errors"
	"os"
	"strings"
)

type StringDB struct {
	Filename string
	Entry    map[string]string
}

func (sdb *StringDB) RemoveOne(handle string) bool {
	lhandle := strings.ToLower(handle)
	if _, ok := sdb.Entry[lhandle]; ok {
		delete(sdb.Entry, lhandle)
		sdb.Save()
		return true
	}
	return false
}

func (sdb *StringDB) SaveOne(handle string, data string) {
	lhandle := strings.ToLower(handle)
	if _, ok := sdb.Entry[lhandle]; ok {
		delete(sdb.Entry, lhandle)
	}
	sdb.Entry[lhandle] = data
	sdb.Save()
}

func (sdb *StringDB) GetOne(handle string) string {
	if entry, ok := sdb.Entry[strings.ToLower(handle)]; ok {
		return entry
	}
	return ""
}

func (sdb *StringDB) HasKey(handle string) bool {
	if _, ok := sdb.Entry[strings.ToLower(handle)]; ok {
		return true
	}
	return false
}

func (sdb *StringDB) GetKeys() []string {
	keys := make([]string, 0)
	if len(sdb.Entry) > 0 {
		for key, _ := range sdb.Entry {
			keys = append(keys, key)
		}
	}
	return keys
}

func (sdb *StringDB) Save() {
	var err error
	var count int
	var data []byte
	var fd *os.File
	data, err = json.Marshal(&sdb.Entry)
	if err != nil {
		logger.Error(fmt.Sprintf("[StringDB.Save()] Couldn't Marshal %s JSON -> %s", sdb.Filename, err.Error()))
		return
	}
	fd, err = os.OpenFile("db/"+sdb.Filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.Error(fmt.Sprintf("[StringDB.Save()] Couldn't open or create db/%s -> %s", sdb.Filename, err.Error()))
		return
	}
	defer fd.Close()
	count, err = fd.Write(data)
	if err != nil {
		logger.Error(fmt.Sprintf("[StringDB.Save()] Couldn't save to db/%s -> %s", sdb.Filename, err.Error()))
		return
	}
	logger.Info(fmt.Sprintf("Wrote %d bytes to db/%s", count, sdb.Filename))
}

func NewStringDB(filename string) *StringDB {
	var err error
	var count int
	var data []byte
	var fi os.FileInfo
	var fd *os.File
	nDB := new(StringDB)
	nDB.Filename = filename
	nDB.Entry = map[string]string{}
	fd, err = os.OpenFile("db/"+nDB.Filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.Error(fmt.Sprintf("[storage.Open()] Couldn't open or create db/%s -> %s", nDB.Filename, err.Error()))
		return nDB
	}
	defer fd.Close()
	if fi, err = fd.Stat(); err != nil {
		logger.Error(fmt.Sprintf("[storage.Open()] Couldn't get file info on db/%s -> %s", nDB.Filename, err.Error()))
		return nDB
	}
	data = make([]byte, fi.Size())
	count, err = fd.Read(data)
	if count > 0 {
		if err = json.Unmarshal(data, &nDB.Entry); err != nil {
			logger.Error(fmt.Sprintf("[storage.Open()] Couldn't Unmarshal json in db/%s -> %s", nDB.Filename, err.Error()))
			return nDB
		}
	}
	return nDB
}
