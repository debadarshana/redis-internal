package core

import (
	"log"
	"time"
)

var store map[string]*Obj
var storeConfig *StoreConfig

type Obj struct {
	Value     interface{}
	ExpiresAt int64 // absolute time when to expire in milliseconds
}

type StoreConfig struct {
	KeysLimit        int
	EvictionStrategy string
}

func init() {
	store = make(map[string]*Obj)
}

// InitStore initializes the store with configuration
func InitStore(config StoreConfig) {
	storeConfig = &config
}

func NewObj(value interface{}, durationMs int64) *Obj {
	expiresAt := int64(-1)
	if durationMs > 0 {
		expiresAt = time.Now().UnixMilli() + durationMs
	}
	return &Obj{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

func Put(k string, obj *Obj) {
	log.Printf("Put called for key '%s', current store size: %d, limit: %d", k, len(store),
		func() int {
			if storeConfig != nil {
				return storeConfig.KeysLimit
			} else {
				return -1
			}
		}())

	// Check if we need to evict before adding new key
	if storeConfig != nil && len(store) >= storeConfig.KeysLimit {
		// Only evict if the key doesn't already exist (we're adding a new key)
		if _, exists := store[k]; !exists {
			log.Printf("Triggering eviction before adding new key '%s'", k)
			Evict(storeConfig.EvictionStrategy)
		} else {
			log.Printf("Key '%s' already exists, updating without eviction", k)
		}
	}

	store[k] = obj
	log.Printf("Key '%s' stored, new store size: %d", k, len(store))
}

func Get(k string) *Obj {
	v := store[k]
	if v != nil {
		if v.ExpiresAt != -1 && time.Now().UnixMilli() >= v.ExpiresAt {
			// Key has expired, delete it
			delete(store, k)
			return nil
		}
		return v
	}
	return nil
}
func Del(k string) bool {
	if _, ok := store[k]; ok {
		delete(store, k)
		return true
	}
	return false
}
