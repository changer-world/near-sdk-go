package collections_test

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	systemMock := system.NewMockSystem()
	env.SetEnv(systemMock)
}

func newLazyOption(key []byte) *collections.LazyOption[string] {
	encode := func(s string) []byte {
		return []byte(s)
	}
	decode := func(data []byte) (string, error) {
		return string(data), nil
	}
	return collections.New[string](key, encode, decode)
}

func TestLazyOption_SetAndGet(t *testing.T) {
	opt := newLazyOption([]byte("my-key"))
	expected := "hello NEAR"

	opt.Set(expected)

	got, err := opt.Get()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("Expected value, got nil")
	}
	if *got != expected {
		t.Fatalf("Expected %q, got %q", expected, *got)
	}
}

func TestLazyOption_IsSome_IsNone(t *testing.T) {
	opt := newLazyOption([]byte("other-key"))

	if opt.IsSome() {
		t.Fatal("Expected IsSome() to be false for empty option")
	}
	if !opt.IsNone() {
		t.Fatal("Expected IsNone() to be true for empty option")
	}

	opt.Set("value")

	if !opt.IsSome() {
		t.Fatal("Expected IsSome() to be true after Set()")
	}
	if opt.IsNone() {
		t.Fatal("Expected IsNone() to be false after Set()")
	}
}

func TestLazyOption_Remove(t *testing.T) {
	opt := newLazyOption([]byte("remove-key"))
	opt.Set("to be removed")

	opt.Remove()

	if opt.IsSome() {
		t.Fatal("Expected IsSome() to be false after Remove()")
	}
	if !opt.IsNone() {
		t.Fatal("Expected IsNone() to be true after Remove()")
	}
}

func TestLazyOption_Take(t *testing.T) {
	opt := newLazyOption([]byte("take-key"))
	val := "temporary"
	opt.Set(val)

	result, err := opt.Take()
	if err != nil {
		t.Fatalf("Take() failed: %v", err)
	}
	if result == nil || *result != val {
		t.Fatalf("Expected %q, got %v", val, result)
	}

	if opt.IsSome() {
		t.Fatal("Expected IsSome() to be false after Take()")
	}
}

func TestLazyOption_MustGet(t *testing.T) {
	opt := newLazyOption([]byte("mustget-key"))
	opt.Set("important")

	result := opt.MustGet()
	if result != "important" {
		t.Fatalf("Expected 'important', got %q", result)
	}
}

func TestLazyOption_Replace(t *testing.T) {
	opt := newLazyOption([]byte("replace-key"))
	opt.Set("old")

	old, err := opt.Replace("new")
	if err != nil {
		t.Fatalf("Replace() failed: %v", err)
	}
	if old == nil || *old != "old" {
		t.Fatalf("Expected old value 'old', got %v", old)
	}

	current := opt.MustGet()
	if current != "new" {
		t.Fatalf("Expected new value 'new', got %q", current)
	}
}
