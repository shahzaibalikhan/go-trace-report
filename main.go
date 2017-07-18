package main

import (
	"github.com/shahzaibalikhan/transpile/bitch"
)

func main() {
	root := "/home/shahzaib/Downloads/notary-automation/traces/"
	// Parser for trace files
	//bitch.ReadDirForParsing(root)

	// Report generator
	//bitch.WorkFlowReport(root)

	bitch.CryptoFinder(root)
}