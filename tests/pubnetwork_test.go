package tests

import (
	"testing"

	"hackergame/objects"
)

func newTestNetwork() objects.PublicNetwork {
	return objects.PublicNetwork{
		Ip:             "203.0.113.1",
		LocalNet:       make(map[string]int),
		ForwardedPorts: make(map[int]string),
		Machines:       0,
	}
}

func TestForwardPort_Success(t *testing.T) {
	network := newTestNetwork()

	got, err := network.ForwardPort(8080, "192.168.1.1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "192.168.1.1" {
		t.Errorf("expected returned localIp %q, got %q", "192.168.1.1", got)
	}
	if network.ForwardedPorts[8080] != "192.168.1.1" {
		t.Errorf("expected port 8080 to be forwarded to 192.168.1.1, got %q", network.ForwardedPorts[8080])
	}
}

func TestForwardPort_AlreadyForwarded(t *testing.T) {
	network := newTestNetwork()

	if _, err := network.ForwardPort(8080, "192.168.1.1"); err != nil {
		t.Fatalf("unexpected error on first forward: %v", err)
	}

	got, err := network.ForwardPort(8080, "192.168.1.2")
	if err == nil {
		t.Fatal("expected error when forwarding an already-forwarded port, got nil")
	}
	if got != "192.168.1.1" {
		t.Errorf("expected existing forwarded ip %q to be returned, got %q", "192.168.1.1", got)
	}
	if network.ForwardedPorts[8080] != "192.168.1.1" {
		t.Errorf("expected port mapping to remain unchanged, got %q", network.ForwardedPorts[8080])
	}
}

func TestForwardPort_DifferentPortsIndependent(t *testing.T) {
	network := newTestNetwork()

	if _, err := network.ForwardPort(80, "192.168.1.1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := network.ForwardPort(443, "192.168.1.2"); err != nil {
		t.Fatalf("unexpected error forwarding a different port: %v", err)
	}

	if network.ForwardedPorts[80] != "192.168.1.1" || network.ForwardedPorts[443] != "192.168.1.2" {
		t.Errorf("expected both ports forwarded independently, got %v", network.ForwardedPorts)
	}
}

func TestAddMachine_Success(t *testing.T) {
	network := newTestNetwork()
	machine := &objects.Machine{MachineID: 42, LocalIp: "192.168.1.2"}

	got, err := network.AddMachine("192.168.1.2", machine)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "192.168.1.2" {
		t.Errorf("expected returned localIp %q, got %q", "192.168.1.2", got)
	}
	if id, ok := network.LocalNet["192.168.1.2"]; !ok || id != 42 {
		t.Errorf("expected LocalNet to map 192.168.1.2 -> 42, got %v (ok=%v)", id, ok)
	}
}

func TestAddMachine_AlreadyExists(t *testing.T) {
	network := newTestNetwork()
	first := &objects.Machine{MachineID: 1, LocalIp: "192.168.1.2"}
	if _, err := network.AddMachine("192.168.1.2", first); err != nil {
		t.Fatalf("unexpected error adding first machine: %v", err)
	}

	second := &objects.Machine{MachineID: 2, LocalIp: "192.168.1.2"}
	_, err := network.AddMachine("192.168.1.2", second)
	if err == nil {
		t.Fatal("expected error when adding a machine with a duplicate localIp, got nil")
	}
	if network.LocalNet["192.168.1.2"] != 1 {
		t.Errorf("expected LocalNet entry to remain unchanged (id=1), got %d", network.LocalNet["192.168.1.2"])
	}
}
