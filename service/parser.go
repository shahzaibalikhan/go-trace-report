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
)

type Raw struct {
	StackFrames map[string] struct {
		Name string `json:"name"`
		Parent int `json:"parent"`
	} `json:"stackFrames"`
}

type workFlow map[string][]Raw
var cryptoRef = make(map[string][]string)
var totalFunctionsCount = make(map[string]int)
var fipsFunctionCount = make(map[string]int)

//var root string = "/home/shahzaib/Downloads/traces/"
func ReadDirForParsing(root string) {
	var filePaths []string

	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		// Fix with regex later
		if strings.Contains(path, "Workflow") && strings.Contains(path, ".out") {
			filePaths = append(filePaths, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Total raw dumps %d\n", len(filePaths))
	timeToSleep := (time.Duration(len(filePaths)) % 7) * time.Second
	fmt.Println("wait for `go trace` to start: ", timeToSleep)

	var wg sync.WaitGroup

	wg.Add(len(filePaths))

	for p := range filePaths {

		go func() {
			defer wg.Done()

			workflowNo := strings.Split(strings.Split(filePaths[p], "Workflow")[1], "/")[0]

			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			randPort := r.Intn(10000+3000) + 3000

			cmdArgs := []string{"tool", "trace", fmt.Sprintf("-http=localhost:%d", randPort), filePaths[p]}

			go executeCommand(cmdArgs)

			time.Sleep(timeToSleep)

			requestData := requestForJSONTrace(randPort)

			if requestData != nil {
				temp := requestJSONUnmarshal(requestData)
				time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
				cryptoFinder(workflowNo, temp)
			}
		}()
	}

	wg.Wait()

	report()

	fmt.Println("All files read!")
}

func executeCommand(cmd []string)  {
	if err := exec.Command("go", cmd...).Run(); err != nil {
		fmt.Println("Error occurred while running command")
		fmt.Fprintln(os.Stderr, err)
	}
}

func requestForJSONTrace(port int) []byte{
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

func cryptoFinder(workflowNo string, data Raw) {

	for i := range data.StackFrames {
		if strings.Contains(data.StackFrames[i].Name, "crypto/") {
			cryptoRef[workflowNo] = append(cryptoRef[workflowNo], data.StackFrames[i].Name)

			if strings.Contains(data.StackFrames[i].Name, "fips") {
				fipsFunctionCount[workflowNo]++
			}

			totalFunctionsCount[workflowNo]++
		}
	}
}

func report() {
	for k := range totalFunctionsCount {

		fipsCallPercent := 0
		if totalFunctionsCount[k] > 0 {
			fipsCallPercent = fipsFunctionCount[k] / totalFunctionsCount[k] * 100;
		}
		fmt.Printf("Workflow: %s \n", k);
		fmt.Printf("fips hit percentage: %d \n", fipsCallPercent);
		fmt.Println("===========================");
	}
}
