package registry

import "io"

// Interface ...
type Interface interface {
	PushImageByID(imageID string) (string, error)
	PushImage(reader io.Reader) (string, error)
	DownloadImage(ipfsHash string) (string, error)
	PullImage(ipfsHash string) (string, error)
}
