package ddrv

import (
	"fmt"
	"io"
	"strings"
)

// do others that not defined in Driver interface
func mbody(reader io.Reader, filename string) (string, io.Reader) {
	boundary := "disgosucks"
	// Set the content type including the boundary
	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", boundary)

	CRLF := "\r\n"
	// fname := uuid.New().String()

	// Assemble all the parts of the multipart form-data
	// This includes the boundary, content disposition with the file name, content type,
	// a blank line to end headers, the actual content (reader), end of content,
	// and end of multipart data
	parts := []io.Reader{
		strings.NewReader("--" + boundary + CRLF),
		strings.NewReader(fmt.Sprintf(`Content-Disposition: form-data; name="file"; filename="%s"`, filename) + CRLF),
		strings.NewReader(fmt.Sprintf(`Content-Type: %s`, "application/octet-stream") + CRLF),
		strings.NewReader(CRLF),
		reader,
		strings.NewReader(CRLF),
		strings.NewReader("--" + boundary + "--" + CRLF),
	}

	// Return the content type and the combined reader of all parts
	return contentType, io.MultiReader(parts...)
}
