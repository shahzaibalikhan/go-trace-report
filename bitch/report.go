package bitch

import (
	"path/filepath"
	"strings"
	"os"
	"fmt"
	"io/ioutil"
)

// Report Generator
func WorkFlowReport(root string) {
	totalFunctionsCount := make(map[string]int)
	fipsFunctionCount := make(map[string]int)

	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {

		if strings.Contains(path, "traces") && strings.Contains(path, "crypto") && strings.Contains(path,
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
				if strings.Contains(data[i], "fips") {
					fipsFunctionCount[workflowNo]++
				}
				totalFunctionsCount[workflowNo]++
			}

		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error occurred while geenrating report: %s", err)
	}

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
