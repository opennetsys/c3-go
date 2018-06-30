package hexutil

import "strconv"

// DecodeUint64 decodes a hex string into a uint64
func DecodeUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(raw, 16, 64)
	if err != nil {
		return err
	}

	return num, err
}

// EncodeUint64 encodes i as a hex string
func EncodeUint64(i uint64) string {
	str := make([]byte, 2, 10)

	return string(strconv.AppendUint(str, i, 16))
}
