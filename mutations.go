package main

import (
	"errors"

	"github.com/graphql-go/graphql"
)

var errNoSlot = errors.New("the slot number was not provided or the slot was not populated")

func mutateSetSlotMap(p graphql.ResolveParams) (interface{}, error) {
	slotInt, ok := p.Args["slot"].(int)
	if !ok {
		return nil, errNoSlot
	}
	filename, ok := p.Args["name"].(string)
	if !ok {
		return nil, errors.New("the name of the poolItem was not provided")
	}
	err := slotMap.setFileForSlot(uint(slotInt), filename)
	if err != nil {
		return false, err
	}
	return true, nil
}

func mutatePlayingState(p graphql.ResolveParams) (interface{}, error) {
	slotInt, ok := p.Args["slot"].(int)
	if !ok {
		return nil, errNoSlot
	}
	slot := uint(slotInt)
	player := slotMap.getPlayer(slot)
	if player == nil {
		return nil, errNoSlot
	}
	//reconstruct playerItem
	playerInput, ok := p.Args["player"].(map[string]interface{})
	if !ok {
		return nil, errors.New("no playerInput was provided")
	}
	//stop aka rewind
	rewind, ok := playerInput["stop"].(bool)
	if ok && rewind {
		player.stop()
	}
	//playing:
	playing, ok := playerInput["playing"].(bool)
	if ok {
		//playing provided
		if playing {
			player.resume()
		} else {
			player.pause()
		}
	}
	//loop
	loop, ok := playerInput["loop"].(bool)
	if ok {
		player.loop = loop
	}
	//volume
	volume, ok := playerInput["volume"].(float64)
	if ok {
		player.setVolume(volume)
	}
	go writeConfig()
	return playerItem{
		Playing: !player.isPaused(),
		Volume:  player.volume(),
		Loop:    player.loop,
	}, nil
}

func mutateDeleteFile(p graphql.ResolveParams) (interface{}, error) {
	file, ok := p.Args["file"].(string)
	if !ok {
		return false, errors.New("filename was not provided")
	}
	ok = deleteFile(file)
	if !ok {
		return false, errors.New("the file could not be deleted. It may be in use")
	}
	return true, nil
}

func mutateClearSlot(p graphql.ResolveParams) (interface{}, error) {
	slotInt, ok := p.Args["slot"].(int)
	if !ok {
		return nil, errNoSlot
	}
	slot := uint(slotInt)
	//check if the slot is populated
	if play := slotMap.getPlayer(slot); play == nil {
		return false, nil
	}
	slotMap.clearSlot(slot)
	return true, nil
}
