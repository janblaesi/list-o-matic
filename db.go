package main

import (
	"bytes"
	"encoding/gob"
	"os"

	"github.com/google/uuid"
)

var lists map[uuid.UUID]TalkingList

func dumpListToFile() {
	var rawBytes bytes.Buffer
	enc := gob.NewEncoder(&rawBytes)

	if err := enc.Encode(lists); err != nil {
		println("Failed to encode database for dumping: ", err.Error())
		return
	}

	fh, err := os.Create("talking_lists")
	if err != nil {
		println("Failed to open database dump file: ", err.Error())
		return
	}
	defer fh.Close()

	_, err = fh.Write(rawBytes.Bytes())
	if err != nil {
		println("Failed to dump current database to file: ", err.Error())
		return
	}
}

func readListFromFile() {
	var rawBytes bytes.Buffer
	dec := gob.NewDecoder(&rawBytes)

	fh, err := os.Open("talking_lists")
	if err != nil {
		println("Failed to open database dump file: ", err.Error())
		return
	}
	defer fh.Close()

	_, err = rawBytes.ReadFrom(fh)
	if err != nil {
		println("Failed to read from database dump file: ", err.Error())
		return
	}

	if err := dec.Decode(&lists); err != nil {
		println("Failed to decode database: ", err.Error())
		return
	}
}
