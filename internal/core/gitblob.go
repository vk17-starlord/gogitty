package core

import (
	"errors"
)

// GitBlob represents a git blob object
type GitBlob struct {
	BlobData []byte
}

// Init initializes the GitBlob object
func (b *GitBlob) Init() error {
	b.BlobData = []byte{}
	return nil
}

// Serialize serializes the GitBlob object into a byte slice (just returns BlobData here)
func (b *GitBlob) Serialize() ([]byte, error) {
	if len(b.BlobData) == 0 {
		return nil, errors.New("no data to serialize")
	}
	return b.BlobData, nil
}

// Deserialize deserializes the byte slice into the GitBlob object
func (b *GitBlob) Deserialize(data []byte) error {
	if len(data) == 0 {
		return errors.New("no data to deserialize")
	}
	b.BlobData = data
	return nil
}

// Format returns the object format, which is "blob" for GitBlob
func (b *GitBlob) Format() string {
	return "blob"
}
