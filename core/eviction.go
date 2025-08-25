package core

import "log"

func evictFirst() {
	log.Printf("Eviction triggered: current store size = %d", len(store))
	for k := range store {
		log.Printf("Evicting key: %s", k)
		delete(store, k)
		return
	}
	log.Println("No keys to evict")
}

func Evict(evictionStrategy string) {
	log.Printf("Evict called with strategy: %s", evictionStrategy)
	switch evictionStrategy {
	case "simple-first":
		evictFirst()
	default:
		evictFirst()
	}
	log.Printf("After eviction: store size = %d", len(store))
}
