package objects

import "errors"

type PublicNetwork struct {
	Ip             string      // publicIp of network
	LocalNet       map[int]int // links machine object to a localID | localID: machineId
	ForwardedPorts map[int]int // port: local id
	Machines       int
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
	for m := range network.LocalNet {
		if m == localID {
			return machine.LocalID, errors.New("Machine with this ip already exists.")
		}
	}

	network.LocalNet[localID] = machine.MachineID

	return localID, nil
}
