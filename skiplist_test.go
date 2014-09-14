package skiplist

import (
	"fmt"
	"testing"
)

type stringKey string

func (s stringKey) LessEq(other Key) bool {
	switch other := other.(type) {
	case stringKey:
		return s <= other
	}
	panic(fmt.Sprintf("invalid comparison of object %v with %v", s, other))
}

func (s stringKey) Eq(other Key) bool {
	return s == other
}

func (s stringKey) Less(other Key) bool {
	switch other := other.(type) {
	case stringKey:
		return s < other
	}
	panic(fmt.Sprintf("invalid comparison of object %v with %v", s, other))
}

func TestEmptySkipList(t *testing.T) {
	sl := NewSkipList(4, 0.5)

	actualLen := sl.Len()
	expectedLen := 0
	if actualLen != expectedLen {
		t.Errorf("actual %v != expected %v (invalid length)", actualLen, expectedLen)
	}

	actual, ok := sl.Find(stringKey("unexistedKey"))
	if actual != defaultValue {
		t.Errorf("actual %v != expected %v (invalid find)", actual, defaultValue)
	}
	if ok {
		t.Error("found element in empty list")
	}

	actual, ok = sl.Delete(stringKey("unexistedKey"))
	if actual != defaultValue {
		t.Errorf("actual %v != expected %v (invalid deleted element)", actual, defaultValue)
	}
	if ok {
		t.Error("invalid delete on empty list")
	}
}

func TestCRUDSkipList(t *testing.T) {
	var (
		sl      = NewSkipList(10, 0.5)
		value   = "value"
		entries = []struct {
			key   Key
			value string
		}{
			{stringKey("1"), "1"},
			{stringKey("2"), "1"},
			{stringKey("3"), "1"},
		}
	)

	for i, entry := range entries {
		sl.Add(entry.key, entry.value)
		actualLen := sl.Len()
		expectedLen := i + 1
		if actualLen != expectedLen {
			t.Fatalf("actual %d != expected %d (invalid length size after add an element)", actualLen, expectedLen)
		}
	}

	it := sl.Iterator()
	for i, entry := range entries {
		if it.Key() != entry.key {
			t.Fatalf("actual %q != expected %q (invalid key return by iterator)", it.Key(), entry.key)
		}
		if it.Value() != entry.value {
			t.Fatalf("actual %q != expected %q (invalid value return by iterator)", it.Value(), entry.value)
		}
		if i != len(entries)-1 && !it.Next() {
			t.Errorf("Iterator should not be to the end")
		}
	}
	if it.Next() {
		t.Errorf("Iterator should be to the end")
	}

	for _, entry := range entries {
		found, ok := sl.Find(entry.key)
		if !ok {
			t.Fatalf("element not added correctly")
		}
		if found != entry.value {
			t.Fatalf("actual %q != expected %q (invalid find)", found, entry.value)
		}
	}

	deleted, ok := sl.Delete(stringKey("unexistedKey"))
	if ok {
		t.Fatal("Should not deleted unexisted key")
	}
	if deleted != defaultValue {
		t.Fatal("Should return default value when the key does not exist")
	}

	for i, entry := range entries {
		deleted, ok := sl.Delete(entry.key)
		if !ok {
			t.Fatalf("element not deleted correctly")
		}
		if deleted != entry.value {
			t.Fatalf("the incorrect deleted element: actual %v != expected %v deleting %v", deleted, entry.value, entry.key)
		}
		actualLen := sl.Len()
		expectedLen := len(entries) - i - 1
		if actualLen != expectedLen {
			t.Fatalf("actual %d != expected %d (invalid length size after add an element)", actualLen, expectedLen)
		}
	}

	it = sl.Iterator()
	if it.Next() {
		t.Fatal("Iterator should fail (after deletion should not have any elements left)")
	}

	for _, entry := range entries {
		found, ok := sl.Find(entry.key)
		if found != defaultValue {
			t.Fatalf("actual %q != expected %q (invalid deleted element)", found, value)
		}
		if ok {
			t.Fatal("invalid delete on empty list")
		}
	}
}
