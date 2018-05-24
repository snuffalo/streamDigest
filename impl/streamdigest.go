package impl

import (
	"github.com/snuffalo/streamDigest/models"
	"database/sql"
	"fmt"
	"strconv"
)


var m = make(map[uint64][]*models.Clip)
var primed = false

func GetDigestByStreamerId(id uint64, db *sql.DB) models.Digest  {
	if !primed {
		primeCache(db)
	}

	var clips = m[id]
	var response = models.Digest{}
	for _, clip := range clips {
		response = append(response, clip)
	}

	return response
}

type AddClipToDigestByStreamerIdResult uint8
const (
	SUCCESS AddClipToDigestByStreamerIdResult = 0
	DUPLICATE_CLIP AddClipToDigestByStreamerIdResult = 1
	INSERT_ERROR AddClipToDigestByStreamerIdResult = 2
)

func AddClipToDigestByStreamerId(c *models.Clip, id uint64, db *sql.DB, log func(string, ...interface{})) AddClipToDigestByStreamerIdResult {
	if !primed {
		primeCache(db)
	}

	for _, clip := range m[id] {
		if IsClipEqual(c, clip) {
			return DUPLICATE_CLIP
		}
	}
	m[id] = append(m[id], c)
	query := fmt.Sprintf("INSERT INTO clips (streamerId, url) VALUES(%d,\"%s\");",id, c.URL)
	_, err := db.Exec(query)
	if err == nil {
		return SUCCESS
	} else {
		return INSERT_ERROR
	}
}

func IsClipEqual(a *models.Clip, b *models.Clip) bool {
	return a.URL == b.URL
}

func primeCache(db *sql.DB) {
	stmt, err := db.Prepare("SELECT streamerId, url from clips;")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		panic(err.Error())
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs :=  make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		key, err := strconv.ParseUint(string(values[0]), 10, len(values[0]))
		if err != nil {
			panic(err.Error())
		}
		m[key] = append(m[key], &models.Clip{URL:string(values[1])})
	}
	if err = rows.Err(); err != nil {
		panic (err.Error())
	}
	primed = true
}