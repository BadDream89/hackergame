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
	if existing, exists := network.ForwardedPorts[port]; exists {
		return existing, errors.New("Port is already forwarded.")
	}

	network.ForwardedPorts[port] = localID

	return localID, nil
}

func (network *PublicNetwork) UnforwardPort(port int) error {
	if _, exists := network.ForwardedPorts[port]; !exists {
		return errors.New("that port is not forwarded")
	}

	delete(network.ForwardedPorts, port)

	return nil
}
