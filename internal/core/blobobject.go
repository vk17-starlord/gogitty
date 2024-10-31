package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
)

type Blob struct {
	Content     []byte
	ContentSize int64
	Header      string
	BlobPath    string
	Hash        string
	Buffer      bytes.Buffer
}

// Serialize the Blob object into a byte slice
func (b *Blob) Serialize(repo string) error {
	// Construct the header with content size and content
	contentAndHeader := fmt.Sprintf("blob %d\x00%s", b.ContentSize, string(b.Content))
	b.Header = contentAndHeader

	// Calculate the SHA-1 hash of the header
	sha := sha1.Sum([]byte(b.Header))
	b.Hash = fmt.Sprintf("%x", sha)

	// Set the path for storing the blob
	b.BlobPath = fmt.Sprintf("%s/objects/%v/%v", repo, b.Hash[0:2], b.Hash[2:])

	// Create a new zlib writer and compress the content only
	z := zlib.NewWriter(&b.Buffer)
	defer z.Close() // Ensure that the zlib writer is closed after use

	_, err := z.Write([]byte(contentAndHeader)) // Write the header to compress
	if err != nil {
		return fmt.Errorf("failed to write to zlib writer: %w", err)
	}

	if err := z.Close(); err != nil {
		return fmt.Errorf("failed to close zlib writer: %w", err)
	}

	return nil
}

// Deserialize byte data into a Blob object
func (b *Blob) Deserialize(data []byte) error {
	b.Content = data
	return nil
}

// Init initializes the Blob object by setting the content and content size.
func (b *Blob) Init(content []byte, size int64) {
	b.Content = content  // Assign content
	b.ContentSize = size // Assign content size
}
