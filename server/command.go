package main

import(
	"encoding/json"
	"io/ioutil"
	"sync"
	"os"
)


var commandStore = &CommandStore{
	File: "commands.json",
}



type CommandStore struct {
	sync.Mutex
	File string
	Commands map[string][]string
}

func (cs *CommandStore) Load() error {
	cs.Lock()
	defer cs.Unlock()

	data, err := ioutil.ReadFile(cs.File)
	if err != nil {
		if os.IsNotExist(err) {
			cs.Commands = make(map[string][]string)
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &cs.Commands)
}


func (cs *CommandStore) Save() error {
	data, err := json.MarshalIndent(cs.Commands, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cs.File, data, 0644)
}


func (cs *CommandStore) AddCommand(clientID, command string) error {
	cs.Lock()
	defer cs.Unlock()
	cs.Commands[clientID] = append(cs.Commands[clientID], command)
	return cs.Save()
}


func (cs *CommandStore) GetCommands(clientID string) []string {
	cs.Lock()
	defer cs.Unlock()
	cmds := cs.Commands[clientID]
	delete(cs.Commands, clientID)
	_ = cs.Save()
	return cmds
}