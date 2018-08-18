package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"

	"github.com/faiface/beep"
)

//getPoolItem returns a poolItem with all the informations about the file in the poolDir.
//If the file could not found or opened, an error is returned
func getPoolItem(filename string) (poolItem, error) {
	ext := filepath.Ext(filename)
	return poolItem{
		Name:      filename[:len(filename)-len(ext)],
		Extension: ext[1:],
		//[1:], because we do not want the dot at the beginning of the string
		Used: slotMap.checkIfUsed(filename),
	}, nil
}

//checkFilename checks the given file, if it exists and can be opened
func checkFilename(filename string) error {
	if filename == "" {
		return errors.New("could not open empty filename string")
	}
	_, err := os.Open(poolDir + "/" + filename)
	if err != nil {
		return err
	}
	return nil
}

func deleteFile(filename string) bool {
	//check if the file is used somewhere
	if slotMap.checkIfUsed(filename) {
		return false
	}
	err := os.Remove(poolDir + "/" + filename)
	if err != nil {
		return false
	}
	return true
}

//errDecode happens if a file can not be read as audio file,
//because the fileformat is not supported, or the file is not an audio file.
//A new player is created with pause = true. This also adds itself to the speaker. (speaker.Play(this))
var errDecode = errors.New("could not decode the given file")

//decode tries to decode the given file. If it could be decoded, the return values are set and ok is true.
//otherwise ok is false
func decode(filename string) (s beep.StreamSeekCloser, format beep.Format, err error) {
	err = checkFilename(filename)
	if err != nil {
		return
	}

	//currently this is the bruteforce method. Open the file and try to decode.
	//could be better solved by reading the filename extension?!
	//The benefit of the bruteforce is: the extension does not matter

	//Try to decode via WAV
	f, err := os.Open(poolDir + "/" + filename)
	if err != nil {
		return
	}
	s, format, err = wav.Decode(f)
	if err == nil {
		return
	}
	//Try to decode via MP3
	f, err = os.Open(poolDir + "/" + filename)
	if err != nil {
		return
	}
	s, format, err = mp3.Decode(f)
	if err == nil {
		return
	}
	err = errDecode
	return
}

func getMyInterfaceAddr() ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	addresses := []net.IP{}
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			addresses = append(addresses, ip)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no address Found, net.InterfaceAddrs: %v", addresses)
	}
	//only need first
	return addresses, nil
}
