package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ReadDirForParsing Reads all work log files, return all array of files created after parsing
func ReadDirForParsing(workFlowDir string) []string {
	var totalRawDumps []string
	var filesWritten []string
	var groupFilePathsByWorkFlow = make(map[int][]string)

	err := filepath.Walk(workFlowDir, func(path string, f os.FileInfo, err error) error {
		// Fix with regex later
		if strings.Contains(path, "Workflow") && strings.Contains(path, ".out") {
			totalRawDumps = append(totalRawDumps, path)
		}

		return nil
	})

	if err != nil {
		panic(fmt.Sprintf("unable to traverse dir: %s err: %v", workFlowDir, err))
	}

	for i := range totalRawDumps {
		workflowNo := strToInt(strings.Split(strings.Split(totalRawDumps[i], "Workflow")[1], "/")[0])
		groupFilePathsByWorkFlow[workflowNo] = append(groupFilePathsByWorkFlow[workflowNo], totalRawDumps[i])
	}

	if len(groupFilePathsByWorkFlow) == 0 {
		fmt.Println("No workflows found in dir", workFlowDir)
		os.Exit(1)
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
					data := jsonToRawObject(requestData)

					dataToWrite := ""
					for d := range data.StackFrames {
						dataToWrite = fmt.Sprintf("%s%s\n", dataToWrite, data.StackFrames[d].Name)
					}

					dir := ensureDir("/functions")
					fileToWrite := fmt.Sprintf("%s/functions-%d-%d", dir, workFlowNo, fun)
					writeFile(fileToWrite, dataToWrite)
					filesWritten = append(filesWritten, fileToWrite)
				}

			}(int(workFlowNo), fun)
		}
	}

	wg.Wait()

	return filesWritten
}

func executeCommand(cmd []string) {
	if err := exec.Command(AbsPath(tracyPath), cmd...).Run(); err != nil {
		Logger("Error occurred while running command")
		Logger(fmt.Sprintf("%v", err))
	}
}

func requestForJSONTrace(port int) []byte {
	res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/jsontrace", port))
	if err != nil {
		Logger("Request error")
		Logger(fmt.Sprintf("%v", err))
		return nil
	}
	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		Logger("Body error")
		Logger(fmt.Sprintf("%v", err))
	}
	res.Body.Close()
	return robots
}

func jsonToRawObject(jsonBlob []byte) Raw {
	r := Raw{}
	err := json.Unmarshal(jsonBlob, &r)

	if err != nil {
		panic(err)
	}
	return r
}
