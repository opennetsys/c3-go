package stringutil

import (
	"bytes"
	"encoding/json"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// CompactJSON ...
func CompactJSON(src []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	re := regexp.MustCompile(`^.*?(\[|\{)`)
	s := re.ReplaceAllString(string(src), "$1")

	re = regexp.MustCompile(`(\]|\})[^]}].*?$`)
	s = re.ReplaceAllString(s, "$1")

	if s == "" {
		return []byte(`{}`), nil
	}

	if err := json.Compact(b, []byte(s)); err != nil {
		log.Printf("[util] fail to compact %s", err)
		return nil, err
	}

	return b.Bytes(), nil
}
