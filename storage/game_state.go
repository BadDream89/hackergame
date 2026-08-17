package storage

import (
	"errors"
	"os"

	"hackergame/objects"
)

func SaveWorld(path string, world *objects.World) error {
	return save(path, world)
}

func LoadWorld(path string) (*objects.World, error) {
	world := &objects.World{
		Networks:    make(map[string]objects.PublicNetwork),
		Machines:    make(map[int]objects.Machine),
		Filesystems: make(map[int]*objects.Filesystem),
		Players:     make(map[string]objects.Player),
	}

	if err := load(path, world); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return world, nil
		}
		return nil, err
	}

	return world, nil
}
