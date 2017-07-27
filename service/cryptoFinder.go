package service

import (
	"io/ioutil"
	"fmt"
	"strings"
)

func LocateCryptoUsages(FilesToRead []string) []string {
	var filesWritten []string

	for i := range FilesToRead {
		dat, _ := ioutil.ReadFile(FilesToRead[i])
		d := findCryptoUsages(string(dat))

		dataToWrite := fmt.Sprintf("%s\n", strings.Join(d, "\n"))
		_ = ensureDir("/crypto")

		fileToWrite := strings.Replace(FilesToRead[i], "/functions-", "/crypto-", 1)
		fileToWrite = strings.Replace(fileToWrite, "/functions/", "/crypto/", 1)

		writeFile(fileToWrite, dataToWrite)
		filesWritten = append(filesWritten, fileToWrite)
	}
	return filesWritten
}

func findCryptoUsages(data string) []string {
	var r []string

	var a = strings.Split(string(data), "\n")
	for i := range a {
		if strings.Contains(a[i], "crypto/") {
			r = append(r, a[i])
		}
	}

	return r
}
