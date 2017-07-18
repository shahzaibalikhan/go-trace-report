package bitch

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//var root string = "/home/shahzaib/Downloads/traces/"
func ReadDirForParsing(root string) {
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		// Fix with regex later
		if strings.Contains(path, "Workflow") && strings.Contains(path, ".out") {
			fmt.Printf("File found: %s\n", path)
			go fileParser(path)
		}
		return nil
	})
	fmt.Println("All files read! returned %v\n", err)
}

func fileParser(path string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randPort := r.Intn(10000+3000) + 3000
	fmt.Println("%d", randPort)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	executeCommand(fmt.Sprintf("go tool trace -http=localhost:%d %s", randPort, path), wg)

}

func executeCommand(cmd string, wg *sync.WaitGroup) {
	fmt.Println(cmd)

	out, err := exec.Command(cmd).Output()

	if err != nil {
		//fmt.Println("error occured")
		//fmt.Println("%s", err)
	} else {

		fmt.Printf("%s", out)
		wg.Done()
		time.Sleep(30000)
	}
}