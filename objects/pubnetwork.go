package objects

import "errors"

type PublicNetwork struct {
	Ip             string         // publicIp of network
	LocalNet       map[string]int // links machine object to a localIp | localIp: machineId
	ForwardedPorts map[int]string // port: local ip
	Machines       int
}

// functions
func (network *PublicNetwork) ForwardPort(port int, localIp string) (string, error) {
	for p := range network.ForwardedPorts {
		if p == port {
			return network.ForwardedPorts[p], errors.New("Port is already forwarded.")
		}
	}

	network.ForwardedPorts[port] = localIp

	return localIp, nil
}

func (network *PublicNetwork) AddMachine(localIp string, machine *Machine) (string, error) {
	for m := range network.LocalNet {
		if m == localIp {
			return machine.LocalIp, errors.New("Machine with this ip already exists.")
		}
	}

	network.LocalNet[localIp] = machine.MachineID

	return localIp, nil
}
