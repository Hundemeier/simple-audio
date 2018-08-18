package main

import (
	"sync"

	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"

	"github.com/faiface/beep"
)

type player struct {
	filename string //the filename of the corresponding file in the poolDir
	ctrl     *beep.Ctrl
	vol      *effects.Volume
	mtx      sync.Mutex
	//stopCallback func(*player) //this callback gets invoked, if the file is at the end, or stop() was used
	loop     bool
	seaker   beep.StreamSeekCloser
	onChange func(*player) //change callback with the new values as parameter. WITHOUT loop currently
}

func newPlayer(filename string) (play *player, err error) {
	play = &player{filename: filename}
	sampled, s, err := newPlayerInternal(filename)
	if err != nil {
		return
	}
	play.seaker = s
	//create a ctrl so the player can pause and resume the playback. Start in paused mode
	play.ctrl = &beep.Ctrl{Streamer: sampled, Paused: true}
	//create volume control:
	play.vol = &effects.Volume{
		Streamer: play.ctrl,
		Base:     2,
	}
	//create the "EndStreamer" that is used for the callback at the end of the Stream of the file
	playInternal(play)
	return
}

//newPlayerInternal does cresate a stremaer with resampler from the filename.
//Does not add this to speaker.Play!
func newPlayerInternal(filename string) (stream beep.Streamer, s beep.StreamSeekCloser, err error) {
	s, format, err := decode(filename)
	if err != nil {
		return
	}

	//use Resampler for the speaker and the given file, so we always use the correct sampleRate
	stream = beep.Resample(3, format.SampleRate, sampleRate, s)
	return
}

//playInternal invokes the speaker.Play method with the player
func playInternal(play *player) {
	speaker.Play(beep.Seq(play.vol, beep.Callback(func() {
		go func() {
			play.rewindInternal()
			playInternal(play) // because when we are in this callback, play is not in speaker.Play anymore
			if play.loop {
				play.ctrl.Paused = false
			} else {
				play.ctrl.Paused = true
			}
		}()
	})))
}

//callOnChnage is a small helper function for internal use!
func (p *player) callOnChange() {
	if p.onChange != nil {
		p.onChange(p)
	}
}

//pause pauses the running playback. If its already paused, nothing happens
func (p *player) pause() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.ctrl.Paused = true
	go p.callOnChange()
}

//resume starts playback of a paused player. If the player is already running, nothing happens.
//Note: this does not start a player that has reached its end again. Use stop() instead.
func (p *player) resume() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.ctrl.Paused = false
	go p.callOnChange()
}

//stop pauses a player and resets to the start of the file.
//Note: this is non-blocking, so it returns when the player is not ready to play.
//A call to player.resume() directly after this won't work.
func (p *player) stop() {
	//because we are calling internal functions we need no mutex lock
	p.pause()
	p.rewindInternal()
	go p.callOnChange()
}

//rewindInternal is an internal function.
func (p *player) rewindInternal() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.seaker.Seek(0)
	go p.callOnChange()
}

//drain removes the player from speaker.Play and closes the file
func (p *player) drain() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.ctrl.Streamer = nil //implicit remove from speaker.Play
	p.seaker.Close()
}

func (p *player) isPaused() bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	return p.ctrl.Paused
}

//setVloume sets the volume of the player. Values greater than 0 will increase
//and values less than 0 will decrease. Range can be: [-7;3]
func (p *player) setVolume(volume float64) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.vol.Volume = volume
	go p.callOnChange()
}

func (p *player) volume() float64 {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	return p.vol.Volume
}
