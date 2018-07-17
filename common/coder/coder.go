package coder

// AppendCode ...
func AppendCode(bytes []byte) []byte {
	return append([]byte{CurrentCode}, bytes...)
}

// StripCode ...
func StripCode(bytes []byte) (byte, []byte, error) {
	if bytes == nil || len(bytes) == 0 {
		return 0, nil, ErrNilBytes
	}

	return bytes[0], bytes[1:], nil
}

// ExtractCode ...
func ExtractCode(bytes []byte) (byte, error) {
	if bytes == nil || len(bytes) == 0 {
		return 0, ErrNilBytes
	}

	return bytes[0], nil
}
