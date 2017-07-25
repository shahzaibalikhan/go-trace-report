package service

import (
	"io/ioutil"
	"fmt"
	"strings"
	"go-crypto-swap/.GOROOT/src/strconv"
	"sort"
)

type val struct {
	funCount float64
	fipCount float64
}

type report map[int]val

func GenerateReport(FilesToRead []string) string {
	var generatedReport string
	var re = report{}

	for i := range FilesToRead {
		dat, _ := ioutil.ReadFile(FilesToRead[i])
		funCount, fipsCount := returnFipsRef(string(dat))

		q := strings.Split(FilesToRead[i], "-")
		q = q[len(q) - 2:]
		workFlowNo, _ := strconv.ParseInt(q[0], 10, 32)
		index := int(workFlowNo)
		if re[index].fipCount == 0 {
			re[index] = val{fipCount:float64(fipsCount), funCount:float64(funCount)}
		} else {
			v := re[index]
			v.funCount = v.funCount + float64(funCount)
			v.fipCount = v.fipCount + float64(fipsCount)
			re[index] = v

		}
	}

	for i := range sortIndex(re) {
		i = i + 1
		generatedReport = generatedReport + fmt.Sprintf("Workflow: %d \n", i);
		generatedReport = generatedReport + fmt.Sprintf("fips hit percentage: %.2f \n", re[i].fipCount / re[i].funCount * 100);
		generatedReport = generatedReport + fmt.Sprintln("===========================");
	}
	return generatedReport
}

func returnFipsRef(data string) (int, int) {
	var fipsCounter = 0

	var a = strings.Split(string(data), "\n")
	for i := range a {
		if strings.Contains(a[i], "fips") {
			fipsCounter++
		}
	}

	return len(a), fipsCounter
}

func sortIndex(m report) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}