package util

import "testing"

func TestHashString(t *testing.T) {
    mockPassword := "angryMonkey"
    mockHash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

    actualHash := HashString(mockPassword)
    if (actualHash != mockHash) {
        t.Errorf("Failed to hash %s; expected %s, got %s\n", mockPassword, mockHash, actualHash)
    }
}