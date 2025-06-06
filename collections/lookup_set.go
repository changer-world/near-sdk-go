package collections

import (
	"crypto/sha256"
	"fmt"

	"github.com/vlmoon99/near-sdk-go/borsh"
	"github.com/vlmoon99/near-sdk-go/env"
)

// LookupSet is a non-iterable set implementation using NEAR contract storage.
type LookupSet[T any] struct {
	prefix []byte
}

// NewLookupSet creates a new LookupSet with the given storage prefix.
func NewLookupSet[T any](prefix []byte) *LookupSet[T] {
	return &LookupSet[T]{prefix: prefix}
}

// getKey creates a storage key from the prefix and the hashed value.
func (s *LookupSet[T]) getKey(element T) []byte {
	data, err := borsh.Serialize(element)
	if err != nil {
		panic(fmt.Sprintf("LookupSet.getKey: failed to serialize element: %v", err))
	}
	hash := sha256.Sum256(data)
	return append(s.prefix, hash[:]...)
}

// Insert adds an element to the set.
// Returns true if the element was not already present.
func (s *LookupSet[T]) Insert(element T) bool {
	key := s.getKey(element)
	exists, _ := env.StorageHasKey(key)
	_, err := env.StorageWrite(key, []byte{1}) // just mark the existence
	if err != nil {
		env.LogString("LookupSet.Insert: failed to write to storage: " + err.Error())
	}
	return !exists
}

// Contains checks whether the set contains the given element.
func (s *LookupSet[T]) Contains(element T) bool {
	key := s.getKey(element)
	exists, err := env.StorageHasKey(key)
	if err != nil {
		env.LogString("LookupSet.Contains: failed to check key: " + err.Error())
		return false
	}
	return exists
}

// Remove deletes the element from the set.
// Returns true if the element existed.
func (s *LookupSet[T]) Remove(element T) bool {
	key := s.getKey(element)
	removed, err := env.StorageRemove(key)
	if err != nil {
		env.LogString("LookupSet.Remove: failed to remove key: " + err.Error())
		return false
	}
	return removed
}
