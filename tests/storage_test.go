package tests

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"hackergame/objects"
	"hackergame/storage"
)

type samplePayload struct {
	Name  string
	Count int
}

func TestSaveLoad_GenericRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.json")
	want := samplePayload{Name: "test", Count: 42}

	if err := storage.Save(path, want); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	var got samplePayload
	if err := storage.Load(path, &got); err != nil {
		t.Fatalf("unexpected error loading: %v", err)
	}

	if got != want {
		t.Errorf("expected loaded payload %+v, got %+v", want, got)
	}
}

func TestLoad_MissingFileErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.json")

	var got samplePayload
	if err := storage.Load(path, &got); err == nil {
		t.Fatal("expected error loading a nonexistent file, got nil")
	}
}

func buildTestFilesystem(t *testing.T, machineID int) *objects.Filesystem {
	t.Helper()

	fs := objects.NewFilesystem(machineID)
	if err := fs.Root.AddChild("home", true); err != nil {
		t.Fatalf("unexpected error adding dir: %v", err)
	}
	home, _ := fs.Root.Child("home")
	if err := home.AddChild("notes.txt", false); err != nil {
		t.Fatalf("unexpected error adding file: %v", err)
	}
	notes, _ := home.Child("notes.txt")
	if err := notes.SetData([]byte("hello world")); err != nil {
		t.Fatalf("unexpected error setting data: %v", err)
	}

	return fs
}

func TestSaveLoadWorld_RoundTrip(t *testing.T) {
	world := newTestWorld()
	machine := world.NewMachine("203.0.113.1", "test-pc")
	world.Filesystems = map[int]*objects.Filesystem{
		machine.MachineID: buildTestFilesystem(t, machine.MachineID),
	}

	path := filepath.Join(t.TempDir(), "world.json")
	if err := storage.SaveWorld(path, world); err != nil {
		t.Fatalf("unexpected error saving world: %v", err)
	}

	loaded, err := storage.LoadWorld(path)
	if err != nil {
		t.Fatalf("unexpected error loading world: %v", err)
	}

	if !reflect.DeepEqual(world.Networks, loaded.Networks) {
		t.Errorf("expected Networks %+v, got %+v", world.Networks, loaded.Networks)
	}
	if !reflect.DeepEqual(world.Machines, loaded.Machines) {
		t.Errorf("expected Machines %+v, got %+v", world.Machines, loaded.Machines)
	}

	loadedFs, ok := loaded.Filesystems[machine.MachineID]
	if !ok {
		t.Fatalf("expected filesystem for machine %d to survive round trip", machine.MachineID)
	}
	if loadedFs.MachineID != machine.MachineID {
		t.Errorf("expected loaded filesystem MachineID %d, got %d", machine.MachineID, loadedFs.MachineID)
	}

	home, ok := loadedFs.Root.Child("home")
	if !ok {
		t.Fatal("expected loaded filesystem to contain /home")
	}
	if !home.IsDir() {
		t.Error("expected /home to remain a directory")
	}

	notes, ok := home.Child("notes.txt")
	if !ok {
		t.Fatal("expected loaded filesystem to contain /home/notes.txt")
	}
	if notes.IsDir() {
		t.Error("expected /home/notes.txt to remain a file")
	}
	if string(notes.Data()) != "hello world" {
		t.Errorf("expected file data %q, got %q", "hello world", notes.Data())
	}
}

func TestLoadWorld_FileDataIsNotBase64Encoded(t *testing.T) {
	world := newTestWorld()
	machine := world.NewMachine("203.0.113.1", "test-pc")

	fs := objects.NewFilesystem(machine.MachineID)
	if err := fs.Root.AddChild("secret.txt", false); err != nil {
		t.Fatalf("unexpected error adding file: %v", err)
	}
	secret, _ := fs.Root.Child("secret.txt")
	plaintext := "the password is hunter2"
	if err := secret.SetData([]byte(plaintext)); err != nil {
		t.Fatalf("unexpected error setting data: %v", err)
	}
	world.Filesystems[machine.MachineID] = fs

	path := filepath.Join(t.TempDir(), "world.json")
	if err := storage.SaveWorld(path, world); err != nil {
		t.Fatalf("unexpected error saving world: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error reading saved file: %v", err)
	}
	if !strings.Contains(string(raw), plaintext) {
		t.Errorf("expected saved file to contain plaintext %q, got:\n%s", plaintext, raw)
	}

	loaded, err := storage.LoadWorld(path)
	if err != nil {
		t.Fatalf("unexpected error loading world: %v", err)
	}

	loadedFs, ok := loaded.Filesystems[machine.MachineID]
	if !ok {
		t.Fatalf("expected filesystem for machine %d to survive round trip", machine.MachineID)
	}
	loadedSecret, ok := loadedFs.Root.Child("secret.txt")
	if !ok {
		t.Fatal("expected loaded filesystem to contain secret.txt")
	}
	if string(loadedSecret.Data()) != plaintext {
		t.Errorf("expected loaded file data %q, got %q", plaintext, loadedSecret.Data())
	}
}

func TestLoadWorld_MissingFileErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.json")

	if _, err := storage.LoadWorld(path); err == nil {
		t.Fatal("expected error loading a nonexistent world file, got nil")
	}
}
