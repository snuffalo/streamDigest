package impl

import "github.com/snuffalo/streamDigest/models"

var m = make(map[uint64][]*models.Clip)

func GetDigestByStreamerId(id uint64) models.Digest  {
	var clips = m[id]
	var response = models.Digest{}
	for _, clip := range clips {
		response = append(response, clip.URL)
	}

	return response
}

func AddClipToDigestByStreamerId(c *models.Clip, id uint64) {
	m[id] = append(m[id], c)
}