package mocks

import "fmt"

// Test struct with BinaryMarshaler support
type BinaryVal struct {
	Data string
}

func (b BinaryVal) MarshalBinary() ([]byte, error) {
	return []byte("bin:" + b.Data), nil
}

type BadBinary struct{}

func (b BadBinary) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("custom binary marshal error")
}
