package tests

import (
	"testing"

	"hackergame/objects"
)

// newHandlerChain creates a host machine plus n-1 additional hops, builds a
// ConnectionHandler via NewConnectionHandler rooted at the host, connects
// through the remaining hops in order, and returns the handler along with
// the machine IDs (ids[0] is the host/main machine).
func newHandlerChain(t *testing.T, world *objects.World, n int) (*objects.ConnectionHandler, []int) {
	t.Helper()

	if n < 1 {
		t.Fatalf("newHandlerChain requires n >= 1, got %d", n)
	}

	host := world.NewMachine("203.0.113.1", "host")
	handler, err := objects.NewConnectionHandler(world, host.MachineID)
	if err != nil {
		t.Fatalf("unexpected error creating connection handler: %v", err)
	}

	ids := []int{host.MachineID}
	for i := 1; i < n; i++ {
		machine := world.NewMachine("203.0.113.1", "hop")
		if err := handler.Connect(world, machine.MachineID); err != nil {
			t.Fatalf("unexpected error connecting hop %d: %v", i, err)
		}
		ids = append(ids, machine.MachineID)
	}

	return &handler, ids
}

func TestNewConnectionHandler_Success(t *testing.T) {
	world := newTestWorld()
	host := world.NewMachine("203.0.113.1", "host")

	handler, err := objects.NewConnectionHandler(world, host.MachineID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if handler.CurrentMachineID != host.MachineID {
		t.Errorf("expected CurrentMachineID %d, got %d", host.MachineID, handler.CurrentMachineID)
	}
	if len(handler.ConnectionPath) != 1 || handler.ConnectionPath[0] != host.MachineID {
		t.Errorf("expected ConnectionPath [%d], got %v", host.MachineID, handler.ConnectionPath)
	}
}

func TestNewConnectionHandler_MachineDoesNotExist(t *testing.T) {
	world := newTestWorld()

	handler, err := objects.NewConnectionHandler(world, 999)
	if err == nil {
		t.Fatal("expected error when creating a handler for a nonexistent machine, got nil")
	}
	if handler.CurrentMachineID != 0 {
		t.Errorf("expected zero-value CurrentMachineID, got %d", handler.CurrentMachineID)
	}
	if len(handler.ConnectionPath) != 0 {
		t.Errorf("expected empty ConnectionPath, got %v", handler.ConnectionPath)
	}
}

func TestConnect_Success(t *testing.T) {
	world := newTestWorld()
	handler, ids := newHandlerChain(t, world, 1)
	target := world.NewMachine("198.51.100.1", "target-pc")

	err := handler.Connect(world, target.MachineID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if handler.CurrentMachineID != target.MachineID {
		t.Errorf("expected CurrentMachineID %d, got %d", target.MachineID, handler.CurrentMachineID)
	}
	want := []int{ids[0], target.MachineID}
	if len(handler.ConnectionPath) != 2 || handler.ConnectionPath[0] != want[0] || handler.ConnectionPath[1] != want[1] {
		t.Errorf("expected ConnectionPath %v, got %v", want, handler.ConnectionPath)
	}
}

func TestConnect_MachineDoesNotExist(t *testing.T) {
	world := newTestWorld()
	handler, ids := newHandlerChain(t, world, 1)

	err := handler.Connect(world, 999)
	if err == nil {
		t.Fatal("expected error when connecting to a nonexistent machine, got nil")
	}
	if handler.CurrentMachineID != ids[0] {
		t.Errorf("expected CurrentMachineID to remain %d, got %d", ids[0], handler.CurrentMachineID)
	}
	if len(handler.ConnectionPath) != 1 {
		t.Errorf("expected ConnectionPath to remain [%d], got %v", ids[0], handler.ConnectionPath)
	}
}

func TestConnect_BuildsPathAcrossMultipleHops(t *testing.T) {
	world := newTestWorld()
	handler, ids := newHandlerChain(t, world, 3)

	if handler.CurrentMachineID != ids[2] {
		t.Errorf("expected CurrentMachineID %d, got %d", ids[2], handler.CurrentMachineID)
	}
	if len(handler.ConnectionPath) != 3 {
		t.Fatalf("expected ConnectionPath of length 3, got %v", handler.ConnectionPath)
	}
	for i, id := range ids {
		if handler.ConnectionPath[i] != id {
			t.Errorf("expected ConnectionPath[%d] = %d, got %d", i, id, handler.ConnectionPath[i])
		}
	}
}

func TestConnect_FailedHopDoesNotBreakSubsequentConnect(t *testing.T) {
	world := newTestWorld()
	handler, ids := newHandlerChain(t, world, 1)
	target := world.NewMachine("198.51.100.1", "target-pc")

	if err := handler.Connect(world, 999); err == nil {
		t.Fatal("expected error connecting to nonexistent machine, got nil")
	}
	if err := handler.Connect(world, target.MachineID); err != nil {
		t.Fatalf("expected no error on valid connect after a failed one, got %v", err)
	}

	if handler.CurrentMachineID != target.MachineID {
		t.Errorf("expected CurrentMachineID %d, got %d", target.MachineID, handler.CurrentMachineID)
	}
	if len(handler.ConnectionPath) != 2 || handler.ConnectionPath[0] != ids[0] || handler.ConnectionPath[1] != target.MachineID {
		t.Errorf("expected ConnectionPath [%d %d], got %v", ids[0], target.MachineID, handler.ConnectionPath)
	}
}

func TestDisconnect_Success(t *testing.T) {
	world := newTestWorld()
	handler, ids := newHandlerChain(t, world, 2)

	err := handler.Disconnect()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if handler.CurrentMachineID != ids[0] {
		t.Errorf("expected CurrentMachineID to revert to %d, got %d", ids[0], handler.CurrentMachineID)
	}
	if len(handler.ConnectionPath) != 1 || handler.ConnectionPath[0] != ids[0] {
		t.Errorf("expected ConnectionPath [%d], got %v", ids[0], handler.ConnectionPath)
	}
}

func TestDisconnect_MultipleTimesWalksPathBackwards(t *testing.T) {
	world := newTestWorld()
	handler, ids := newHandlerChain(t, world, 3)

	if err := handler.Disconnect(); err != nil {
		t.Fatalf("unexpected error on first disconnect: %v", err)
	}
	if handler.CurrentMachineID != ids[1] {
		t.Errorf("expected CurrentMachineID %d after first disconnect, got %d", ids[1], handler.CurrentMachineID)
	}

	if err := handler.Disconnect(); err != nil {
		t.Fatalf("unexpected error on second disconnect: %v", err)
	}
	if handler.CurrentMachineID != ids[0] {
		t.Errorf("expected CurrentMachineID %d after second disconnect, got %d", ids[0], handler.CurrentMachineID)
	}
	if len(handler.ConnectionPath) != 1 {
		t.Errorf("expected ConnectionPath of length 1, got %v", handler.ConnectionPath)
	}
}

func TestDisconnect_FromMainMachineErrors(t *testing.T) {
	world := newTestWorld()
	handler, ids := newHandlerChain(t, world, 1)

	err := handler.Disconnect()
	if err == nil {
		t.Fatal("expected error when disconnecting from the main machine, got nil")
	}
	if handler.CurrentMachineID != ids[0] {
		t.Errorf("expected CurrentMachineID to remain %d, got %d", ids[0], handler.CurrentMachineID)
	}
	if len(handler.ConnectionPath) != 1 {
		t.Errorf("expected ConnectionPath to remain unchanged, got %v", handler.ConnectionPath)
	}
}

func TestDisconnect_OnZeroValueHandlerErrors(t *testing.T) {
	handler := &objects.ConnectionHandler{}

	err := handler.Disconnect()
	if err == nil {
		t.Fatal("expected error when disconnecting a zero-value handler with no connection path, got nil")
	}
}

func TestDisconnectFromAll_Success(t *testing.T) {
	world := newTestWorld()
	handler, ids := newHandlerChain(t, world, 3)

	err := handler.DisconnectFromAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if handler.CurrentMachineID != ids[0] {
		t.Errorf("expected CurrentMachineID to revert to main machine %d, got %d", ids[0], handler.CurrentMachineID)
	}
	if len(handler.ConnectionPath) != 1 || handler.ConnectionPath[0] != ids[0] {
		t.Errorf("expected ConnectionPath [%d], got %v", ids[0], handler.ConnectionPath)
	}
}

func TestDisconnectFromAll_AlreadyAtMainMachine(t *testing.T) {
	world := newTestWorld()
	handler, ids := newHandlerChain(t, world, 1)

	err := handler.DisconnectFromAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if handler.CurrentMachineID != ids[0] {
		t.Errorf("expected CurrentMachineID to remain %d, got %d", ids[0], handler.CurrentMachineID)
	}
	if len(handler.ConnectionPath) != 1 {
		t.Errorf("expected ConnectionPath to remain length 1, got %v", handler.ConnectionPath)
	}
}

// This documents current behavior rather than desired behavior: DisconnectFromAll
// indexes ConnectionPath[0] without checking length first. NewConnectionHandler always
// seeds ConnectionPath with the host machine, so this only bites a caller who builds a
// ConnectionHandler{} directly instead of going through the constructor.
func TestDisconnectFromAll_PanicsOnZeroValueHandler(t *testing.T) {
	handler := &objects.ConnectionHandler{}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected DisconnectFromAll to panic on an empty ConnectionPath (documents current unguarded behavior)")
		}
	}()

	_ = handler.DisconnectFromAll()
}
