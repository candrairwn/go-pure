package utils

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// responseRecorder is a wrapper around [http.ResponseWriter] that records the status and bytes written during the response.
// It implements the [http.ResponseWriter] interface by embedding the original ResponseWriter.
type ResponseRecorder struct {
	http.ResponseWriter
	Status   int
	NumBytes int
}

// Header implements the [http.ResponseWriter] interface.
func (re *ResponseRecorder) Header() http.Header {
	return re.ResponseWriter.Header()
}

// Write implements the [http.ResponseWriter] interface.
func (re *ResponseRecorder) Write(b []byte) (int, error) {
	re.NumBytes += len(b)
	return re.ResponseWriter.Write(b)
}

// WriteHeader implements the [http.ResponseWriter] interface.
func (re *ResponseRecorder) WriteHeader(statusCode int) {
	re.Status = statusCode
	re.ResponseWriter.WriteHeader(statusCode)
}

// Hijack implements the [http.Hijacker] interface.
func (re *ResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := re.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter does not support the Hijacker interface")
	}
	return hijacker.Hijack()
}
