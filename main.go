package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shahzaibalikhan/transpile/service"
)

const usageMessage = "" +
	`Usage of 'go tool trace':
Given a trace file produced by 'go test':
	go test -trace=trace.out pkg

Args:
	path to workflow dir
`

func main() {
	start := time.Now()

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageMessage)
		os.Exit(2)
	}

	flag.Parse()

	switch flag.NArg() {
	case 1:
		service.WorkFlowDir = service.AbsPath(flag.Arg(0))
		fmt.Println("Traversing directory for workflows: ", service.WorkFlowDir)
	default:
		flag.Usage()
	}

	// Parser for trace files
	filesWritten := service.ReadDirForParsing(service.WorkFlowDir)

	service.Logger("Function Files written")
	service.Logger(strings.Join(filesWritten, "\n"))

	cryptoRefWritten := service.LocateCryptoUsages(filesWritten)

	service.Logger("Crypto Files written")
	service.Logger(strings.Join(cryptoRefWritten, "\n"))

	report := service.GenerateReport(cryptoRefWritten)

	service.Logger("Final report")
	fmt.Println(report)

	fmt.Println("Time to execute ", time.Since(start))
	os.Exit(0)
}
