package proxy_reader

import (
	"os"

	"github.com/cheggaaa/pb/v3"
)

// ProxyReader holds a file and a progress bar. We use it to implement an
// io.Reader when uploading files to s3. This gives us the opportunity
// to write our own Read and ReadAt methods and update the progress bar
// throughout.
type ProxyReader struct {
	file        *os.File
	progressBar *pb.ProgressBar
}

// New creates a new ProxyReader assigning values to the private file and the *pb.ProgressBar fields.
func New(file *os.File, bar *pb.ProgressBar) *ProxyReader {
	return &ProxyReader{
		file:        file,
		progressBar: bar,
	}
}

// Read must be implemented for ProxyReader, however we just pass the slice
// of bytes argument right through to the file's Read method. We need no extra
// functionality here.
func (pr *ProxyReader) Read(p []byte) (int, error) {
	return pr.file.Read(p)
}

// ReadAt passes the slice of bytes right to file.ReadAt() -> kind of like a "super" call
// and then updates the ProxyReader's progress bar on each read
func (pr *ProxyReader) ReadAt(p []byte, off int64) (int, error) {
	// Call original ReadAt method with the slice of bytes
	n, err := pr.file.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	// No idea why the number of bytes read needs to be divided by 2, but read somewhere
	// that "maybe request is read once on sign and actually sends call ReadAt again"
	pr.progressBar.Add(n / 2)

	return n, err
}

// Seek must be implemented for ProxyReader, however we just pass offset and whence args
// right through to the file's Seek method. We need no extra functionality here.
func (pr *ProxyReader) Seek(offset int64, whence int) (int64, error) {
	return pr.file.Seek(offset, whence)
}
