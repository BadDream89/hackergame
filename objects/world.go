package objects

import (
	"fmt"
	"math/rand/v2"
)

type World struct {
	Networks map[string]PublicNetwork
}

type MachinesStorage struct {
	Machines map[int]Machine // store machines by their id
}

type Machine struct {
	MachineID int    // id of filesystem in database
	PublicIp  string // public ip of machine
	LocalIp   string // local ip of machine
}

func (world *World) NewMachine(publicIp string) Machine {
	id := int(rand.Uint32())

	if _, exists := world.Networks[publicIp]; !exists {
		world.Networks[publicIp] = PublicNetwork{
			Ip:             publicIp,
			LocalNet:       make(map[string]int),
			ForwardedPorts: make(map[int]string),
			Machines:       1,
		}
	}

	pubnet := world.Networks[publicIp]
	pubnet.Machines += 1
	world.Networks[publicIp] = pubnet

	localIp := fmt.Sprintf("192.168.1.%d", pubnet.Machines)

	return Machine{
		MachineID: id,
		PublicIp:  publicIp,
		LocalIp:   localIp,
	}
}
