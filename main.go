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
	name        string
	gender      string
	position    string
	wantsLiving bool
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

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		fmt.Printf("%q\n", line[0])
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
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
}
