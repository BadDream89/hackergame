package tests

import (
	"testing"

	"hackergame/objects"
)

func newTestWorld() *objects.World {
	return &objects.World{Networks: make(map[string]objects.PublicNetwork)}
}

func TestNewMachine_CreatesNetworkIfNotExists(t *testing.T) {
	world := newTestWorld()

	machine := world.NewMachine("203.0.113.1", "test-pc")

	if machine.PublicIp != "203.0.113.1" {
		t.Errorf("expected PublicIp %q, got %q", "203.0.113.1", machine.PublicIp)
	}
	if machine.MachineName != "test-pc" {
		t.Errorf("expected MachineName %q, got %q", "test-pc", machine.MachineName)
	}
	if _, exists := world.Networks["203.0.113.1"]; !exists {
		t.Fatal("expected a network to be created for a new publicIp")
	}
	if machine.LocalID != 1 {
		t.Errorf("expected first machine's LocalID %d, got %d", 1, machine.LocalID)
	}
}

func TestNewMachine_ReturnsNonZeroMachineID(t *testing.T) {
	world := newTestWorld()

	machine := world.NewMachine("203.0.113.1", "test-pc")

	// rand.Uint32() can technically be 0, but astronomically unlikely; this guards against
	// an obviously broken id generator (e.g. one that always returns 0).
	if machine.MachineID == 0 {
		t.Error("expected a non-zero MachineID")
	}
}

func TestNewMachine_ExistingNetworkReusesIt(t *testing.T) {
	world := newTestWorld()
	world.Networks["203.0.113.1"] = objects.PublicNetwork{
		Ip:             "203.0.113.1",
		LocalNet:       make(map[int]int),
		ForwardedPorts: make(map[int]string),
		Machines:       5,
	}

	machine := world.NewMachine("203.0.113.1", "test-pc")

	if machine.PublicIp != "203.0.113.1" {
		t.Errorf("expected PublicIp %q, got %q", "203.0.113.1", machine.PublicIp)
	}
	if machine.LocalID != 6 {
		t.Errorf("expected LocalID derived from existing Machines count %d, got %d", 6, machine.LocalID)
	}
}

func TestNewMachine_PersistsMachinesCountAcrossCalls(t *testing.T) {
	world := newTestWorld()

	first := world.NewMachine("203.0.113.1", "pc-1")
	second := world.NewMachine("203.0.113.1", "pc-2")
	third := world.NewMachine("203.0.113.1", "pc-3")

	if first.LocalID != 1 {
		t.Errorf("expected first machine LocalID %d, got %d", 1, first.LocalID)
	}
	if second.LocalID != 2 {
		t.Errorf("expected second machine LocalID %d, got %d", 2, second.LocalID)
	}
	if third.LocalID != 3 {
		t.Errorf("expected third machine LocalID %d, got %d", 3, third.LocalID)
	}
	if world.Networks["203.0.113.1"].Machines != 3 {
		t.Errorf("expected stored network Machines count to persist at 3, got %d", world.Networks["203.0.113.1"].Machines)
	}
}

func TestNewMachine_DifferentPublicIpsAreIndependent(t *testing.T) {
	world := newTestWorld()

	a1 := world.NewMachine("203.0.113.1", "a1")
	b1 := world.NewMachine("198.51.100.1", "b1")
	a2 := world.NewMachine("203.0.113.1", "a2")

	if a1.LocalID != 1 || a2.LocalID != 2 {
		t.Errorf("expected network 203.0.113.1 local ids 1,2, got %d,%d", a1.LocalID, a2.LocalID)
	}
	if b1.LocalID != 1 {
		t.Errorf("expected independent network 198.51.100.1 to start its own local id count at 1, got %d", b1.LocalID)
	}
}
