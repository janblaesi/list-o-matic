// List-O-Matic Talking List Management System
// Copyright (C) 2021 Jan Blaesi
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
