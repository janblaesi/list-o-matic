package main

import (
	"bytes"
	"encoding/gob"
	"os"

	"github.com/google/uuid"
)

var lists map[uuid.UUID]TalkingList

func dumpListToFile() {
	var raw_bytes bytes.Buffer
	enc := gob.NewEncoder(&raw_bytes)

	if err := enc.Encode(lists); err != nil {
		println(err.Error())
		return
	}

	fh, err := os.Create("talking_lists")
	if err != nil {
		println(err.Error())
		return
	}
	defer fh.Close()

	fh.Write(raw_bytes.Bytes())
}

func readListFromFile() {
	var raw_bytes bytes.Buffer
	dec := gob.NewDecoder(&raw_bytes)

	fh, err := os.Open("talking_lists")
	if err != nil {
		println(err.Error())
		return
	}
	defer fh.Close()

	raw_bytes.ReadFrom(fh)

	if err := dec.Decode(&lists); err != nil {
		println(err.Error())
		return
	}
}
