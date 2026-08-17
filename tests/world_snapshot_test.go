package tests

import (
	"os"
	"testing"

	"hackergame/objects"
	"hackergame/storage"
)

// buildFilesystemTree adds a small directory tree with file contents under
// root, so the sample save file has something to look at for a given machine.
func buildFilesystemTree(t *testing.T, root *objects.DataNode, hostname string) {
	t.Helper()

	if err := root.AddChild("etc", true); err != nil {
		t.Fatalf("unexpected error adding /etc: %v", err)
	}
	etc, _ := root.Child("etc")
	if err := etc.AddChild("hostname", false); err != nil {
		t.Fatalf("unexpected error adding /etc/hostname: %v", err)
	}
	hostnameFile, _ := etc.Child("hostname")
	if err := hostnameFile.SetData([]byte(hostname)); err != nil {
		t.Fatalf("unexpected error setting /etc/hostname data: %v", err)
	}

	if err := root.AddChild("home", true); err != nil {
		t.Fatalf("unexpected error adding /home: %v", err)
	}
	home, _ := root.Child("home")
	if err := home.AddChild("notes.txt", false); err != nil {
		t.Fatalf("unexpected error adding /home/notes.txt: %v", err)
	}
	notes, _ := home.Child("notes.txt")
	if err := notes.SetData([]byte("todo: don't get hacked")); err != nil {
		t.Fatalf("unexpected error setting /home/notes.txt data: %v", err)
	}
}

// TestSaveWorld_ProducesInspectableFile builds a world with multiple public
// networks, multiple machines per network, forwarded ports, and a filesystem
// per machine, then saves it to tests/testdata/sample_world.json so its
// content can be inspected by hand. The file is intentionally left on disk
// (not written to t.TempDir()) rather than cleaned up.
func TestSaveWorld_ProducesInspectableFile(t *testing.T) {
	world := &objects.World{
		Networks:    make(map[string]objects.PublicNetwork),
		Machines:    make(map[int]objects.Machine),
		Filesystems: make(map[int]*objects.Filesystem),
		Players:     make(map[string]objects.Player),
	}

	webServer := world.NewMachine("203.0.113.1", "web-server")
	dbServer := world.NewMachine("203.0.113.1", "db-server")
	laptop := world.NewMachine("198.51.100.1", "laptop")

	network1 := world.Networks["203.0.113.1"]
	if _, err := network1.ForwardPort(8080, webServer.LocalID); err != nil {
		t.Fatalf("unexpected error forwarding port 8080: %v", err)
	}
	world.Networks["203.0.113.1"] = network1

	for _, machine := range []objects.Machine{webServer, dbServer, laptop} {
		fs := objects.NewFilesystem(machine.MachineID)
		buildFilesystemTree(t, fs.Root, machine.MachineName)
		world.Filesystems[machine.MachineID] = fs
	}

	path := "testdata/sample_world.json"
	if err := storage.SaveWorld(path, world); err != nil {
		t.Fatalf("unexpected error saving world: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected save file to exist at %s, got error: %v", path, err)
	}
	if info.Size() == 0 {
		t.Errorf("expected save file to be non-empty, got 0 bytes")
	}

	t.Logf("saved world to %s (%d bytes)", path, info.Size())
}
