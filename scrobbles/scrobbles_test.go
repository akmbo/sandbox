package scrobbles

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromJSON(t *testing.T) {
	assert := assert.New(t)

	t.Run("test valid log", func(t *testing.T) {
		file, err := os.Open("testdata/test_log.json")
		assert.NoError(err)
		defer file.Close()

		got, err := FromJSON(file)

		assert.NoError(err)

		assert.Equal("test user", got.Username)

		for i := 0; i < 5000; i++ {
			num := strconv.Itoa(i + 1)
			want := Scrobble{
				Track:  "track " + num,
				Artist: "artist " + num,
				Album:  "album " + num,
				Date:   time.Unix(0, msToNs(1649270931000+int64(i))),
			}
			assert.Equal(want, got.Scrobbles[i])
		}
	})

	t.Run("test invalid log", func(t *testing.T) {
		file, err := os.Open("testdata/invalid_log.txt")
		assert.NoError(err)
		defer file.Close()

		_, err = FromJSON(file)
		assert.Error(err)
	})
}

func BenchmarkFromJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file, err := os.Open("testdata/test_log.json")
		if err != nil {
			b.Fatalf("got error opening file")
		}
		defer file.Close()

		if _, err := FromJSON(file); err != nil {
			b.Fatal("got error but wasn't expecting one")
		}
	}
}
