package objects

import (
	"math/rand/v2"
)

type World struct {
	Networks map[string]PublicNetwork
}

type MachinesStorage struct {
	Machines map[int]Machine // store machines by their id
}

type Machine struct {
	MachineID   int    // id of filesystem in database
	PublicIp    string // public ip of machine
	LocalID     int    // local ip of machine
	MachineName string // name of the computer obviously
}

func (world *World) NewMachine(publicIp string, machineName string) Machine {
	id := int(rand.Uint32())

	// list it in the public network
	if _, exists := world.Networks[publicIp]; !exists {
		world.Networks[publicIp] = PublicNetwork{
			Ip:             publicIp,
			LocalNet:       make(map[int]int),
			ForwardedPorts: make(map[int]int),
			Machines:       0,
		}
	}

	// increment the machines count in the public network for generating a new local id
	pubnet := world.Networks[publicIp]
	pubnet.Machines += 1
	world.Networks[publicIp] = pubnet

	localID := pubnet.Machines

	// return new object
	return Machine{
		MachineID:   id,
		PublicIp:    publicIp,
		LocalID:     localID,
		MachineName: machineName,
	}
}
