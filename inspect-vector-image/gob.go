package main

import (
	"encoding/gob"
	"os"
)

func decodeGob(name string, e interface{}) error {
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
