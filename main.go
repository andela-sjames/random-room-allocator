package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Employee struct defined
type employee struct {
	name     string
	gender   string
	position string
	optSpace bool
}

func closeFile(f *os.File) {
	fmt.Println("closing")
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func generateObject(f *os.File) {
	fmt.Println("generating data object")

	employees := make(map[string][]map[string]string)

	var fellowSlice []map[string]string
	var staffSlice []map[string]string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if line[2] == "STAFF" {
			staffMap := make(map[string]string)
			staffMap["name"] = line[0]
			staffMap["gender"] = line[1]
			staffMap["position"] = line[2]

			staffSlice = append(staffSlice, staffMap)
		}

		if line[2] == "FELLOW" {
			fellowsMap := make(map[string]string)
			optSpace := "true"
			if line[3] == "Y" {
				optSpace = "false"
			}

			fellowsMap["name"] = line[0]
			fellowsMap["gender"] = line[1]
			fellowsMap["position"] = line[2]
			fellowsMap["optSpace"] = optSpace

			fellowSlice = append(fellowSlice, fellowsMap)
		}

		// fmt.Printf("%q\n", line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	employees["staff"] = staffSlice
	employees["fellow"] = fellowSlice

	fmt.Printf("%q\n", employees)
}

// FileParser gets data from inputfile.
type fileParser struct {
	filepath string
}

func (fp *fileParser) GetEmployees() {
	f, err := os.Open(fp.filepath)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer closeFile(f)
	generateObject(f)

}

func main() {
	inputfile := &fileParser{filepath: "inputA.txt"}
	inputfile.GetEmployees()
	// fmt.Printf("%q\n", employees)
}
