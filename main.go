package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ParkingLot struct {
	capacity          int
	slots             map[int]string // slot number -> registration number
	nextAvailableSlot int
}

func NewParkingLot() *ParkingLot {
	return &ParkingLot{
		capacity:          0,
		slots:             make(map[int]string),
		nextAvailableSlot: 1,
	}
}

func (p *ParkingLot) CreateParkingLot(capacity int) {
	p.capacity = capacity
	p.slots = make(map[int]string)
	p.nextAvailableSlot = 1
	fmt.Printf("Created a parking lot with %d slots\n", capacity)
}

func (p *ParkingLot) Park(registrationNumber string) {
	if len(p.slots) >= p.capacity {
		fmt.Println("Sorry, parking lot is full")
		return
	}

	// Find the next available slot
	for p.slots[p.nextAvailableSlot] != "" {
		p.nextAvailableSlot++
	}

	p.slots[p.nextAvailableSlot] = registrationNumber
	fmt.Printf("Allocated slot number: %d\n", p.nextAvailableSlot)
}

func (p *ParkingLot) Leave(registrationNumber string, hours int) {
	slotNumber := 0
	for slot, reg := range p.slots {
		if reg == registrationNumber {
			slotNumber = slot
			break
		}
	}

	if slotNumber == 0 {
		fmt.Printf("Registration number %s not found\n", registrationNumber)
		return
	}

	// Calculate parking charge
	charge := 10
	if hours > 2 {
		charge += (hours - 2) * 10
	}

	// Free the slot
	delete(p.slots, slotNumber)
	// Update nextAvailableSlot if necessary
	if slotNumber < p.nextAvailableSlot {
		p.nextAvailableSlot = slotNumber
	}

	fmt.Printf("Registration number %s with Slot Number %d is free with Charge $%d\n",
		registrationNumber, slotNumber, charge)
}

func (p *ParkingLot) Status() {
	fmt.Println("Slot No.    Registration No.")
	// Create a slice of slot numbers for sorting
	slots := make([]int, 0, len(p.slots))
	for slot := range p.slots {
		slots = append(slots, slot)
	}
	// Sort slots
	for i := 0; i < len(slots)-1; i++ {
		for j := i + 1; j < len(slots); j++ {
			if slots[i] > slots[j] {
				slots[i], slots[j] = slots[j], slots[i]
			}
		}
	}
	// Print sorted slots
	for _, slot := range slots {
		fmt.Printf("%d           %s\n", slot, p.slots[slot])
	}
}

func processCommands(filename string) {
	parkingLot := NewParkingLot()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: File %s not found\n", filename)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		switch parts[0] {
		case "create_parking_lot":
			capacity, _ := strconv.Atoi(parts[1])
			parkingLot.CreateParkingLot(capacity)

		case "park":
			parkingLot.Park(parts[1])

		case "leave":
			hours, _ := strconv.Atoi(parts[2])
			parkingLot.Leave(parts[1], hours)

		case "status":
			parkingLot.Status()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error processing file: %v\n", err)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}
	processCommands(os.Args[1])
}
