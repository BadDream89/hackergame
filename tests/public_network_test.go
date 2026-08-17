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

func TestUnforwardPort_Success(t *testing.T) {
	network := newTestNetwork()
	if _, err := network.ForwardPort(8080, 1); err != nil {
		t.Fatalf("unexpected error forwarding port: %v", err)
	}

	if err := network.UnforwardPort(8080); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, exists := network.ForwardedPorts[8080]; exists {
		t.Errorf("expected port 8080 to be removed from ForwardedPorts, got %v", network.ForwardedPorts)
	}
}

func TestUnforwardPort_NotForwardedErrors(t *testing.T) {
	network := newTestNetwork()

	err := network.UnforwardPort(8080)
	if err == nil {
		t.Fatal("expected error when unforwarding a port that was never forwarded, got nil")
	}
}

func TestUnforwardPort_DoesNotAffectOtherPorts(t *testing.T) {
	network := newTestNetwork()
	if _, err := network.ForwardPort(80, 1); err != nil {
		t.Fatalf("unexpected error forwarding port 80: %v", err)
	}
	if _, err := network.ForwardPort(443, 2); err != nil {
		t.Fatalf("unexpected error forwarding port 443: %v", err)
	}

	if err := network.UnforwardPort(80); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, exists := network.ForwardedPorts[80]; exists {
		t.Error("expected port 80 to be removed")
	}
	if network.ForwardedPorts[443] != 2 {
		t.Errorf("expected port 443 to remain forwarded to localID 2, got %d", network.ForwardedPorts[443])
	}
}
