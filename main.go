package main

import (
	"github.com/shahzaibalikhan/transpile/service"
)

func main() {
	root := "/home/shahzaib/Downloads/notary-automation/traces/"
	// Parser for trace files
	service.ReadDirForParsing(root)
}
