package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Device struct {
	Name      string  `json:"name"`
	Sensor    string  `json:"sensor"`
	Threshold float64 `json:"threshold"`
}

type Room struct {
	Name    string   `json:"name"`
	Devices []Device `json:"devices"`
}

type Rooms struct {
	Rooms []Room `json:"rooms"`
}

func (r *Rooms) GetRoom(name string) *Room {
	for _, room := range r.Rooms {
		if room.Name == name {
			return &room
		}
	}
	return nil
}

func ParseRooms(path string) (*Rooms, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	// we initialize our Users array
	rooms := &Rooms{}
	err = json.Unmarshal(byteValue, rooms)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
