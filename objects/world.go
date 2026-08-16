package objects

import (
	"math/rand/v2"
)

type World struct {
	Networks    map[string]PublicNetwork
	Machines    map[int]Machine
	Filesystems map[int]*Filesystem // keyed by MachineID
}

type Machine struct {
	MachineID   int    // id of filesystem in database
	PublicIp    string // public ip of machine
	LocalID     int    // local id of machine within its public network
	MachineName string // name of the computer obviously
}

func (world *World) NewMachine(publicIp string, machineName string) Machine {
	id := int(rand.Uint32())

	// list it in the public network
	if _, exists := world.Networks[publicIp]; !exists {
		world.Networks[publicIp] = PublicNetwork{
			PublicIp:          publicIp,
			MachinesByLocalID: make(map[int]int),
			ForwardedPorts:    make(map[int]int),
			LocalIDCounter:    0,
		}
	}

	// increment the local id counter in the public network for generating a new local id
	pubnet := world.Networks[publicIp]
	pubnet.LocalIDCounter += 1
	world.Networks[publicIp] = pubnet

	localID := pubnet.LocalIDCounter

	machineObj := Machine{
		MachineID:   id,
		PublicIp:    publicIp,
		LocalID:     localID,
		MachineName: machineName,
	}
	world.Machines[id] = machineObj

	// return new object
	return machineObj
}

func (world *World) GetMachine(machineID int) (Machine, bool) {
	machine, exists := world.Machines[machineID]
	return machine, exists
}
