package allocator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

func AllocateToOffice(e *EmployeeSlice, offices []string, unAllocatedToOffice chan<- EmployeeSlice, wg *sync.WaitGroup) {
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

func AllocateToMaleHostels(mhs *EmployeeSlice, maleHostels []string, unAllocatedToMaleHostels chan<- EmployeeSlice, wg *sync.WaitGroup) {
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

func AllocateToFemaleHostels(fhs *EmployeeSlice, femaleHostels []string, unAllocatedToFemaleHostels chan<- EmployeeSlice, wg *sync.WaitGroup) {

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

func GetUnallocatedemployees(officeSpace <-chan EmployeeSlice, maleHostels <-chan EmployeeSlice, femaleHostels <-chan EmployeeSlice, wg *sync.WaitGroup) {
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
