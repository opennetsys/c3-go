package stringutil

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"
)

// CompactJSON ...
func CompactJSON(src []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	re := regexp.MustCompile(`^.*?(\[|\{)`)
	s := re.ReplaceAllString(string(src), "$1")

	re = regexp.MustCompile(`(\]|\})[^]}].*?$`)
	s = re.ReplaceAllString(s, "$1")

	if err := json.Compact(b, []byte(s)); err != nil {
		log.Println("fail to compact", err)
		return nil, err
	}

	return b.Bytes(), nil
}
