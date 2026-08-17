package objects

import "errors"

// do not use for creating new connection handler
type ConnectionHandler struct {
	CurrentMachineID int   // id of current machine
	ConnectionPath   []int // path of connecting through the machines
}

func NewConnectionHandler(world *World, hostMachineID int) (ConnectionHandler, error) {
	if _, exists := world.GetMachine(hostMachineID); !exists {
		return ConnectionHandler{}, errors.New("machine does not exist")
	}

	return ConnectionHandler{
		CurrentMachineID: hostMachineID,
		ConnectionPath:   []int{hostMachineID},
	}, nil
}

func (conn *ConnectionHandler) Connect(world *World, machineID int) error {
	if _, exists := world.GetMachine(machineID); !exists {
		return errors.New("machine does not exist")
	}

	conn.ConnectionPath = append(conn.ConnectionPath, machineID)
	conn.CurrentMachineID = machineID

	return nil
}

func (conn *ConnectionHandler) Disconnect() error {
	connPathLen := len(conn.ConnectionPath)
	if connPathLen < 2 {
		return errors.New("cannot disconnect from the main machine")
	}

	conn.CurrentMachineID = conn.ConnectionPath[connPathLen-2]
	conn.ConnectionPath = conn.ConnectionPath[:connPathLen-1]

	return nil
}

func (conn *ConnectionHandler) DisconnectFromAll() error {
	conn.CurrentMachineID = conn.ConnectionPath[0]
	conn.ConnectionPath = conn.ConnectionPath[:1]

	return nil
}
