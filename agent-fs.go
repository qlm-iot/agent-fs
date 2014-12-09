package main

import (
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/qlm-iot/qlm/df"
	"github.com/qlm-iot/qlm/mi"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func wsServerConnector(address string) (chan []byte, chan []byte) {
	send := make(chan []byte)
	receive := make(chan []byte)
	go func() {
		for {
			select {
			case rawMsg := <-send:
				var h http.Header
				println(string(rawMsg))
				conn, _, err := websocket.DefaultDialer.Dial(address, h)
				if err == nil {
					if err := conn.WriteMessage(websocket.BinaryMessage, rawMsg); err != nil {
						receive <- []byte(err.Error())
					}
					_, content, err := conn.ReadMessage()
					if err == nil {
						receive <- content
					} else {
						receive <- []byte(err.Error())
					}
				} else {
					receive <- []byte(err.Error())
				}
			}
		}
	}()
	return send, receive
}

func createWriteRequest(qlm  []byte) []byte {
	ret, _ := mi.Marshal(mi.OmiEnvelope{
		Version: "1.0",
		Ttl: -1,
		Write: &mi.WriteRequest{
			MsgFormat: "odf",
			TargetType: "device",
			Message: &mi.Message{
				Data: string(qlm),
			},
		},
	})
	return ret
}

func main() {

	// Command line
	if len(os.Args) != 3 {
		fmt.Printf("Usage: sensor-fs <root directory> <core address>\n")
		return
	}
	root := os.Args[1]
	address := os.Args[2]

	send, receive := wsServerConnector(address)

	go func(channel chan []byte) {
		for {
			msg := <-channel
			fmt.Printf(string(msg))
		}
	}(receive)

	// Read the fs tree starting at root to a qlm struct
	qlm, _ := ParseFs(root)

	// Initial message
	//bytes, _ := df.Marshal(qlm)
	//send <- bytes

	// Create a new fs watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Watch the whole tree
	err = filepath.Walk(root, func(path string, _ os.FileInfo, _ error) error {
		//log.Print("watching: ", path)
		err = watcher.Add(path)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	changed := false

	// Send an update whenever the fs changes
	done := make(chan bool)
	go func() {

		//then := time.Now()
		epsilon, _ := time.ParseDuration("10ms")

		for {
			select {
			case <-watcher.Events:
				//elapsed := time.Since(then)
				//then = time.Now()
				//log.Print(event.Name, " ", event.Op, " ", elapsed)
				if changed {
					// Ignore multiple events in quick succession
					// Any change will cause the whole tree to be sent
					// so we don't even check what has changed
				} else {
					changed = true
					go func() {
						// wait a moment for the operation to complete
						time.Sleep(epsilon)
						qlm, _ = ParseFs(root)
						msg, _ := df.Marshal(qlm)
						bytes := createWriteRequest(msg)
						send <- bytes
						changed = false
					}()
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	<-done
}
