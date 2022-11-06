package resources

import (
	"bytes"
	"compress/zlib"
	"io"
)

func CompressContent(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(b)
	if err != nil {
		return []byte{}, err
	}
	err = w.Close()
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func DecompressContent(b []byte) ([]byte, error) {
	buf := bytes.NewReader(b)
	r, err := zlib.NewReader(buf)
	if err != nil {
		return []byte{}, err
	}

	res, err := io.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}

	return res, nil
}
