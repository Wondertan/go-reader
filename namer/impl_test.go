package namer

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

const testSize = 100

func TestReadNamer_Read(t *testing.T) {
	for s1 := 1; s1 <= testSize; s1++ {
		for s2 := 1; s2 <= testSize; s2++ {
			data := make([]byte, s1)
			_, _ = rand.Read(data)

			namer := NewReadNamer(data, t.Name())

			data2 := make([]byte, s2)
			n, err := namer.Read(data2)
			if err != nil && err != io.EOF {
				t.Fatal(err)
			}

			if len(data) > len(data2) {
				if !bytes.Equal(data[:n], data2) {
					t.Fail()
				}
			} else {
				if !bytes.Equal(data, data2[:len(data)]) {
					t.Fail()
				}
			}
		}
	}
}

func TestBatchReadNamer_Read(t *testing.T) {
	var res1 []byte
	var res2 []byte
	rs := make([]ReadNamer, testSize)
	batch := NewBatchReadNamer(rs)

	for s1 := 1; s1 <= testSize; s1++ {
		data := make([]byte, s1)
		_, _ = rand.Read(data)
		res1 = append(res1, data...)

		rs[s1-1] = NewReadNamer(data, t.Name())
	}

	for s1 := testSize; s1 >= 1; s1-- {
		data := make([]byte, s1)
		_, _ = batch.Read(data)
		res2 = append(res2, data...)
	}

	if !bytes.Equal(res1, res2) {
		t.Fail()
	}
}
