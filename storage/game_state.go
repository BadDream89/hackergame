package storage

import "hackergame/objects"

func SaveWorld(path string, world *objects.World) error {
	return Save(path, world)
}

func LoadWorld(path string) (*objects.World, error) {
	world := &objects.World{
		Networks:    make(map[string]objects.PublicNetwork),
		Machines:    make(map[int]objects.Machine),
		Filesystems: make(map[int]*objects.Filesystem),
	}

	if err := Load(path, world); err != nil {
		return nil, err
	}

	return world, nil
}
