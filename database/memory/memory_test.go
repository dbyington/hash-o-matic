package memory

import (
    "testing"
    //"time"
    //"os"
)


type Database interface {
    GetNextId() (Id int, err error)
    SaveHashWithId(hashString string, Id int) (realId int, err error)
    GetHashById(Id int) (hashString string, err error)
}

var DB = new(MemoryDatabase)

var dbi = Database(DB)

var hashString = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
var reversedHashString = "==Q7fg+A2us7dF0zJlkZPZbKv0ZCKVz8B5Q0bEqdFjv6jFIRYzBlPaRaTa/pzLVK+xyErAQDtwVdzlUg56BWhHEZ"

func TestGetNextId(t *testing.T) {
    Id, err := dbi.GetNextId()
    if err != nil || Id != 1 {
        t.Errorf("Error getting first Id; expected 1 got %d", Id)
    }
}

func TestSaveHashWithId(t *testing.T) {
    Id, err := dbi.GetNextId()
    if err != nil {
        t.Error("Failed to get next Id:", err)
    }

    Id, err = dbi.SaveHashWithId(hashString, Id)
    if err != nil {
        t.Error("Error saving hashString:", err)
    }
    Id, err = dbi.GetNextId()
    savedId, err := dbi.SaveHashWithId(reversedHashString, Id)
    if err != nil {
        t.Error("Error saving hashString:", err)
    }

    if err != nil {
        t.Error("Failed to get next Id:", err)
    }
    if savedId != Id {
        t.Errorf("Failed to save hashString; expected Id %d got %d\n", Id, savedId)
    }

    savedId, err = dbi.SaveHashWithId(hashString, 42)
    if err == nil {
        t.Error("Expected idErr: 'Bad Id', got:", err)
    }
}

func TestGetHashById(t *testing.T) {
    Id, err := dbi.GetNextId()
    returnId, err := dbi.SaveHashWithId(hashString, Id)
    if err != nil || returnId != Id{
        t.Error("Error saving hashString:", err, returnId)
    }
    Id, err = dbi.GetNextId()
    returnId, err = dbi.SaveHashWithId(reversedHashString, Id)

    testId := returnId - 1
    savedHash, err := dbi.GetHashById(testId)
    if err != nil {
        t.Error("Problem getting hash", err)
    }
    if savedHash != hashString {
        t.Errorf("Saved hash mismatch; expected %s got %s\n", hashString, savedHash)
    }
}