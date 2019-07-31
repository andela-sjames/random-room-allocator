package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
)

func allocateToOffice(e *EmployeeSlice, offices []string, unAllocatedToOffice chan<- EmployeeSlice, wg *sync.WaitGroup) {
	var file []byte
	var allocationSlice []Space

	for _, office := range offices {
		spc := &Space{
			Name:       office,
			MaxPersons: 6,
			SpaceType:  "officeRoom",
		}

		spc.addOfficeMembers(&e)

		// append allocation to slice
		allocationSlice = append(allocationSlice, *spc)
	}

	// write allocation to json file
	file, _ = json.MarshalIndent(allocationSlice, "", " ")
	_ = ioutil.WriteFile("officeAllocation.json", file, 0644)

	fmt.Println(`Office allocation file "officeAllocation.json" created :)`)

	unAllocatedToOffice <- *e
	wg.Done()
}

func allocateToMaleHostels(mhs *EmployeeSlice, maleHostels []string, unAllocatedToMaleHostels chan<- EmployeeSlice, wg *sync.WaitGroup) {
	var file []byte
	var allocationSlice []Space

	for _, maleHostel := range maleHostels {
		spc := &Space{
			Name:       maleHostel,
			MaxPersons: 4,
			SpaceType:  "maleRoom",
		}

		spc.addMaleMembers(&mhs)

		// append allocation to slice
		allocationSlice = append(allocationSlice, *spc)
	}

	// write allocation to json file
	file, _ = json.MarshalIndent(allocationSlice, "", " ")
	_ = ioutil.WriteFile("maleHostelAllocation.json", file, 0644)

	fmt.Println(`Male allocation file "maleHostelAllocation.json" created :)`)

	unAllocatedToMaleHostels <- *mhs
	wg.Done()
}

func allocateToFemaleHostels(fhs *EmployeeSlice, femaleHostels []string, unAllocatedToFemaleHostels chan<- EmployeeSlice, wg *sync.WaitGroup) {

	var file []byte
	var allocationSlice []Space

	for _, femaleHostel := range femaleHostels {
		spc := &Space{
			Name:       femaleHostel,
			MaxPersons: 4,
			SpaceType:  "femaleRoom",
		}

		spc.addFemaleMembers(&fhs)

		// append allocation to slice
		allocationSlice = append(allocationSlice, *spc)
	}

	// write allocation to json file
	file, _ = json.MarshalIndent(allocationSlice, "", " ")
	_ = ioutil.WriteFile("femaleHostelAllocation.json", file, 0644)

	fmt.Println(`Female allocation file "femaleHostelAllocation.json" created :)`)

	unAllocatedToFemaleHostels <- *fhs
	wg.Done()
}

func getUnallocatedemployees(officeSpace <-chan EmployeeSlice, maleHostels <-chan EmployeeSlice, femaleHostels <-chan EmployeeSlice, wg *sync.WaitGroup) {
	fmt.Println("Waiting to recieve no allocation updates...")

	var file []byte

	officeLeftOvers := <-officeSpace
	maleHostelsLeftOvers := <-maleHostels
	femaleHostelsLeftOvers := <-femaleHostels

	employees := make(NotAllocatedSlice)
	var resultSlice []Employee

	if len(officeLeftOvers) > 0 {
		// office leftovers
		for _, item := range officeLeftOvers {
			emp := &Employee{
				Name:     item["name"],
				Gender:   item["gender"],
				Position: item["position"],
			}
			resultSlice = append(resultSlice, *emp)
		}

		employees["staff"] = resultSlice
	}

	if len(maleHostelsLeftOvers) > 0 {
		// maleHostelsLeftOvers leftovers
		for _, item := range maleHostelsLeftOvers {
			emp := &Employee{
				Name:        item["name"],
				Gender:      item["gender"],
				Position:    item["position"],
				LivingSpace: item["livingSpace"],
			}
			resultSlice = append(resultSlice, *emp)
		}

		employees["maleFellows"] = resultSlice
	}

	if len(femaleHostelsLeftOvers) > 0 {
		for _, item := range maleHostelsLeftOvers {
			emp := &Employee{
				Name:        item["name"],
				Gender:      item["gender"],
				Position:    item["position"],
				LivingSpace: item["livingSpace"],
			}
			resultSlice = append(resultSlice, *emp)
		}

		employees["femaleFellows"] = resultSlice
	}

	// write allocation to json file
	file, _ = json.MarshalIndent(employees, "", " ")
	_ = ioutil.WriteFile("NoAllocation.json", file, 0644)

	fmt.Println(`No allocation file "NoAllocation.json" created :)`)
	wg.Done()
}

func main() {
	inputfile := &FileParser{Filepath: "inputA.txt"}
	e := inputfile.GetEmployees()
	var eSlice, maleHostelSlice, femaleHostelSlice EmployeeSlice
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

	// declare the hostels and offices here
	femaleHostel := []string{"ruby", "platinum", "jade", "pearl", "diamond"}
	maleHostel := []string{"topaz", "silver", "gold", "onyx", "opal"}
	office := []string{
		"Carat", "Anvil", "Crucible",
		"Kiln", "Forge", "Foundry",
		"Furnace", "Boiler",
		"Mint", "Vulcan",
	}

	// define channels here
	unAllocatedToOffice := make(chan EmployeeSlice)
	unAllocatedToMaleHostels := make(chan EmployeeSlice)
	unAllocatedToFemaleHostels := make(chan EmployeeSlice)

	// add a wait group here
	var wg sync.WaitGroup
	wg.Add(4)

	go allocateToOffice(&eSlice, office, unAllocatedToOffice, &wg)
	go allocateToMaleHostels(&maleHostelSlice, maleHostel, unAllocatedToMaleHostels, &wg)
	go allocateToFemaleHostels(&femaleHostelSlice, femaleHostel, unAllocatedToFemaleHostels, &wg)
	go getUnallocatedemployees(unAllocatedToOffice, unAllocatedToMaleHostels, unAllocatedToFemaleHostels, &wg)

	wg.Wait()
	fmt.Println("Done!")
}
