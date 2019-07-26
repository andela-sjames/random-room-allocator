package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type employeeMap map[string][]map[string]interface{}
type membersSpace map[string]interface{}

type space struct {
	name       string
	maxPersons int
}

type office struct {
	*space
	officeMembers membersSpace
}

type maleRoom struct {
	*space
	maleMembers membersSpace
}

type femaleRoom struct {
	*space
	femaleMembers membersSpace
}

func closeFile(f *os.File) {
	fmt.Println("closing")
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func generateObject(f *os.File) employeeMap {
	fmt.Println("generating data object")

	employees := make(map[string][]map[string]interface{})
	var staffSlice []map[string]interface{}

	var maleFellowSlice []map[string]interface{}
	var femaleFellowSlice []map[string]interface{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if line[3] == "STAFF" {
			staffMap := make(map[string]interface{})
			staffMap["name"] = fmt.Sprintf("%s %s ", line[0], line[1])
			staffMap["gender"] = line[2]
			staffMap["position"] = line[3]

			staffSlice = append(staffSlice, staffMap)
		}

		if line[3] == "FELLOW" {
			maleFellowsMap := make(map[string]interface{})
			femaleFellowsMap := make(map[string]interface{})
			livingSpace := true
			if line[3] == "Y" {
				livingSpace = false
			}

			if line[2] == "M" {
				maleFellowsMap["name"] = fmt.Sprintf("%s %s", line[0], line[1])
				maleFellowsMap["gender"] = line[2]
				maleFellowsMap["position"] = line[3]
				maleFellowsMap["livingSpace"] = livingSpace

				maleFellowSlice = append(maleFellowSlice, maleFellowsMap)
			}

			if line[2] == "F" {
				femaleFellowsMap["name"] = fmt.Sprintf("%s %s", line[0], line[1])
				femaleFellowsMap["gender"] = line[2]
				femaleFellowsMap["position"] = line[3]
				femaleFellowsMap["livingSpace"] = livingSpace

				femaleFellowSlice = append(femaleFellowSlice, femaleFellowsMap)
			}

		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	employees["staff"] = staffSlice
	employees["maleFellows"] = maleFellowSlice
	employees["femaleFellows"] = femaleFellowSlice
	return employees
}

// FileParser gets data from inputfile.
type fileParser struct {
	filepath string
}

func (fp *fileParser) GetEmployees() employeeMap {
	f, err := os.Open(fp.filepath)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		panic(err)
	}
	defer closeFile(f)
	return generateObject(f)
}

func main() {
	inputfile := &fileParser{filepath: "inputA.txt"}
	e := inputfile.GetEmployees()
	fmt.Println(e)
}
