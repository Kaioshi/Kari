// basic JSON storage
package storage

import (
	"Kari/lib/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// StringListDB
type StringListDB struct {
	Filename string
	Entry    map[string]string
}

func (sl *StringListDB) GetOne(handle string) string {
	lhandle := strings.ToLower(handle)
	if entry, ok := sl.Entry[lhandle]; ok {
		return entry
	}
	return ""
}

func (sl *StringListDB) RemoveOne(handle string) bool {
	lhandle := strings.ToLower(handle)
	if _, ok := sl.Entry[lhandle]; ok {
		delete(sl.Entry, lhandle)
		sl.Save()
		return true
	}
	return false
}

func (sl *StringListDB) SaveOne(handle string, data string) {
	lhandle := strings.ToLower(handle)
	if _, ok := sl.Entry[lhandle]; ok {
		delete(sl.Entry, lhandle)
	}
	sl.Entry[lhandle] = data
	sl.Save()
}

func (sl *StringListDB) HasKey(handle string) bool {
	if _, ok := sl.Entry[strings.ToLower(handle)]; ok {
		return true
	}
	return false
}

func (sl *StringListDB) GetKeys() []string {
	keys := make([]string, 0)
	if len(sl.Entry) > 0 {
		for key, _ := range sl.Entry {
			keys = append(keys, key)
		}
	}
	return keys
}

func (sl *StringListDB) Save() {
	var err error
	var entries string
	if len(sl.Entry) > 0 {
		for key, entry := range sl.Entry {
			entries += fmt.Sprintf("%s %s\n", key, entry)
		}
		entries = entries[:len(entries)-1] // no trailing \n
		err = ioutil.WriteFile("db/"+sl.Filename, []byte(entries), 0666)
	} else {
		err = ioutil.WriteFile("db/"+sl.Filename, []byte{}, 0666)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("[StringListDB.Save()] Couldn't save to db/%s -> %s", sl.Filename, err.Error()))
		return
	}
	logger.Info(fmt.Sprintf("Saved db/%s ..", sl.Filename))
}

func NewStringListDB(filename string) *StringListDB {
	var err error
	var index int
	var data []byte
	var entries []string
	nDB := new(StringListDB)
	nDB.Filename = filename
	nDB.Entry = map[string]string{}
	data, err = ioutil.ReadFile("db/" + nDB.Filename)
	if err != nil {
		logger.Error(fmt.Sprintf("[storage.NewStringListDB()] Couldn't read db/%s -> %s", nDB.Filename, err.Error()))
		return nDB
	}
	if len(data) > 0 {
		entries = strings.Split(string(data), "\n")
		for _, line := range entries {
			if line == "" {
				continue
			}
			index = strings.Index(line, " ")
			nDB.Entry[line[:index]] = line[index+1:]
		}
	}
	return nDB
}

// StringDB
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
	var data []byte
	data, err = json.Marshal(&sdb.Entry)
	if err != nil {
		logger.Error(fmt.Sprintf("[StringDB.Save()] Couldn't Marshal %s JSON -> %s", sdb.Filename, err.Error()))
		return
	}
	if len(data) > 0 {
		err = ioutil.WriteFile("db/"+sdb.Filename, data, 0666)
	} else {
		err = ioutil.WriteFile("db/"+sdb.Filename, []byte{}, 0666)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("[StringDB.Save()] Couldn't save db/%s -> %s", sdb.Filename, err.Error()))
		return
	}
	logger.Info(fmt.Sprintf("Saved db/%s ..", sdb.Filename))
}

func NewStringDB(filename string) *StringDB {
	var err error
	var data []byte
	nDB := new(StringDB)
	nDB.Filename = filename
	nDB.Entry = map[string]string{}
	data, err = ioutil.ReadFile("db/" + nDB.Filename)
	if err != nil {
		logger.Error(fmt.Sprintf("[storage.NewStringDB()] Couldn't read db/%s -> %s", nDB.Filename, err.Error()))
		return nDB
	}
	if len(data) > 0 {
		if err = json.Unmarshal(data, &nDB.Entry); err != nil {
			logger.Error(fmt.Sprintf("[storage.Open()] Couldn't Unmarshal json in db/%s -> %s", nDB.Filename, err.Error()))
			return nDB
		}
	}
	return nDB
}
