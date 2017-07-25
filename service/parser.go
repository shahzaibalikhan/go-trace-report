package service

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"sync"
	"strconv"
)

type WorkFlowData struct {
	WorkFlowNumber int
	Data           []Raw
}

type Raw struct {
	StackFrames map[string]struct {
		Name   string `json:"name"`
		Parent int `json:"parent"`
	} `json:"stackFrames"`
}

type workFlow map[string][]Raw

func ReadDirForParsing(workFlowDir string) []string {
	var totalRawDumps []string
	var filesWritten []string
	var groupFilePathsByWorkFlow = make(map[int][]string)

	_ = filepath.Walk(workFlowDir, func(path string, f os.FileInfo, err error) error {
		// Fix with regex later
		if strings.Contains(path, "Workflow") && strings.Contains(path, ".out") {
			totalRawDumps = append(totalRawDumps, path)
		}

		return nil
	})

	for i := range (totalRawDumps) {
		workflowNo := strings.Split(strings.Split(totalRawDumps[i], "Workflow")[1], "/")[0]
		number, _ := strconv.ParseInt(workflowNo, 10, 0)
		groupFilePathsByWorkFlow[int(number)] = append(groupFilePathsByWorkFlow[int(number)], totalRawDumps[i])
	}

	fmt.Printf("Total workflow %d\n", len(groupFilePathsByWorkFlow))

	fmt.Printf("Total workflow dumps %d\n", len(totalRawDumps))
	timeToSleep := 1 * time.Second
	fmt.Println("wait for `go trace` to start: ", timeToSleep)

	var wg sync.WaitGroup

	wg.Add(len(totalRawDumps))

	for workFlowNo := range groupFilePathsByWorkFlow {
		for fun := range groupFilePathsByWorkFlow[workFlowNo] {
			go func(workFlowNo int, fun int) {
				defer wg.Done()
				randPort := randomPortNumber()
				cmdArgs := []string{fmt.Sprintf("-http=localhost:%d", randPort),
					groupFilePathsByWorkFlow[int(workFlowNo)][fun]}
				go executeCommand(cmdArgs)

				time.Sleep(timeToSleep)

				requestData := requestForJSONTrace(randPort)

				if requestData != nil {
					data := requestJSONUnmarshal(requestData)
					time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)

					dataToWrite := ""
					for d := range data.StackFrames {
						dataToWrite = fmt.Sprintf("%s%s\n", dataToWrite, data.StackFrames[d].Name)
					}
					fileToWrite := fmt.Sprintf("/tmp/functions-%d-%d", workFlowNo, fun)
					filesWritten = append(filesWritten, fileToWrite)
					_ = ioutil.WriteFile(fileToWrite, []byte(dataToWrite), 0644)
				}

			}(int(workFlowNo), fun)
		}
	}

	wg.Wait()

	return filesWritten
}

func executeCommand(cmd []string) {
	absPath, _ := filepath.Abs("./tracy")

	if err := exec.Command(absPath, cmd...).Run(); err != nil {
		fmt.Println("Error occurred while running command")
		fmt.Fprintln(os.Stderr, err)
	}
}

func requestForJSONTrace(port int) []byte {
	res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/jsontrace", port))
	if err != nil {
		fmt.Println("Request error")
		fmt.Println(err)
		return nil
	} else {
		robots, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Body error")
			fmt.Println(err)
		}
		res.Body.Close()
		return robots
	}
}

func requestJSONUnmarshal(jsonBlob []byte) Raw {
	r := Raw{}
	err := json.Unmarshal(jsonBlob, &r)

	if err != nil {
		panic(err)
	}
	return r
}

func randomPortNumber() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(10000 + 3000) + 3000
}