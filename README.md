# simple-audio
This is a simple audio player that can be remotely controlled via a webinterface.

This software is currently under development, but the webinterface already works mostly.

In the future this player should be remotely controlled via sACN (DMX).

You can find downloads on the [releases](https://github.com/Hundemeier/simple-audio/releases) page.

## Philosophy

The player consists of two elements:
* pool
* slots

For more information about the internal workings of these components, 
read the [Internal](https://github.com/Hundemeier/simple-audio#internal) section.

You can upload audio files (*.mp3 and *.wav) to the pool for future use.
You can then assign a file from the pool to one slot. Currently there are the slots from 1-512 available.
The settings, like volume or looping, are stored to the slots and will be reseted when 
changing to a new file for the slot.

## Internal

This player is in its core just an executable. So the webinterface is baked into this file. 

However, if you upload a file via the webinterface it will create a folder named 
"pool" directly besides this executable and stores the uploaded files there. 
So you should not have another folder named "pool" or "Pool" next to the executabel or the 
behaviour of this software is not specified.

In addition to the pool-folder, a config.json file is created as soon as something in 
the configuration of the slots has changed. This means that this file is constantly updated 
and should not be removed or manually changed.

The program reads the configuration on startup. So for a clean start, simply remove the config.json file.
