package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddAndGet(t *testing.T) {
	cache := NewCache(5 * time.Second)

	// Test adding and retrieving an item
	testKey := "test-key"
	testVal := []byte("test-value")

	cache.Add(testKey, testVal)

	val, found := cache.Get(testKey)
	if !found {
		t.Errorf("Expected to find key %s in cache, but it wasn't found", testKey)
	}

	if string(val) != string(testVal) {
		t.Errorf("Expected value %s, got %s", string(testVal), string(val))
	}

	// Test that a non-existent key returns not found
	_, found = cache.Get("non-existent-key")
	if found {
		t.Error("Expected non-existent key to return not found")
	}
}

func TestCacheExpiration(t *testing.T) {
	interval := 100 * time.Millisecond
	cache := NewCache(interval)

	testKey := "expiring-key"
	testVal := []byte("expiring-value")

	cache.Add(testKey, testVal)

	_, found := cache.Get(testKey)
	if !found {
		t.Error("Expected to find item immediately after adding")
	}

	time.Sleep(interval * 2)

	_, found = cache.Get(testKey)
	if found {
		t.Error("Expected item to be removed after expiration interval")
	}
}
