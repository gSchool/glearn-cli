package proxy_reader

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/cheggaaa/pb/v3"
)

type ProxyReader struct {
	file        *os.File        // file
	progressBar *pb.ProgressBar // Every call to read can update progress bar
	totalBytes  int64           // Total number of bytes in file
	bytesRead   int64           // Total # of bytes transferred
}

// New creates a new ProxyReader assigning the file, the file's size, and the
// *pb.ProgressBar passed in to it
func New(file *os.File, bar *pb.ProgressBar) (*ProxyReader, error) {
	// Obtain FileInfo so we can look at length in bytes
	fileStats, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("Could not obtain file stats for %s", file.Name())
	}

	return &ProxyReader{
		file:        file,
		progressBar: bar,
		totalBytes:  fileStats.Size(),
	}, nil
}

// Read passes the slice of bytes to read right to file.Read() -> kind of like a "super" call
// with no additional functionality
func (pr *ProxyReader) Read(p []byte) (int, error) {
	return pr.file.Read(p)
}

func (pr *ProxyReader) ReadAt(p []byte, off int64) (int, error) {
	// Call original ReadAt method with slice the slice of bytes
	n, err := pr.file.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	// Update the proxy reader's bytes read with this chunk
	atomic.AddInt64(&pr.bytesRead, int64(n))

	// No idea why the number of bytes read needs to be divided by 2, but read somewhere
	// that "maybe request is read once on sign and actually sends call ReadAt again"
	pr.progressBar.Add(n / 2)

	return n, err
}

func (pr *ProxyReader) Seek(offset int64, whence int) (int64, error) {
	return pr.file.Seek(offset, whence)
}
