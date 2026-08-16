package tests

import (
	"testing"

	"hackergame/objects"
)

func newTestNetwork() objects.PublicNetwork {
	return objects.PublicNetwork{
		PublicIp:          "203.0.113.1",
		MachinesByLocalID: make(map[int]int),
		ForwardedPorts:    make(map[int]int),
		LocalIDCounter:    0,
	}
}

func TestForwardPort_Success(t *testing.T) {
	network := newTestNetwork()

	got, err := network.ForwardPort(8080, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != 1 {
		t.Errorf("expected returned localID %d, got %d", 1, got)
	}
	if network.ForwardedPorts[8080] != 1 {
		t.Errorf("expected port 8080 to be forwarded to localID 1, got %d", network.ForwardedPorts[8080])
	}
}

func TestForwardPort_AlreadyForwarded(t *testing.T) {
	network := newTestNetwork()

	if _, err := network.ForwardPort(8080, 1); err != nil {
		t.Fatalf("unexpected error on first forward: %v", err)
	}

	got, err := network.ForwardPort(8080, 2)
	if err == nil {
		t.Fatal("expected error when forwarding an already-forwarded port, got nil")
	}
	if got != 1 {
		t.Errorf("expected existing forwarded localID %d to be returned, got %d", 1, got)
	}
	if network.ForwardedPorts[8080] != 1 {
		t.Errorf("expected port mapping to remain unchanged, got %d", network.ForwardedPorts[8080])
	}
}

func TestForwardPort_DifferentPortsIndependent(t *testing.T) {
	network := newTestNetwork()

	if _, err := network.ForwardPort(80, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := network.ForwardPort(443, 2); err != nil {
		t.Fatalf("unexpected error forwarding a different port: %v", err)
	}

	if network.ForwardedPorts[80] != 1 || network.ForwardedPorts[443] != 2 {
		t.Errorf("expected both ports forwarded independently, got %v", network.ForwardedPorts)
	}
}

func TestAddMachine_Success(t *testing.T) {
	network := newTestNetwork()
	machine := &objects.Machine{MachineID: 42, LocalID: 1}

	got, err := network.AddMachine(1, machine)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != 1 {
		t.Errorf("expected returned localID %d, got %d", 1, got)
	}
	if id, ok := network.MachinesByLocalID[1]; !ok || id != 42 {
		t.Errorf("expected MachinesByLocalID to map 1 -> 42, got %v (ok=%v)", id, ok)
	}
}

func TestAddMachine_AlreadyExists(t *testing.T) {
	network := newTestNetwork()
	first := &objects.Machine{MachineID: 1, LocalID: 1}
	if _, err := network.AddMachine(1, first); err != nil {
		t.Fatalf("unexpected error adding first machine: %v", err)
	}

	second := &objects.Machine{MachineID: 2, LocalID: 1}
	_, err := network.AddMachine(1, second)
	if err == nil {
		t.Fatal("expected error when adding a machine with a duplicate localID, got nil")
	}
	if network.MachinesByLocalID[1] != 1 {
		t.Errorf("expected MachinesByLocalID entry to remain unchanged (id=1), got %d", network.MachinesByLocalID[1])
	}
}

func TestAddMachine_DifferentIDsIndependent(t *testing.T) {
	network := newTestNetwork()
	m1 := &objects.Machine{MachineID: 10, LocalID: 1}
	m2 := &objects.Machine{MachineID: 20, LocalID: 2}

	if _, err := network.AddMachine(1, m1); err != nil {
		t.Fatalf("unexpected error adding first machine: %v", err)
	}
	if _, err := network.AddMachine(2, m2); err != nil {
		t.Fatalf("unexpected error adding second machine: %v", err)
	}

	if network.MachinesByLocalID[1] != 10 || network.MachinesByLocalID[2] != 20 {
		t.Errorf("expected both machines registered independently, got %v", network.MachinesByLocalID)
	}
}
