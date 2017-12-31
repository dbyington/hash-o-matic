package memory

import (
    "time"
    "log"
    "errors"
)

type HashRecord struct {
    Hash string
    Ready bool
    Created time.Time
}

type MemoryDatabase struct {}

var Hashes []HashRecord

func (m MemoryDatabase) GetNextId() (Id int, err error) {
    Hashes = append(Hashes, HashRecord{})
    Id = len(Hashes)
    realId, _ := checkId(Id)
    Hashes[realId].Ready = false
    Hashes[realId].Created = time.Now()
    return Id, err
}

func checkId(Id int) (realId int, err error) {
    // the real Id is Id-1, since Id is the length of the array
    if len(Hashes) >= Id {
        realId = Id - 1
    } else {
        log.Printf("Id check failed: %d is invalid.\n", Id)
        return 0, errors.New("Bad Id")
    }
    return realId, nil
}

func (m MemoryDatabase) SaveHashWithId(hashString string, Id int) (realId int, err error) {
    realId, idErr := checkId(Id)
    if idErr != nil {
        return 0, idErr
    }
    Hashes[realId].Hash = hashString
    Hashes[realId].Ready = true
    return Id, err
}

func (m MemoryDatabase) GetHashById(Id int) (hashString string, err error) {
    realId, idErr := checkId(Id)
    if idErr != nil {
        return "", idErr
    }
    if Hashes[realId].Ready {
        hashString = Hashes[realId].Hash
        return hashString, err
    }
    return "", errors.New("hash string not ready")
}
