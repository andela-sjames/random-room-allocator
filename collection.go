package main

// NotAllocatedSlice type defined
type NotAllocatedSlice map[string][]Employee

// Employee struct defined
type Employee struct {
	Name        interface{} `json:"Name"`
	Gender      interface{} `json:"Gender"`
	Position    interface{} `json:"Position"`
	LivingSpace interface{} `json:"LivingSpace"`
}

// Space struct defined
type Space struct {
	Name       string        `json:"Name"`
	MaxPersons int           `json:"MaxPersons"`
	Members    EmployeeSlice `json:"Members"`
	SpaceType  string        `json:"Type"`
}

func (sp *Space) addOfficeMembers(e **EmployeeSlice) {
	var lastItem EmployeeDataMap
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

func (sp *Space) addMaleMembers(e **EmployeeSlice) {
	var lastItem EmployeeDataMap
	for i := 0; i < sp.MaxPersons; i++ {
		if len(**e) > 0 {
			// read the last item from the slice
			lastItem = (**e)[len(**e)-1]

			// add only if member wants living space
			if lastItem["livingSpace"] == true {
				sp.Members = append(sp.Members, lastItem)
			}
			// remove the last item of the slice
			**e = (**e)[:len(**e)-1]
		}
	}
}

func (sp *Space) addFemaleMembers(e **EmployeeSlice) {
	var lastItem EmployeeDataMap
	for i := 0; i < sp.MaxPersons; i++ {
		if len(**e) > 0 {
			// read the last item from the slice
			lastItem = (**e)[len(**e)-1]

			// add only if member wants living space
			if lastItem["livingSpace"] == true {
				sp.Members = append(sp.Members, lastItem)
			}
			// remove the last item of the slice
			**e = (**e)[:len(**e)-1]
		}
	}
}
