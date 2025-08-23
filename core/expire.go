package core

import (
	"log"
	"time"
)

func expireSample() float32 {
	//a sample size is 20
	var limit int = 20
	var expiredCnt int = 0

	for key, obj := range store {
		if obj.ExpiresAt != -1 {
			limit--
			if obj.ExpiresAt <= time.Now().UnixMilli() {
				delete(store, key)
				expiredCnt++
			}
		}
		if limit == 0 {
			break
		}
	}
	return float32(expiredCnt) / float32(limit)
}

func DeleteExpireKeys() {
	// Sampling approach: https://redis.io/commands/expire/
	//delte a sample and if the number of delete from the sample is more than 25% delete again from a new sample

	for {
		frac := expireSample()

		if frac < 0.25 {
			break
		}
	}
	log.Println("Deleted the expired Keys. total keys ", len(store))
}
