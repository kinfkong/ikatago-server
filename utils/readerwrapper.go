package utils

import (
	"io"

	"log"

	"github.com/kinfkong/ikatago-server/errors"
)

type CloseHandler func(err error)
type IOReaderWrapper struct {
	reader         io.Reader
	pr             io.Reader
	pw             io.Writer
	clientClosed   bool
	onClientClosed CloseHandler
}

func NewIOReaderWrapper(reader io.Reader) *IOReaderWrapper {
	pr, pw := io.Pipe()

	result := IOReaderWrapper{
		reader:       reader,
		pr:           pr,
		pw:           pw,
		clientClosed: false,
	}
	// keep reading
	go func() {
		buf := make([]byte, 4096)
		var resultErr error = nil
		for {
			n, err := result.reader.Read(buf)
			if err != nil {
				resultErr = err
				break
			}
			// write to the pw
			written := 0
			for written < n {
				wn, err := pw.Write(buf[written:n])
				if err != nil {
					log.Printf("ERROR: FAILED TO WRITE: %v", err)
					resultErr = err
					break
				}
				written += wn
			}
			if written < n {
				log.Printf("ERROR: failed written enough bytes: %v, %v", written, n)
				resultErr = errors.CreateError(500, "invalid_write")
				break
			}
		}
		pr.Close()
		pw.Close()
		result.clientClosed = true
		if result.onClientClosed != nil {
			result.onClientClosed(resultErr)
		}
	}()
	return &result
}

func (reader *IOReaderWrapper) Read(p []byte) (n int, err error) {
	return reader.pr.Read(p)
}
