package database

import (
    "github.com/dbyington/hash-o-matic/database/memory"
)


type Database interface {
    GetNextId() (Id int, err error)
    SaveHashWithId(hashString string, Id int) (realId int, err error)
    GetHashById(Id int) (hashString string, err error)
}

var DB = new(memory.MemoryDatabase)

var dbInterface = Database(DB)
var GetNextId = dbInterface.GetNextId
var GetHashById = dbInterface.GetHashById
var SaveHashWithId = dbInterface.SaveHashWithId

/*
Methods needed:
GetNextId() (int, err)
GetHashById(Id) (hashString, err)
GetHashIdByHash(hashString) (hashString, err)
SaveHashWithId(hashString, Id) (err)
 */