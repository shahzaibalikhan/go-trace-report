package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func unixTime() int64 {
	return time.Now().UnixNano()
}

func randomPortNumber() int {
	r := rand.New(rand.NewSource(unixTime()))
	return r.Intn(10000+3000) + 3000
}

func strToInt(s string) int {
	number, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("str to int conversion failed: %s, err: %v", s, err))
	}
	return int(number)
}

// AbsPath returns Absolute path to provided location
func AbsPath(p string) string {
	absPath, err := filepath.Abs(p)

	if err != nil {
		panic(fmt.Sprintf("Absolute path conversion failed: %s, err: %v", p, err))
	}

	ValidateDirExist(absPath)

	return absPath
}

func writeFile(path, data string) {
	err := ioutil.WriteFile(path, []byte(data), 0644)

	if err != nil {
		panic(fmt.Sprintf("Error occurred while writing file: %s err: %v", path, err))
	}
}

func ensureDir(dir string) string {
	dirPath := filepath.Join(WorkFlowDir, dir)
	os.MkdirAll(dirPath, os.ModePerm)
	return dirPath
}

func fileReader(p string) string {
	data, err := ioutil.ReadFile(p)

	if err != nil {
		panic(fmt.Sprintf("Error occurred while reading file: %s err: %v", p, err))
	}

	return string(data)
}

// ValidateDirExist checks if provided dir exist
func ValidateDirExist(p string) {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		panic(fmt.Sprintf("Dir does not exist: %s err: %v", p, err))
	}
}

// Logger used for debug logs
func Logger(s string) {
	if os.Getenv("DEBUG") == "1" {
		log.Println(s)
	}
}
