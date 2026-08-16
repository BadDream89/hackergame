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

	machine := world.NewMachine("203.0.113.1")

	if machine.PublicIp != "203.0.113.1" {
		t.Errorf("expected PublicIp %q, got %q", "203.0.113.1", machine.PublicIp)
	}
	if _, exists := world.Networks["203.0.113.1"]; !exists {
		t.Fatal("expected a network to be created for a new publicIp")
	}
	if machine.LocalIp != "192.168.1.2" {
		t.Errorf("expected first machine's LocalIp %q, got %q", "192.168.1.2", machine.LocalIp)
	}
}

func TestNewMachine_ReturnsNonZeroMachineID(t *testing.T) {
	world := newTestWorld()

	machine := world.NewMachine("203.0.113.1")

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
		LocalNet:       make(map[string]int),
		ForwardedPorts: make(map[int]string),
		Machines:       5,
	}

	machine := world.NewMachine("203.0.113.1")

	if machine.PublicIp != "203.0.113.1" {
		t.Errorf("expected PublicIp %q, got %q", "203.0.113.1", machine.PublicIp)
	}
	if machine.LocalIp != "192.168.1.6" {
		t.Errorf("expected LocalIp derived from existing Machines count %q, got %q", "192.168.1.6", machine.LocalIp)
	}
}

func TestNewMachine_PersistsMachinesCountAcrossCalls(t *testing.T) {
	world := newTestWorld()

	first := world.NewMachine("203.0.113.1")
	second := world.NewMachine("203.0.113.1")
	third := world.NewMachine("203.0.113.1")

	if first.LocalIp != "192.168.1.2" {
		t.Errorf("expected first machine LocalIp %q, got %q", "192.168.1.2", first.LocalIp)
	}
	if second.LocalIp != "192.168.1.3" {
		t.Errorf("expected second machine LocalIp %q, got %q", "192.168.1.3", second.LocalIp)
	}
	if third.LocalIp != "192.168.1.4" {
		t.Errorf("expected third machine LocalIp %q, got %q", "192.168.1.4", third.LocalIp)
	}
	if world.Networks["203.0.113.1"].Machines != 4 {
		t.Errorf("expected stored network Machines count to persist at 4, got %d", world.Networks["203.0.113.1"].Machines)
	}
}
