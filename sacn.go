package main

import "github.com/Hundemeier/go-sacn/sacn"

//sACNhandler handles incoming sACN data.
func sACNhandler(oldData, newData sacn.DataPacket) {
	if universe != nil && uint(newData.Universe()) == *universe {
		//slotnumber is index + 1:
		//inspect all channels and set them accordingly
		for i, val := range newData.Data() {
			slot := uint(i + 1)
			setSlotDMX(slot, val)
		}
	}
}

//this sets a channel after the DMX value passed
//Channel usage:
//0-85: stop
//86-171: pause
//172-255: play
func setSlotDMX(slot uint, value byte) {
	play := slotMap.getPlayer(slot)
	if play != nil {
		switch {
		case value <= 85:
			play.stop()
		case value > 85 && value <= 171:
			play.pause()
		case value > 171:
			play.resume()
		}
	}
}
