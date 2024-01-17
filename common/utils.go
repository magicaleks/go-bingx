package common

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func DecodeGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decodedMsg, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return decodedMsg, nil
}
