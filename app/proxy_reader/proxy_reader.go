package proxy_reader

import (
	"os"
)

type ProxyReader struct {
	file        *os.File        // file
	progressBar *pb.ProgressBar // Every call to read can update progress bar
}

// New creates a new ProxyReader assigning the file, the file's size, and the
// *pb.ProgressBar passed in to it
func New(file *os.File, bar *pb.ProgressBar) *ProxyReader {
	return &ProxyReader{
		file:        file,
		progressBar: bar,
	}
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

	// No idea why the number of bytes read needs to be divided by 2, but read somewhere
	// that "maybe request is read once on sign and actually sends call ReadAt again"
	pr.progressBar.Add(n / 2)

	return n, err
}

func (pr *ProxyReader) Seek(offset int64, whence int) (int64, error) {
	return pr.file.Seek(offset, whence)
}
