package collections

import (
	"errors"
	"log"

	"github.com/vlmoon99/near-sdk-go/env"
)

// EncoderFunc defines a function to serialize a value of type T into bytes.
type EncoderFunc[T any] func(T) []byte

// DecoderFunc defines a function to deserialize bytes into a value of type T.
type DecoderFunc[T any] func([]byte) (T, error)

// LazyOption is a persistent optional value stored in contract storage.
type LazyOption[T any] struct {
	key    []byte
	encode EncoderFunc[T]
	decode DecoderFunc[T]
}

// New creates a new LazyOption with a given key and serialization functions.
func New[T any](key []byte, encode EncoderFunc[T], decode DecoderFunc[T]) *LazyOption[T] {
	return &LazyOption[T]{key: key, encode: encode, decode: decode}
}

// IsSome checks if the value exists in storage.
func (l *LazyOption[T]) IsSome() bool {
	exists, err := env.StorageHasKey(l.key)
	if err != nil {
		log.Println("IsSome: error checking key existence:", err)
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
		log.Println("Get: error checking key existence:", err)
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	data, err := env.StorageRead(l.key)
	if err != nil {
		log.Println("Get: error reading from storage:", err)
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	val, err := l.decode(data)
	if err != nil {
		log.Println("Get: error decoding data:", err)
		return nil, err
	}
	return &val, nil
}

// Set inserts or updates the value in storage.
func (l *LazyOption[T]) Set(value T) bool {
	existed, err := env.StorageHasKey(l.key)
	if err != nil {
		log.Println("Set: error checking key existence:", err)
		existed = false
	}

	data := l.encode(value)
	_, err = env.StorageWrite(l.key, data)
	if err != nil {
		log.Println("Set: error writing to storage:", err)
	}
	return existed
}

// Replace replaces the value and returns the old one.
func (l *LazyOption[T]) Replace(value T) (*T, error) {
	old, err := l.Get()
	if err != nil {
		log.Println("Replace: error getting old value:", err)
		return nil, err
	}

	l.Set(value)
	return old, nil
}

// Remove deletes the value from storage.
func (l *LazyOption[T]) Remove() bool {
	existed, err := env.StorageHasKey(l.key)
	if err != nil {
		log.Println("Remove: error checking key existence:", err)
		return false
	}
	if !existed {
		return false
	}

	removed, err := env.StorageRemove(l.key)
	if err != nil {
		log.Println("Remove: error removing from storage:", err)
		return false
	}
	return removed
}

// Take removes the value and returns it.
func (l *LazyOption[T]) Take() (*T, error) {
	val, err := l.Get()
	if err != nil {
		log.Println("Take: error getting value:", err)
		return nil, err
	}
	l.Remove()
	return val, nil
}

// MustGet returns the value or panics if it's missing or decoding fails.
func (l *LazyOption[T]) MustGet() T {
	val, err := l.Get()
	if err != nil {
		panic(err)
	}
	if val == nil {
		panic(errors.New("lazyoption: value not found"))
	}
	return *val
}
