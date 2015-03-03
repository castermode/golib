package cache

import (
	"fmt"
	"github.com/nicholaskh/assert"
	"math/rand"
	"testing"
)

type simpleStruct struct {
	int
	string
}

type complexStruct struct {
	int
	simpleStruct
}

var getTests = []struct {
	name       string
	keyToAdd   interface{}
	keyToGet   interface{}
	expectedOk bool
}{
	{"string_hit", "myKey", "myKey", true},
	{"string_miss", "myKey", "nonsense", false},
	{"simple_struct_hit", simpleStruct{1, "two"}, simpleStruct{1, "two"}, true},
	{"simeple_struct_miss", simpleStruct{1, "two"}, simpleStruct{0, "noway"}, false},
	{"complex_struct_hit", complexStruct{1, simpleStruct{2, "three"}},
		complexStruct{1, simpleStruct{2, "three"}}, true},
}

func TestGet(t *testing.T) {
	for _, tt := range getTests {
		lru := NewLruCache(0)
		lru.Set(tt.keyToAdd, 1234)
		val, ok := lru.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestDel(t *testing.T) {
	lru := NewLruCache(0)
	lru.Set("myKey", 1234)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestRemove failed.  Expected %d, got %v", 1234, val)
	}

	lru.Del("myKey")
	if _, ok := lru.Get("myKey"); ok {
		t.Fatal("TestRemove returned a removed entry")
	}
}

func TestInc(t *testing.T) {
	lru := NewLruCache(10)
	counter := lru.Inc("foo")
	assert.Equal(t, 1, counter)
	counter = lru.Inc("foo")
	assert.Equal(t, 2, counter)
	lru.Del("foo")
	counter = lru.Inc("foo")
	assert.Equal(t, 1, counter)
}

func BenchmarkCreateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("mc_stress:%d", rand.Int())
	}
}

func BenchmarkSet(b *testing.B) {
	b.ReportAllocs()
	lru := NewLruCache(0)
	var key string
	for i := 0; i < b.N; i++ {
		key = fmt.Sprintf("mc_stress:%d", rand.Int())
		lru.Set(key, 5)
	}
	b.SetBytes(int64(len(key)))
}

func BenchmarkGet(b *testing.B) {
	b.ReportAllocs()
	lru := NewLruCache(0)
	var key string
	for i := 0; i < b.N; i++ {
		key = fmt.Sprintf("mc_stress:%d", rand.Int())
		lru.Set(key, 5)
		lru.Get(key)
	}
	b.SetBytes(int64(len(key)))
}
