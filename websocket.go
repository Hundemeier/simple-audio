package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { //always anser request no matter what the origin is
		return true //always return true, so that every request is allowed
	},
}
var websockets = newWsStore()

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	websockets.addConn(c)
}

func handleInputEvent(event map[string]*json.RawMessage) {
	var slotInt int
	var slot uint
	if err := json.Unmarshal(*event["slot"], &slotInt); err == nil {
		slot = uint(slotInt)
	} else {
		return
	}
	//get the player for the slot
	player := slotMap.getPlayer(slot)
	if player == nil {
		//player is nil and the slot number must be invalid
		return
	}
	//get the "player" object from the json object
	playEvent := make(map[string]*json.RawMessage)
	err := json.Unmarshal(*event["player"], &playEvent)
	if err != nil {
		//no player object provided
		return
	}
	//stop
	if stop := unmarshalBool(playEvent, "stop"); stop != nil && *stop {
		player.stop()
	}
	//playing
	if playing := unmarshalBool(playEvent, "playing"); playing != nil {
		if *playing {
			player.resume()
		} else {
			player.pause()
		}
	}
	//loop
	if loop := unmarshalBool(playEvent, "loop"); loop != nil {
		player.loop = *loop
	}
	//volume
	var volume float64
	if playEvent["volume"] != nil {
		if err = json.Unmarshal(*playEvent["volume"], &volume); err == nil {
			player.setVolume(volume)
		}
	}
}

func unmarshalBool(data map[string]*json.RawMessage, key string) *bool {
	result := false
	if data[key] == nil {
		return nil
	}
	err := json.Unmarshal(*data[key], &result)
	if err != nil {
		return nil
	}
	return &result
}
