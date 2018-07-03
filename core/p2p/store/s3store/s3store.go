package s3store

import (
	"errors"

	s3ds "github.com/ipfs/go-ds-s3"
)

type Config struct {
	Domain    string
	AccessKey string
	SecretKey string
	Bucket    string
}

func New(cfg *Config) (*s3ds.S3Bucket, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	return s3ds.NewS3Datastore(&s3ds.Config{
		Domain:    cfg.Domain,
		AccessKey: cfg.AccessKey,
		SecretKey: cfg.SecretKey,
		Bucker:    cfg.Bucket,
	}), nil
}
