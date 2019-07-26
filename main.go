package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type employeeMap map[string][]map[string]interface{}
type employeeDataMap map[string]interface{}
type employeeSlice []map[string]interface{}

// Space struct defined
type Space struct {
	name       string
	maxPersons int
	members    employeeSlice
	spaceType  string
}

func (sp *Space) addMembers(e **employeeSlice) {
	var lastItem employeeDataMap
	fmt.Println(len(**e))
	for i := 0; i < sp.maxPersons; i++ {
		// pop an item from the slice
		if len(**e) > 0 {
			// read the last item from the slice
			lastItem = (**e)[len(**e)-1]

			// remove the last item of the slice
			**e = (**e)[:len(**e)-1]
		}
		sp.members = append(sp.members, lastItem)
	}
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

	employees := make(employeeMap)
	var staffSlice employeeSlice

	var maleFellowSlice employeeSlice
	var femaleFellowSlice employeeSlice

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if line[3] == "STAFF" {
			staffMap := make(employeeDataMap)
			staffMap["name"] = fmt.Sprintf("%s %s ", line[0], line[1])
			staffMap["gender"] = line[2]
			staffMap["position"] = line[3]

			staffSlice = append(staffSlice, staffMap)
		}

		if line[3] == "FELLOW" {
			maleFellowsMap := make(employeeDataMap)
			femaleFellowsMap := make(employeeDataMap)
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

func allocateToOffice(e *employeeSlice, offices []string) {
	for _, office := range offices {
		fmt.Println(office)
		fmt.Println(&e, "intial address")

		ofc := &Space{
			name:       office,
			maxPersons: 6,
			spaceType:  "officeRoom",
		}

		ofc.addMembers(&e)
	}

	fmt.Println(e)
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
	var eSlice employeeSlice
	for _, val := range e {
		eSlice = append(eSlice, val...)
	}

	// shuffle employee slice for random office allocation
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(eSlice), func(i, j int) { eSlice[i], eSlice[j] = eSlice[j], eSlice[i] })

	// declare the hostels and office here
	// femaleHostel := []string{"ruby", "platinum", "jade", "pearl", "diamond"}
	// maleHostel := []string{"topaz", "silver", "gold", "onyx", "opal"}
	office := []string{"Carat", "Anvil", "Crucible", "Kiln", "Forge", "Foundry", "Furnace", "Boiler", "Mint", "Vulcan"}
	go allocateToOffice(&eSlice, office)
	fmt.Scanln()

}
