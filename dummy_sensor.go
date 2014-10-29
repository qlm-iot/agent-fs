package main

import "fmt"
import "log"
import "os"
import "path/filepath"
import "time"

func main() {

	// Check number of arguments
	if len(os.Args) != 2 {
		fmt.Printf("Usage: dummy_sensor <empty directory>\n")
		return
	}

	// Check that argument is an empty directory
	// Does it exist?
	root, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	// Is it a directory?
	if !root.IsDir() {
		log.Fatal(root.Name(), " is not a directory")
	}
	// Is it empty?
	err = filepath.Walk(root.Name(), func (path string, _ os.FileInfo, _ error) error {
		//fmt.Printf("%s %s\n", root.Name(), path)
		if path != root.Name() {
			return fmt.Errorf("%s is not empty", root.Name())
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Absolute path of the root directory
	rootpath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Generate Object ID from hostname and pid
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	objectid := fmt.Sprintf("DummySensor-%s-%d", hostname, os.Getpid())

	// Generate directory structure
	path := rootpath + "/Objects/" + objectid
	//fmt.Printf("%s\n", path)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Create a file that represents the Object attributes
	objectattributefile, err := os.Create(path + "/attributes")
	if err != nil {
		log.Fatal(err)
	}
	_, err = objectattributefile.WriteString("<Object udef=\"a.p.9\" type=\"Application Software Product\">\n")
	if err != nil {
		log.Fatal(err)
	}			
	err = objectattributefile.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Create a file that represents the id element(s)
	objectidfile, err := os.Create(path + "/id")
	if err != nil {
		log.Fatal(err)
	}
	_, err = objectidfile.WriteString("<id>" + objectid + "<\\id>\n")
	if err != nil {
		log.Fatal(err)
	}			
	err = objectidfile.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Create a file that represents the description element
	objectdescriptionfile, err := os.Create(path + "/description")
	if err != nil {
		log.Fatal(err)
	}
	_, err = objectdescriptionfile.WriteString("<description>A dummy sensor<\\description>\n")
	if err != nil {
		log.Fatal(err)
	}			
	err = objectdescriptionfile.Close()
	if err != nil {
		log.Fatal(err)
	}

	// InfoItem name
	infoitemname := "SystemTime"

	// Create a file that represents the InfoItem
	infoitemfile, err := os.Create(path + "/" + infoitemname)
	if err != nil {
		log.Fatal(err)
	}

	// Opening tag for InfoItem
	infoitemopen := "<InfoItem name=\"" + infoitemname + "\">"

	// Description for InfoItem
	infoitemdescription := "<description>System time for " + hostname + "<\\description>"

	// Periodically update the InfoItem
	ticker := time.NewTicker(time.Millisecond * 1000)
	go func() {
		for t := range ticker.C {
			_, err = infoitemfile.Seek(0, 0)
			if err != nil {
				log.Fatal(err)
			}
			_, err = infoitemfile.WriteString(infoitemopen + "\n" + "  " + infoitemdescription + "\n" + "  <value>" + t.String() + "<\\value>\n<\\InfoItem>\n")
			if err != nil {
				log.Fatal(err)
			}			
		}
	}()


	// Informative message
	fmt.Printf("%s running at %s, press enter to quit\n", objectid, rootpath)	

	// Wait for keypress
	_, _ = fmt.Scanf("\n")

	// Stop the ticker
	ticker.Stop()

	// Close the file
	err = infoitemfile.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Clean up
	//fmt.Printf("%s\n", rootpath + "/Objects/")
	err = os.RemoveAll(rootpath + "/Objects/")
	if err != nil {
		log.Fatal(err)
	}
}
