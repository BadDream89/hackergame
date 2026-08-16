package tests

import (
	"testing"

	"hackergame/objects"
)

func TestNewNode_Success(t *testing.T) {
	root := objects.NewDataNode("/", true)

	err := root.NewNode("file.txt", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	child, ok := root.Child("file.txt")
	if !ok {
		t.Fatal("expected child node to be added to parent")
	}
	if child.Name() != "file.txt" || child.IsDir() != false {
		t.Errorf("expected child with name %q isDir=false, got name=%q isDir=%v", "file.txt", child.Name(), child.IsDir())
	}
}

func TestNewNode_LazilyInitializesChildrenMap(t *testing.T) {
	root := objects.NewDataNode("/", true)

	if err := root.NewNode("dir1", true); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if root.ChildCount() != 1 {
		t.Fatalf("expected 1 child, got %d", root.ChildCount())
	}
	if _, ok := root.Child("dir1"); !ok {
		t.Fatal("expected dir1 to be present in children")
	}
}

func TestNewNode_DuplicateName(t *testing.T) {
	root := objects.NewDataNode("/", true)

	if err := root.NewNode("file.txt", false); err != nil {
		t.Fatalf("unexpected error on first insert: %v", err)
	}

	err := root.NewNode("file.txt", true)
	if err == nil {
		t.Fatal("expected error when adding a duplicate name, got nil")
	}
	if root.ChildCount() != 1 {
		t.Errorf("expected children to remain unchanged, got %d entries", root.ChildCount())
	}
}

func TestNewNode_OnNonDirectory(t *testing.T) {
	file := objects.NewDataNode("file.txt", false)

	err := file.NewNode("child", false)
	if err == nil {
		t.Fatal("expected error when adding a node under a non-directory, got nil")
	}
}

func TestChangeFileData_Success(t *testing.T) {
	file := objects.NewDataNode("file.txt", false)
	data := []byte("hello world")

	err := file.ChangeFileData(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(file.Data()) != "hello world" {
		t.Errorf("expected file data %q, got %q", "hello world", file.Data())
	}
}

func TestChangeFileData_OnDirectory(t *testing.T) {
	dir := objects.NewDataNode("dir", true)

	err := dir.ChangeFileData([]byte("nope"))
	if err == nil {
		t.Fatal("expected error when changing data on a directory node, got nil")
	}
	if dir.Data() != nil {
		t.Errorf("expected directory data to remain nil, got %v", dir.Data())
	}
}
