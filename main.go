package main

import (
	"github.com/shahzaibalikhan/transpile/service"
	"fmt"
)

func main() {
	//workFlowDir := "/home/shahzaib/Downloads/notary-automation/traces/"
	workFlowDir := "/home/shahzaib/Downloads/test/traces/"
	// Parser for trace files
	filesWritten := service.ReadDirForParsing(workFlowDir)

	fmt.Println("Function Files written")
	fmt.Println(filesWritten)

	cryptoRefWritten := service.CryptoFinder(filesWritten)

	fmt.Println("Crypto Files written")
	fmt.Println(cryptoRefWritten)

	report := service.GenerateReport(cryptoRefWritten)

	fmt.Println("Final report")
	fmt.Println(report)

}
