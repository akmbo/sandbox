package scrobbles

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	assert := assert.New(t)

	t.Run("test valid log", func(t *testing.T) {
		path := "testdata/test_log.json"
		got, err := FromFile(path)

		assert.NoError(err)

		assert.Equal("test user", got.Username)

		for i := 0; i < 5000; i++ {
			num := strconv.Itoa(i + 1)
			want := Scrobble{
				Track:  "track " + num,
				Artist: "artist " + num,
				Album:  "album " + num,
				Date:   1649270931000 + uint(i),
			}
			assert.Equal(want, got.Scrobbles[i])
		}
	})

	t.Run("test invalid log", func(t *testing.T) {
		path := "testdata/invalid_log.txt"

		_, err := FromFile(path)
		assert.Error(err)
	})
}

func BenchmarkFromFile(b *testing.B) {
	path := "testdata/test_log.json"
	for i := 0; i < b.N; i++ {
		if _, err := FromFile(path); err != nil {
			b.Fatal("got error but wasn't expecting one")
		}
	}
}
