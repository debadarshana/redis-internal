package core

import (
	"time"
)

var store map[string]*Obj

type Obj struct {
	Value     interface{}
	ExpiresAt int64 // absolute time when to expire in milliseconds
}

func init() {
	store = make(map[string]*Obj)
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
	store[k] = obj
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
