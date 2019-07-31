# random-room-allocator

A random room allocation system in go.

## Algorithm to randomly allocate company members to randomly generated rooms

### Some context here

- Store office space allocation to a json file
- Store male hostel allocation to a json file
- Store female hostel allocations to ja son file
- Store unallocated employee data to a json file


The code has a pre-defined number of offices for allocation and a predefined number of hostels to allocate to male and female fellows who opted for it.

```Text
Offices prepopulated are:'Carat', 'Anvil', 'Crucible', 'Kiln', 'Forge', 'Foundry', 'Furnace', 'Boiler', 'Mint', 'Vulcan'
Malerooms prepopulated are:'topaz', 'silver', 'gold', 'onyx', 'opal'
Femalerooms prepopulated are: 'ruby', 'platinum', 'jade', 'pearl', 'diamond'
```


```Text
Basic Conditions:
1. No staff should be allocated to Male or Female Rooms
2. No Male or Female room should exceed 4 persons
3. No office allocation should exceed 6 persons
```

### Sample input file structure

```Text
BOLA AHMED   M FELLOW Y
JOHN OBI     M FELLOW N
ISSAC NNADI  M STAFF
CRIBS JANE   F FELLOW Y
```

### build go file

go build go run {fileparser,main,collection,allocation}.go

### run go code

go run {fileparser,main,collection,allocation}.go
