package objects

import "errors"

type PublicNetwork struct {
	PublicIp          string      // publicIp of network
	MachinesByLocalID map[int]int // links machine object to a localID | localID: machineId
	ForwardedPorts    map[int]int // port: local id
	LocalIDCounter    int         // last local id assigned to a machine on this network
}

// functions
func (network *PublicNetwork) ForwardPort(port int, localID int) (int, error) {
	for p := range network.ForwardedPorts {
		if p == port {
			return network.ForwardedPorts[p], errors.New("Port is already forwarded.")
		}
	}

	network.ForwardedPorts[port] = localID

	return localID, nil
}

func (network *PublicNetwork) AddMachine(localID int, machine *Machine) (int, error) {
	for m := range network.MachinesByLocalID {
		if m == localID {
			return machine.LocalID, errors.New("Machine with this ip already exists.")
		}
	}

	network.MachinesByLocalID[localID] = machine.MachineID

	return localID, nil
}
