package main

import (
	"os"
	"log"
	"fmt"
	"net"
	"time"
	"path/filepath"
	"../qlm/df"
	"../fsnotify"
)

func main() {

	// Command line
	if len(os.Args) != 4 {
		fmt.Printf("Usage: sensor-fs <root directory> <core address> <core port>\n")
		return
	}
	root := os.Args[1]
	address := os.Args[2]
	port := os.Args[3]

	// Connect to core
	conn, err := net.Dial("tcp", address + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Read the fs tree starting at root to a qlm struct
	qlm, _ := ParseFs(root);

	// Initial message
	bytes, _ := df.Marshal(qlm)
	conn.Write(bytes)

	// Create a new fs watcher
	watcher, _ := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Watch the whole tree
	err = filepath.Walk(root, func(path string, _ os.FileInfo, _ error) error {
		log.Print("watching: ", path)
		err = watcher.Add(path)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Send an update whenever the fs changes
	done := make(chan bool)
	go func() {
		var last fsnotify.Event

		then := time.Now()
		epsilon, _ := time.ParseDuration("5ms")

		for {
			select {
			case event := <-watcher.Events:
				elapsed := time.Since(then)
				then = time.Now()
				log.Print(event.Name, " ", event.Op, " ", elapsed)
				if (event.Name == last.Name && elapsed < epsilon) {
					// Ignore multiple events from the same operation
					// Hack needed to get around fsnotify sending 
					// multiple events for each fs operation
					// eg. date > file triggers 4 write events
				} else {
					go func() {
						// wait a moment for the operation to complete
						time.Sleep(epsilon)
						qlm, _ = ParseFs(root);
						bytes, _ := df.Marshal(qlm)
						conn.Write(bytes)
					}()
				}
				last = event
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	<-done
}
