package objects

import "errors"

type Filesystem struct {
	MachineID int
	Root      *DataNode
}

type DataNode struct {
	name     string
	isDir    bool
	children map[string]*DataNode
	data     []byte
}

func NewDataNode(name string, isDir bool) *DataNode {
	return &DataNode{
		name:  name,
		isDir: isDir,
	}
}

func ensureChildAddable(node *DataNode, name string) error {
	if !node.isDir {
		return errors.New("cannot add node to a non-directory node")
	}

	if node.children == nil {
		node.children = make(map[string]*DataNode)
	}

	if _, exists := node.children[name]; exists {
		return errors.New("node with this name already exists")
	}

	return nil
}

func (node *DataNode) AddChild(name string, isDir bool) error {
	if err := ensureChildAddable(node, name); err != nil {
		return err
	}

	node.children[name] = &DataNode{
		name:     name,
		isDir:    isDir,
		children: nil,
		data:     nil,
	}

	return nil
}

func (node *DataNode) SetData(data []byte) error {
	if node.isDir {
		return errors.New("cannot add data to the directory node")
	}

	node.data = data

	return nil
}

func (node *DataNode) Name() string {
	return node.name
}

func (node *DataNode) IsDir() bool {
	return node.isDir
}

func (node *DataNode) Data() []byte {
	return node.data
}

func (node *DataNode) Child(name string) (*DataNode, bool) {
	child, ok := node.children[name]
	return child, ok
}

func (node *DataNode) ChildCount() int {
	return len(node.children)
}
