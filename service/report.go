package service

import (
	"fmt"
	"sort"
	"strings"
)

type val struct {
	funCount float64
	fipCount float64
}

type report map[int]val

// GenerateReport Generates final report containing total function count and FIPS call references
func GenerateReport(FilesToRead []string) string {
	var re = report{}

	for i := range FilesToRead {
		totalFunCounts, fipsCount := returnFipsRef(fileReader(FilesToRead[i]))

		q := strings.Split(FilesToRead[i], "-")
		q = q[len(q)-2:]

		index := strToInt(q[0])

		if re[index].funCount == 0 {
			re[index] = val{fipCount: float64(fipsCount), funCount: float64(totalFunCounts)}
		} else {
			v := re[index]
			v.funCount = v.funCount + float64(totalFunCounts)
			v.fipCount = v.fipCount + float64(fipsCount)
			re[index] = v
		}
	}

	return createReport(re)
}

func returnFipsRef(data string) (int, int) {
	var counter = 0

	var a = strings.Split(string(data), "\n")

	for i := range a {
		if strings.Contains(a[i], "fips") {
			counter++
		}
	}

	return len(a), counter
}

func sortIndex(m report) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}

func createReport(re report) string {
	generatedReport := "\n************************************\n\n"

	for i := range sortIndex(re) {
		i = i + 1
		generatedReport = generatedReport + fmt.Sprintf("Workflow: %d \n", i)
		generatedReport = generatedReport + fmt.Sprintf("fips hit percentage: %.2f \n", re[i].fipCount/re[i].funCount*100)
		generatedReport = generatedReport + fmt.Sprintln("===========================")
	}

	return generatedReport
}
