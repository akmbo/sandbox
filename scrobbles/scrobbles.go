package scrobbles

import (
	"encoding/json"
	"io"
	"os"
)

type Scrobble struct {
	Track  string `json:"track"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Date   uint   `json:"date"`
}

type ScrobbleLog struct {
	Username  string     `json:"username"`
	Scrobbles []Scrobble `json:"scrobbles"`
}

func FromFile(path string) (sLog ScrobbleLog, err error) {
	file, err := os.Open(path)
	if err != nil {
		return sLog, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		if err = decoder.Decode(&sLog); err == io.EOF {
			break
		} else if err != nil {
			return sLog, err
		}
	}

	return sLog, nil
}
