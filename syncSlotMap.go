package main

import (
	"sync"
)

type syncSlotMap struct {
	mtx            sync.Mutex
	m              map[uint]*player
	onPlayerChange func(slot uint, play *player) //this gets called whenever one of the players in the set is changed
}

func newSyncSlotMap() syncSlotMap {
	return syncSlotMap{
		m: make(map[uint]*player),
	}
}

//setFileForSlot creates a new player for the given slot and uses the given filename for populating it.
//If the file does not exist or another problem occurs, an error will be returned
func (s *syncSlotMap) setFileForSlot(slot uint, filename string) error {
	//check if file exists, othwerise nothing happens (nothing gets deleted)
	if err := checkFilename(filename); err != nil {
		return err
	}
	s.clearSlot(slot) //clear the slot.
	s.mtx.Lock()
	defer s.mtx.Unlock()
	tmpPlay, err := newPlayer(filename)
	if err != nil {
		return err
	}
	tmpPlay.onChange = func(play *player) {
		if s.onPlayerChange != nil {
			s.onPlayerChange(slot, tmpPlay)
		}
	}
	s.m[slot] = tmpPlay
	go writeConfig()
	return nil
}

//clearSlot clears the slot from a player. If no player exists, nothing happens.
func (s *syncSlotMap) clearSlot(slot uint) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.m[slot] != nil {
		s.m[slot].drain()
		delete(s.m, slot)
	}
}

//getPlayer returns the player for the given slot. Can be nil, if no player was set.
func (s *syncSlotMap) getPlayer(slot uint) (play *player) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	play = s.m[slot]
	return
}

//getSlotItems returns a list with all items on the syncSlotMap
func (s *syncSlotMap) getSlotItems() (items []slotItem) {
	s.mtx.Lock()
	items = make([]slotItem, 0, len(s.m))
	for slot := range s.m {
		s.mtx.Unlock()
		/*//getPoolItem is using implicit the mutex of this map too, so we temporarily unlock the mutex
		s.mtx.Unlock()
		poolItem, err := getPoolItem(player.filename)
		if err != nil {
			//if an error occured when reading the poolItem skip it
			continue
		}
		item := slotItem{
			Slot: slot,
			Item: poolItem,
			Player: playerItem{
				Playing: !player.isPaused(),
				Volume:  player.volume(),
				Loop:    player.loop,
			},
		}
		items = append(items, item)
		s.mtx.Lock()*/
		item, err := s.getSlotItem(slot)
		if err != nil {
			continue //when error, skip the slot
		}
		items = append(items, item)
		s.mtx.Lock()
	}
	s.mtx.Unlock()
	return
}

func (s *syncSlotMap) getSlotItem(slot uint) (item slotItem, err error) {
	player := s.getPlayer(slot)
	if player == nil {
		return
	}
	poolItem, err := getPoolItem(player.filename)
	if err != nil {
		return
	}
	item = slotItem{
		Slot: slot,
		Item: poolItem,
		Player: playerItem{
			Playing: !player.isPaused(),
			Volume:  player.volume(),
			Loop:    player.loop,
			Current: player.currentSample(),
			Length:  player.maxSample(),
			Rate:    player.sampleRate,
		},
	}
	return
}

//checkIfUsed returns true if the file is used in the syncSlotMap
func (s *syncSlotMap) checkIfUsed(filename string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for _, player := range s.m {
		if player.filename == filename {
			return true
		}
	}
	return false
}
