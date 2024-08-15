package dht

import (
	"testing"
)
func TestGeneral(t *testing.T) {
	//newChord
	chord, err := NewChord("root", 3099)
	if err != nil {
		t.Fatalf("could not create new chord: %v", err)
		return
	}
	t.Logf("Response: %v", chord)

	//joinChord
	j_chord, err := JoinChord("node2", "root", 3099)
	if err != nil {
		t.Fatalf("could not join chord: %v", err)
		return
	}
	t.Logf("Response: %v", j_chord)

	//set
	err = chord.Set("key1", "value1")
	if err != nil {
		t.Fatalf("could not call Set: %v", err)
		return
	}

	//get
	r, err := chord.Get("key1")
	if err != nil {
		t.Fatalf("could not call Get: %v", err)
		return
	}
	if r != "value1" {
		t.Fatalf("wrong value: received %s, expected value1", r)
		return
	}
	t.Logf("Response: %v", r)

	//delete
	err = chord.Delete("key1")
	if err != nil {
		t.Fatalf("could not call delete: %v", err)
		return
	}

	//getAllKeys
	err = chord.Set("key1", "value1")
	if err != nil {
		t.Fatalf("could not call Set: %v", err)
		return
	}
	err = chord.Set("key2", "value2")
	if err != nil {
		t.Fatalf("could not call Set: %v", err)
		return
	}
	keys, err := chord.GetAllKeys()
	if len(keys) != 2 {
		t.Fatalf("wrong value: received %d, expected 2", len(keys))
		return
	}
	t.Logf("Response: %v", keys)

	//IsFirst
	is_first, err := chord.IsFirst()
	if err != nil {
		t.Fatalf("could not call IsFirst: %v", err)
		return
	}
	if !is_first {
		t.Fatalf("not first as expected")
	}
}
