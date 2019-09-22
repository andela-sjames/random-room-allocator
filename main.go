package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/andela-sjames/random-room-allocator/allocator"
)

func main() {
	inputfile := &allocator.FileParser{Filepath: "inputA.txt"}
	e := inputfile.GetEmployees()
	var eSlice, maleHostelSlice, femaleHostelSlice allocator.EmployeeSlice
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
	unAllocatedToOffice := make(chan allocator.EmployeeSlice)
	unAllocatedToMaleHostels := make(chan allocator.EmployeeSlice)
	unAllocatedToFemaleHostels := make(chan allocator.EmployeeSlice)

	// add a wait group here
	var wg sync.WaitGroup
	wg.Add(4)

	go allocator.AllocateToOffice(&eSlice, office, unAllocatedToOffice, &wg)
	go allocator.AllocateToMaleHostels(&maleHostelSlice, maleHostel, unAllocatedToMaleHostels, &wg)
	go allocator.AllocateToFemaleHostels(&femaleHostelSlice, femaleHostel, unAllocatedToFemaleHostels, &wg)
	go allocator.GetUnallocatedemployees(unAllocatedToOffice, unAllocatedToMaleHostels, unAllocatedToFemaleHostels, &wg)

	wg.Wait()
	fmt.Println("Done!")
}
