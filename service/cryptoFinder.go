package service

import (
	"io/ioutil"
	"fmt"
	"strings"
)

func CryptoFinder(FilesToRead []string) []string {
	var filesWritten []string

	for i := range FilesToRead {
		dat, _ := ioutil.ReadFile(FilesToRead[i])
		d := returnCryptoRef(string(dat))

		dataToWrite := fmt.Sprintf("%s\n", strings.Join(d, "\n"))

		fileToWrite := strings.Replace(FilesToRead[i], "/functions", "/crypto", 1)
		_ = ioutil.WriteFile(fileToWrite, []byte(dataToWrite), 0644)
		filesWritten = append(filesWritten, fileToWrite)
	}
	return filesWritten
}

func returnCryptoRef(data string) []string {
	var r []string

	var a = strings.Split(string(data), "\n")
	for i := range a {
		if strings.Contains(a[i], "crypto/") {
			r = append(r, a[i])
		}
	}

	return r
}
