package bitch

import (
	"path/filepath"
	"strings"
	"os"
	"fmt"
	"io/ioutil"
)

// Crypto ref Finder
func CryptoFinder(root string) {
	cryptoRef := make(map[string][]string)

	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {

		if strings.Contains(path, "traces") && strings.Contains(path, "function") && strings.Contains(path,
			"txt") {
			fmt.Printf("File found: %s\n", path)
			workflowNo := strings.Split(strings.Split(path, "Workflow")[1], "/")[0]
			fmt.Printf("work flow number: %s\n", workflowNo)

			f, err := ioutil.ReadFile(path)

			if err != nil {
				panic(err)
			}

			data := strings.Split(string(f), "\n")

			for i := 0; i < len(data); i++ {
				if strings.Contains(data[i], "crypto/") {
					cryptoRef[workflowNo] = append(cryptoRef[workflowNo], data[i])
				}
			}

		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	for k := range cryptoRef {
		fmt.Printf("Workflow: %s \n", k);
		fmt.Println("===========================");
		err := ioutil.WriteFile("/tmp/data" + k, []byte(strings.Join(cryptoRef[k], "\n")), 0644)
		if err != nil {
			panic(err)
		}
	}
}