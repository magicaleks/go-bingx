package common

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func DecodeGzip(data []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	decodedMsg, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decodedMsg), nil
}
