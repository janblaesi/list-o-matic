package main

import (
	"github.com/google/uuid"
)

// The representation of all available talking lists in RAM.
// Can be saved and/or retrieved from disk.
var lists map[uuid.UUID]TalkingList

// Initialize this pseudo database and try to read entries from file
func setupDatabase() error {
	lists = make(map[uuid.UUID]TalkingList)
	return readListFromFile()
}

// Dump the current database from RAM to disk
func dumpListToFile() error {
	return dumpJsonToFile(&lists, cfg.Database.TalkingListsPath)
}

// Read the current list from disk to RAM
func readListFromFile() error {
	return parseJsonFromFile(&lists, cfg.Database.TalkingListsPath)
}
