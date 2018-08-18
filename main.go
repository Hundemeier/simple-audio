package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/fatih/color"
)

const poolDir = "./pool"
const sampleRate = beep.SampleRate(44100)

//slotMap contains the mapping of a slot to a filename in the poolDir
var slotMap = newSyncSlotMap()

var port = flag.Uint("port", 8080, "the port on which the webinterface is listening. Only use port 80, when no other application is using this port!")

func main() {
	slotMap.onPlayerChange = func(slot uint, play *player) {
		item, err := slotMap.getSlotItem(slot)
		if err == nil {
			go websockets.writeToWebsockets(item)
		}
	}

	flag.Parse()

	initSpeaker()

	makeSlotMapFromConfig()

	initGraphql()
	initWebService()
}

func initSpeaker() {
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/15))
	if err != nil {
		log.Fatal(err)
	}
}

func initWebService() {
	fmt.Println("Starting server...")
	server := http.Server{
		Addr: fmt.Sprintf(":%v", *port),
	}
	fmt.Println("Serving at:")
	fmt.Printf("\thttp://127.0.0.1:%v\n", *port)
	addrs, _ := getMyInterfaceAddr()
	for _, addr := range addrs {
		fmt.Printf("\thttp://%v%v\n", addr, server.Addr)
	}
	//fmt.Println("Close with \033[47;30m Ctrl+C \033[0m")
	fmt.Print("Close with ")
	color.Set(color.FgBlack)
	color.Set(color.BgWhite)
	fmt.Print(" Ctrl+C ")
	color.Unset()
	fmt.Print("\n")

	//http.Handle("/", http.FileServer(http.Dir("webgui/")))
	http.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "webgui/dist"}))
	http.HandleFunc("/pool/upload", uploadHandler)
	http.HandleFunc("/websocket", handleWebsocket)
	log.Fatal(server.ListenAndServe())
}

func makeSlotMapFromConfig() {
	conf := readConfig()
	for _, item := range conf.SlotMap {
		slotMap.setFileForSlot(item.Slot, item.Filename)
		player := slotMap.getPlayer(item.Slot)
		if player != nil {
			player.setVolume(item.Player.Volume)
			player.loop = item.Player.Loop
		}
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		httpFile, head, err := r.FormFile("uploadfile")
		if err != nil {
			//fmt.Println("Err:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer httpFile.Close()
		if head.Size != 0 {
			//check if dir "pool" exists, otherwise create directory
			if _, err := os.Stat(poolDir); os.IsNotExist(err) {
				err = os.MkdirAll(poolDir, 0755)
				if err != nil {
					fmt.Println("Error:", err)
				}
			}

			path := poolDir + "/" + head.Filename
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				//fmt.Println("Err:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer f.Close()
			_, err = io.Copy(f, httpFile)
			if err != nil {
				//fmt.Println("Err:", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}
