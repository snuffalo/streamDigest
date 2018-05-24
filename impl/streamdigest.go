package impl

import (
	"github.com/snuffalo/streamDigest/models"
	"database/sql"
	"fmt"
	"strconv"
	"github.com/go-redis/redis"
)

var primed = false

func setClipInCache(id uint64, url string, r *redis.Client) {
	_, err := r.LPush(string(id), url).Result()
	if (err != nil) {
		panic(err.Error())
	}
}

func getDigestInCache(id uint64, r *redis.Client) models.Digest {
	results, err := r.LRange(string(id), 0, -1).Result()
	if err != nil {
		panic(err.Error())
	}

	var digest = models.Digest{}
	for _, result := range results {
		digest = append(digest, &models.Clip{URL:result})
	}

	return digest
}

func GetDigestByStreamerId(id uint64, db *sql.DB, r *redis.Client) models.Digest  {
	if !primed {
		primeCache(db, r)
	}

	return getDigestInCache(id, r)
}

type AddClipToDigestByStreamerIdResult uint8
const (
	SUCCESS AddClipToDigestByStreamerIdResult = 0
	DUPLICATE_CLIP AddClipToDigestByStreamerIdResult = 1
	INSERT_ERROR AddClipToDigestByStreamerIdResult = 2
)

func AddClipToDigestByStreamerId(c *models.Clip, id uint64, db *sql.DB, r *redis.Client, log func(string, ...interface{})) AddClipToDigestByStreamerIdResult {
	if !primed {
		primeCache(db, r)
	}

	currentDigest := getDigestInCache(id, r)

	for _, clip := range currentDigest {
		if IsClipEqual(c, clip) {
			return DUPLICATE_CLIP
		}
	}
	setClipInCache(id, c.URL, r)
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

func primeCache(db *sql.DB, r *redis.Client) {
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
		//m[key] = append(m[key], &models.Clip{URL:string(values[1])})
		setClipInCache(key, string(values[1]), r)
	}
	if err = rows.Err(); err != nil {
		panic (err.Error())
	}
	primed = true
}