package main

import (
	"encoding/gob"
	"os"
)

func loadGob(name string, e interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	err = dec.Decode(e)
	if err != nil {
		return err
	}

	return nil
}

func saveGob(name string, e interface{}) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	err = enc.Encode(e)
	if err != nil {
		return err
	}

	return nil
}
