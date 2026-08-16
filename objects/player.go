package objects

type Player struct {
	MainMachine int // id of main machine
	Nickname    string
}

type Connection struct {
	CurrentMachine int   // id of current machine
	ConnectionPath []int // path of connecting through the machines
}
