package objects

import (
	"errors"
	"math/rand/v2"
)

type World struct {
	Networks    map[string]PublicNetwork
	Machines    map[int]Machine
	Filesystems map[int]*Filesystem // keyed by MachineID
	Players     map[string]Player   // nickname: player
}

type Machine struct {
	MachineID   int    // id of filesystem in database
	PublicIp    string // public ip of machine
	LocalID     int    // local id of machine within its public network
	MachineName string // name of the computer obviously
}

type Player struct {
	MainMachineID int // id of main machine
	Nickname      string
}

// creates a new machine object
// adds it in the World.Machines
// lists on the PublicNetwork
func (world *World) NewMachine(publicIp string, machineName string) Machine {
	id := int(rand.Uint32())

	if _, exists := world.Networks[publicIp]; !exists {
		world.Networks[publicIp] = PublicNetwork{
			PublicIp:          publicIp,
			MachinesByLocalID: make(map[int]int),
			ForwardedPorts:    make(map[int]int),
			LocalIDCounter:    0,
		}
	}

	// list it in the public network
	// increment the local id counter in the public network for generating a new local id
	pubnet := world.Networks[publicIp]
	pubnet.LocalIDCounter += 1
	localID := pubnet.LocalIDCounter
	pubnet.MachinesByLocalID[localID] = id

	// update public network in the world
	world.Networks[publicIp] = pubnet

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

func (world *World) NewPlayer(nickname string, machineName string) error {
	if _, exists := world.Players[nickname]; exists {
		return errors.New("that player already exists")
	}

	playerMachine := world.NewMachine("", machineName)
	player := Player{
		MainMachineID: playerMachine.MachineID,
		Nickname:      nickname,
	}

	world.Players[nickname] = player

	return nil
}
