package dataset

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

// Collection names, which are also the dataset folders and the bundle keys.
const (
	CollectionFish   = "fish"
	CollectionPlants = "plants"
)

// Collection describes one kind of record: where its schema lives and how its bundle is named.
type Collection struct {
	// Name is the dataset folder, the bundle key and the bundle file name (fish.json, plants.json).
	Name string
	// SchemaFile is relative to the parent of the dataset folder.
	SchemaFile string
}

// collections lists every kind of record the dataset holds.
var collections = map[string]Collection{
	CollectionFish:   {Name: CollectionFish, SchemaFile: filepath.Join("schema", "fish.schema.json")},
	CollectionPlants: {Name: CollectionPlants, SchemaFile: filepath.Join("schema", "plant.schema.json")},
}

// CollectionFor returns the collection named by the dataset folder, e.g. ../plants.
func CollectionFor(datasetDir string) (Collection, error) {
	name := filepath.Base(filepath.Clean(datasetDir))
	collection, found := collections[name]
	if !found {
		return Collection{}, fmt.Errorf("unknown dataset folder %q, expected one of: %s", name, collectionNames())
	}
	return collection, nil
}

// collectionNames lists the known folders for error messages.
func collectionNames() string {
	names := make([]string, 0, len(collections))
	for name := range collections {
		names = append(names, name)
	}
	slices.Sort(names)
	return strings.Join(names, ", ")
}
