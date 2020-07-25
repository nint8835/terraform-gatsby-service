package service

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
)

// CreateTarFile creates an in-memory TAR file for the purposes of creating a plain-text file within a Docker container.
func CreateTarFile(path string, contents string) (*tar.Reader, error) {
	tarBuffer := bytes.NewBuffer([]byte{})

	tarBufferWriter := bufio.NewWriter(tarBuffer)
	tarWriter := tar.NewWriter(tarBufferWriter)

	header := &tar.Header{
		Name: path,
		Mode: 0600,
		Size: int64(len(contents)),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return nil, fmt.Errorf("failed to write file header: %w", err)
	}

	if _, err := tarWriter.Write([]byte(contents)); err != nil {
		return nil, fmt.Errorf("failed to write contents to tar file: %w", err)
	}

	if err := tarWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}

	tarBufferReader := bufio.NewReader(tarBuffer)
	return tar.NewReader(tarBufferReader), nil
}
