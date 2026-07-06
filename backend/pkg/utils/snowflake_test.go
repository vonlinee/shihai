package utils

import "testing"

func TestGenerateIDReturnsNonZeroID(t *testing.T) {
	id := GenerateID()

	if id == 0 {
		t.Fatal("GenerateID() = 0, want non-zero ID")
	}
}

func TestGenerateIDReturnsUniqueIDs(t *testing.T) {
	seen := make(map[uint64]struct{}, 1000)

	for i := 0; i < 1000; i++ {
		id := GenerateID()
		if _, ok := seen[id]; ok {
			t.Fatalf("GenerateID() returned duplicate ID %d", id)
		}
		seen[id] = struct{}{}
	}
}
