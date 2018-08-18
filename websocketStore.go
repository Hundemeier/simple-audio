package main

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type wsStore struct {
	list     []*websocket.Conn
	mtx      sync.Mutex
	writeMtx sync.Mutex
}

func newWsStore() *wsStore {
	return &wsStore{
		list: make([]*websocket.Conn, 0),
	}
}

//addConn adds the given Connection to the list if it exists not already
func (s *wsStore) addConn(conn *websocket.Conn) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.contains(conn) {
		return
	}
	s.list = append(s.list, conn)
	go s.readWebsocket(conn)
}

func (s *wsStore) contains(conn *websocket.Conn) bool {
	for _, ws := range s.list {
		if ws == conn {
			return true
		}
	}
	return false
}

func (s *wsStore) remove(conn *websocket.Conn) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	//get the index
	index := -1
	for i, ws := range s.list {
		if ws == conn {
			index = i
			break
		}
	}
	if index < 0 {
		//no valid index was found
		return
	}
	//remove
	if index == len(s.list)-1 {
		//if we are here it is the last element in the list
		s.list = s.list[:len(s.list)-1]
	} else {
		s.list = append(s.list[:index], s.list[index+1:]...)
	}
}

//getList returns a copy of the list with all websocket.Conns
func (s *wsStore) getList() []*websocket.Conn {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	copy := make([]*websocket.Conn, 0, len(s.list))
	for _, ws := range s.list {
		copy = append(copy, ws)
	}
	return copy
}

func (s *wsStore) writeToWebsockets(item slotItem) {
	s.writeMtx.Lock()
	defer s.writeMtx.Unlock()
	list := s.getList()
	for _, conn := range list {
		conn.WriteJSON(item)
	}
}

//readWebsocket is currently an internal function to the wsStore struct
func (s *wsStore) readWebsocket(conn *websocket.Conn) {
	for {
		//check if the store still contains this conn so we can leave this goroutine
		if !s.contains(conn) {
			break
		}
		//read
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if msgType != websocket.TextMessage {
			continue //skip if not a textmessgae
		}
		input := make(map[string]*json.RawMessage)
		err = json.Unmarshal(data, &input)
		if err != nil {
			//could not read message as json
			continue
		}
		go handleInputEvent(input)
	}
	conn.Close()
	s.remove(conn)
}
