package ttlib

import (
	"encoding/binary"
	"fmt"
	"io"
)

// readBytes reads a uint16 number of bytes to read from the reader,
// then reads that number of bytes into a buffer to return.  If the
// number of bytes is greater than `limit`, an error is returned.
func readBytes(conn io.Reader, limit uint16) ([]byte, error) {
	var length uint16
	if err := binary.Read(conn, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	if length > limit {
		return nil, fmt.Errorf("Data too long: %v bytes.", length)
	}

	buf := make([]byte, length)

	if _, err := conn.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

// writeBytes is the counterpart to readBytes. It writes the data
// length to the writer, then writes the data.
func writeBytes(conn io.Writer, data []byte) error {
	err := binary.Write(conn, binary.LittleEndian, uint16(len(data)))
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
