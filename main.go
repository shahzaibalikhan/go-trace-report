package main

import (
	"github.com/shahzaibalikhan/transpile/service"
	"fmt"
	"flag"
	"os"
	"strings"
	"time"
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
		fmt.Println("Traversing Dir for workflows: ", service.WorkFlowDir)
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

	fmt.Printf("Time to execute %s", time.Since(start))
}
