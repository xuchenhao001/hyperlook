package cadvisor

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/google/cadvisor/cache/memory"
	"github.com/google/cadvisor/storage"
	"log"
)

var (
	storageDriver   = flag.String("storage_driver", "", fmt.Sprintf("Storage `driver` to use. Data is always cached shortly in memory, this controls where data is pushed besides the local cache. Empty means none. Options are: <empty>, %s", strings.Join(storage.ListDrivers(), ", ")))
	storageDuration = flag.Duration("storage_duration", 2*time.Minute, "How long to keep data stored (Default: 2min).")
)

// NewMemoryStorage creates a memory storage with an optional backend storage option.
func NewMemoryStorage() (*memory.InMemoryCache, error) {
	backendStorage, err := storage.New(*storageDriver)
	if err != nil {
		return nil, err
	}
	if *storageDriver != "" {
		log.Printf("Using backend storage type %q", *storageDriver)
	}
	log.Printf("Caching stats in memory for %v", *storageDuration)
	return memory.New(*storageDuration, backendStorage), nil
}