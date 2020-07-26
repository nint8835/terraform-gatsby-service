package service

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"time"
)

// CreateTarFile creates an in-memory TAR file for the purposes of creating a plain-text file within a Docker container.
func CreateTarFile(path string, contents string) (io.Reader, error) {
	tarBuffer := new(bytes.Buffer)

	tarBufferWriter := bufio.NewWriter(tarBuffer)
	tarWriter := tar.NewWriter(tarBufferWriter)

	header := &tar.Header{
		Name:    path,
		Mode:    0777,
		Size:    int64(len(contents)),
		ModTime: time.Now(),
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

	return tarBuffer, nil
}
