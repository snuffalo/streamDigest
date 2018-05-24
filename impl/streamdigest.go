package impl

import (
	"github.com/snuffalo/streamDigest/models"
	"database/sql"
	"fmt"
)


var m = make(map[uint64][]*models.Clip)

func GetDigestByStreamerId(id uint64) models.Digest  {
	var clips = m[id]
	var response = models.Digest{}
	for _, clip := range clips {
		response = append(response, clip)
	}

	return response
}

func AddClipToDigestByStreamerId(c *models.Clip, id uint64, db *sql.DB, log func(string, ...interface{})) bool {
	for _, clip := range m[id] {
		if IsClipEqual(c, clip) {
			return false
		}
	}
	m[id] = append(m[id], c)
	query := fmt.Sprintf("INSERT INTO clips (url) VALUES(\"%s\");", c.URL)
	_, err := db.Exec(query)
	if err == nil {
		return true
	} else {
		panic(err)
	}
}

func IsClipEqual(a *models.Clip, b *models.Clip) bool {
	return a.URL == b.URL
}