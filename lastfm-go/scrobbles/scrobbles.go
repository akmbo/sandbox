package scrobbles

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type Scrobble struct {
	Track  string    `json:"track"`
	Artist string    `json:"artist"`
	Album  string    `json:"album"`
	Date   time.Time `json:"date"`
}

func (s *Scrobble) GetTrackURL(username string) string {
	base := "https://www.last.fm/user/%s/library/music/%s/_/%s"
	link := fmt.Sprintf(base, username, s.Artist, s.Track)
	return strings.ReplaceAll(link, " ", "+")
}

func (s *Scrobble) GetAlbumURL(username string) string {
	base := "https://www.last.fm/user/%s/library/music/%s/%s"
	link := fmt.Sprintf(base, username, s.Artist, s.Album)
	return strings.ReplaceAll(link, " ", "+")
}

func (s *Scrobble) GetArtistURL(username string) string {
	base := "https://www.last.fm/user/%s/library/music/%s"
	link := fmt.Sprintf(base, username, s.Artist)
	return strings.ReplaceAll(link, " ", "+")
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
