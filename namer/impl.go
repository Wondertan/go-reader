package namer

import (
	"bytes"
	"io"
)

// Basic ReadNamer implementation
type readNamer struct {
	data []byte
	name string
}

// Basic ReadNamer constructor
func NewReadNamer(data []byte, name string) ReadNamer {
	return &readNamer{data: data, name: name}
}

// Read fills given byte slice with underlying basic ReadNamer buffer,
// returns EOF on bytes end
func (r *readNamer) Read(b []byte) (int, error) {
	if len(r.data) > len(b) {
		return copy(b, r.data[:len(b)]), nil
	} else {
		return copy(b[:len(r.data)], r.data), io.EOF
	}
}

// Name returns the name tighten to basic ReadNamer
func (r *readNamer) Name() string {
	return r.name
}

// ReadNamer implementation which allows to read from many ReadNamers
type batchReadNamer struct {
	readers []ReadNamer
	index   uint64

	buffer *bytes.Buffer
}

// BatchReadNamer constructor
func NewBatchReadNamer(readers []ReadNamer) ReadNamer {
	return &batchReadNamer{readers: readers, index: 0, buffer: new(bytes.Buffer)}
}

// Read fill given byte slice from underlying ReadNamers
// Read also remembers what was read
func (r *batchReadNamer) Read(b []byte) (int, error) {
	var err error

	defer func() {
		e := recover()
		if e == nil {
			return
		}

		if panicErr, ok := e.(error); ok {
			err = panicErr
			return
		}
	}()

	for {
		if r.buffer.Len() >= len(b) || int(r.index) == len(r.readers) {
			break
		}

		_, err = r.buffer.ReadFrom(r.readers[r.index])
		if err != nil {
			return 0, err
		}

		r.index++
	}

	n, err := r.buffer.Read(b)
	if err != nil && err != io.EOF {
		return n, err
	}

	if err == nil &&
		int(r.index) == len(r.readers) &&
		r.buffer.Len() == 0 {
		err = io.EOF
	}

	return n, err
}

// Name returns the name of the underlying ReadNamer which is currently in use
func (r *batchReadNamer) Name() string {
	return r.readers[r.index].Name()
}
