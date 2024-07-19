package scrobbles

import (
	"encoding/json"
	"io"
	"time"
)

type Scrobble struct {
	Track  string    `json:"track"`
	Artist string    `json:"artist"`
	Album  string    `json:"album"`
	Date   time.Time `json:"date"`
}

func (s *Scrobble) UnmarshalJSON(b []byte) error {
	type alias Scrobble
	aux := &struct {
		Date int64 `json:"date"`
		*alias
	}{
		alias: (*alias)(s),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	s.Date = time.Unix(0, msToNs(aux.Date))

	return nil
}

type ScrobbleLog struct {
	Username  string     `json:"username"`
	Scrobbles []Scrobble `json:"scrobbles"`
}

func msToNs(milliseconds int64) int64 {
	return milliseconds * 1000000
}

func FromJSON(reader io.Reader) (sLog ScrobbleLog, err error) {
	decoder := json.NewDecoder(reader)
	for {
		if err = decoder.Decode(&sLog); err == io.EOF {
			break
		} else if err != nil {
			return sLog, err
		}
	}

	return sLog, nil
}
