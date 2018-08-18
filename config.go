package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

const configFile = "./config.json"

var configMutex = sync.Mutex{}

type config struct {
	SlotMap []slotConfig
}

type slotConfig struct {
	Slot     uint
	Filename string
	Player   playerItem
}

func writeConfig() (err error) {
	configMutex.Lock()
	defer configMutex.Unlock()
	slotItems := slotMap.getSlotItems()
	slotMap := make([]slotConfig, 0, len(slotItems))
	for _, item := range slotItems {
		slotMap = append(slotMap, slotConfig{
			Filename: item.Item.Name + "." + item.Item.Extension,
			Slot:     item.Slot,
			Player:   item.Player,
		})
	}
	conf := config{
		SlotMap: slotMap,
	}
	//Write data to file:
	data, err := json.Marshal(conf)
	if err != nil {
		return
	}
	f, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	f.Truncate(0) //clear existing data in the file
	_, err = f.Write(data)
	return
}

func readConfig() (conf config) {
	configMutex.Lock()
	defer configMutex.Unlock()
	conf.SlotMap = make([]slotConfig, 0)
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}
	json.Unmarshal(file, &conf)
	return
}
