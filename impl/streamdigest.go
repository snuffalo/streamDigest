package impl

import "github.com/snuffalo/streamDigest/models"

func GetDigestByStreamerId() *models.Digest  {
	return &models.Digest{Message:"hello world!"}
}