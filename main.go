package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type employeeMap map[string][]map[string]interface{}
type employeeDataMap map[string]interface{}
type employeeSlice []map[string]interface{}

// Space struct defined
type Space struct {
	Name       string        `json:"Name"`
	MaxPersons int           `json:"MaxPersons"`
	Members    employeeSlice `json:"Members"`
	SpaceType  string        `json:"Type"`
}

func (sp *Space) addMembers(e **employeeSlice) {
	var lastItem employeeDataMap
	for i := 0; i < sp.MaxPersons; i++ {
		if len(**e) > 0 {
			// read the last item from the slice
			lastItem = (**e)[len(**e)-1]
			// remove the last item of the slice
			**e = (**e)[:len(**e)-1]
		}
		sp.Members = append(sp.Members, lastItem)
	}
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

func allocateToOffice(e *employeeSlice, offices []string, unAllocatedToOffice chan<- interface{}, wg *sync.WaitGroup) {
	var file []byte
	var allocationSlice []Space

	for _, office := range offices {
		spc := &Space{
			Name:       office,
			MaxPersons: 6,
			SpaceType:  "officeRoom",
		}

		spc.addMembers(&e)

		// append allocation to slice
		allocationSlice = append(allocationSlice, *spc)
	}

	// write allocation to json file
	file, _ = json.MarshalIndent(allocationSlice, "", " ")
	_ = ioutil.WriteFile("officeAllocation.json", file, 0644)

	unAllocatedToOffice <- e
	wg.Done()
}

func allocateToMaleHostels(mhs *employeeSlice, maleHostels []string, unAllocatedToMaleHostels chan<- interface{}, wg *sync.WaitGroup) {
	var file []byte
	var allocationSlice []Space

	for _, maleHostel := range maleHostels {
		spc := &Space{
			Name:       maleHostel,
			MaxPersons: 4,
			SpaceType:  "maleRoom",
		}

		spc.addMembers(&mhs)

		// append allocation to slice
		allocationSlice = append(allocationSlice, *spc)
	}

	// write allocation to json file
	file, _ = json.MarshalIndent(allocationSlice, "", " ")
	_ = ioutil.WriteFile("maleHostelAllocation.json", file, 0644)

	unAllocatedToMaleHostels <- mhs
	wg.Done()
}

func allocateToFemaleHostels(fhs *employeeSlice, femaleHostels []string, unAllocatedToFemaleHostels chan<- interface{}, wg *sync.WaitGroup) {

	var file []byte
	var allocationSlice []Space

	for _, femaleHostel := range femaleHostels {
		spc := &Space{
			Name:       femaleHostel,
			MaxPersons: 4,
			SpaceType:  "femaleRoom",
		}

		spc.addMembers(&fhs)

		// append allocation to slice
		allocationSlice = append(allocationSlice, *spc)
	}

	// write allocation to json file
	file, _ = json.MarshalIndent(allocationSlice, "", " ")
	_ = ioutil.WriteFile("femaleHostelAllocation.json", file, 0644)

	unAllocatedToFemaleHostels <- fhs
	wg.Done()
}

func getUnallocatedemployees(officeSpace <-chan interface{}, maleHostels <-chan interface{}, femaleHostels <-chan interface{}, wg *sync.WaitGroup) {
	fmt.Println("Waiting to recieve updates...")
	officeLeftOvers := <-officeSpace
	maleHostelsLeftOvers := <-maleHostels
	femaleHostelsLeftOvers := <-femaleHostels

	fmt.Println(officeLeftOvers, maleHostelsLeftOvers, femaleHostelsLeftOvers, "There are here gentlefella")
	wg.Done()

}

func main() {
	inputfile := &fileParser{filepath: "inputA.txt"}
	e := inputfile.GetEmployees()
	var eSlice, maleHostelSlice, femaleHostelSlice employeeSlice
	for _, val := range e {
		eSlice = append(eSlice, val...)
	}

	maleHostelSlice, _ = e["maleFellows"]
	femaleHostelSlice, _ = e["femaleFellows"]

	// shuffle employee slice for random office allocation
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(eSlice), func(i, j int) { eSlice[i], eSlice[j] = eSlice[j], eSlice[i] })

	// shuffle maleHostelSlice slice for random male hostel allocation
	rand.Shuffle(len(maleHostelSlice), func(i, j int) {
		maleHostelSlice[i], maleHostelSlice[j] = maleHostelSlice[j], maleHostelSlice[i]
	})

	// shuffle femaleHostelSlice slice for random female hostel allocation
	rand.Shuffle(len(femaleHostelSlice), func(i, j int) {
		femaleHostelSlice[i], femaleHostelSlice[j] = femaleHostelSlice[j], femaleHostelSlice[i]
	})

	// declare the hostels and office here
	femaleHostel := []string{"ruby", "platinum", "jade", "pearl", "diamond"}
	maleHostel := []string{"topaz", "silver", "gold", "onyx", "opal"}
	office := []string{
		"Carat", "Anvil", "Crucible",
		"Kiln", "Forge", "Foundry",
		"Furnace", "Boiler",
		"Mint", "Vulcan",
	}

	// define channels here
	unAllocatedToOffice := make(chan interface{})
	unAllocatedToMaleHostels := make(chan interface{})
	unAllocatedToFemaleHostels := make(chan interface{})
	var wg sync.WaitGroup // add a wait group here

	wg.Add(4)

	go allocateToOffice(&eSlice, office, unAllocatedToOffice, &wg)
	go allocateToMaleHostels(&maleHostelSlice, maleHostel, unAllocatedToMaleHostels, &wg)
	go allocateToFemaleHostels(&femaleHostelSlice, femaleHostel, unAllocatedToFemaleHostels, &wg)

	go getUnallocatedemployees(unAllocatedToOffice, unAllocatedToMaleHostels, unAllocatedToFemaleHostels, &wg)
	wg.Wait()
}
