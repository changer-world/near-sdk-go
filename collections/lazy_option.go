package collections

import (
	"errors"

	"github.com/vlmoon99/near-sdk-go/borsh"
	"github.com/vlmoon99/near-sdk-go/env"
)

// LazyOption is a persistent optional value stored in contract storage.
type LazyOption[T any] struct {
	key []byte
}

// NewLazyOption creates a new LazyOption with a given key.
func NewLazyOption[T any](key []byte) *LazyOption[T] {
	return &LazyOption[T]{key: key}
}

// IsSome checks if the value exists in storage.
func (l *LazyOption[T]) IsSome() bool {
	exists, err := env.StorageHasKey(l.key)
	if err != nil {
		env.LogString("LazyOption.IsSome: error checking key existence: " + err.Error())
		return false
	}
	return exists
}

// IsNone returns true if the value does not exist.
func (l *LazyOption[T]) IsNone() bool {
	return !l.IsSome()
}

// Get retrieves the value from storage.
func (l *LazyOption[T]) Get() (*T, error) {
	exists, err := env.StorageHasKey(l.key)
	if err != nil {
		env.LogString("LazyOption.Get: error checking key existence: " + err.Error())
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	data, err := env.StorageRead(l.key)
	if err != nil {
		env.LogString("LazyOption.Get: error reading from storage: " + err.Error())
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var val T
	err = borsh.Deserialize(data, &val)
	if err != nil {
		env.LogString("LazyOption.Get: error deserializing data: " + err.Error())
		return nil, err
	}
	return &val, nil
}

// Set inserts or updates the value in storage.
func (l *LazyOption[T]) Set(value T) bool {
	existed, err := env.StorageHasKey(l.key)
	if err != nil {
		env.LogString("LazyOption.Set: error checking key existence: " + err.Error())
		existed = false
	}

	data, err := borsh.Serialize(value)
	if err != nil {
		env.LogString("LazyOption.Set: error serializing data: " + err.Error())
		return existed
	}

	_, err = env.StorageWrite(l.key, data)
	if err != nil {
		env.LogString("LazyOption.Set: error writing to storage: " + err.Error())
	}
	return existed
}

// Replace replaces the value and returns the old one.
func (l *LazyOption[T]) Replace(value T) (*T, error) {
	old, err := l.Get()
	if err != nil {
		env.LogString("LazyOption.Replace: error getting old value: " + err.Error())
		return nil, err
	}

	l.Set(value)
	return old, nil
}

// Remove deletes the value from storage.
func (l *LazyOption[T]) Remove() bool {
	existed, err := env.StorageHasKey(l.key)
	if err != nil {
		env.LogString("LazyOption.Remove: error checking key existence: " + err.Error())
		return false
	}
	if !existed {
		return false
	}

	removed, err := env.StorageRemove(l.key)
	if err != nil {
		env.LogString("LazyOption.Remove: error removing from storage: " + err.Error())
		return false
	}
	return removed
}

// Take removes the value and returns it.
func (l *LazyOption[T]) Take() (*T, error) {
	val, err := l.Get()
	if err != nil {
		env.LogString("LazyOption.Take: error getting value: " + err.Error())
		return nil, err
	}
	l.Remove()
	return val, nil
}

// MustGet returns the value or panics if it's missing or deserialization fails.
func (l *LazyOption[T]) MustGet() T {
	val, err := l.Get()
	if err != nil {
		panic(err)
	}
	if val == nil {
		panic(errors.New("LazyOption.MustGet: value not found"))
	}
	return *val
}
