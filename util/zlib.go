package util

import (
	"bytes"
	"compress/zlib"
	"io"
)

// DoZlibCompress 进行zlib压缩
func DoZlibCompress(src string) string {
	if len(src) == 0 {
		return ""
	}
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_, err := w.Write([]byte(src))
	if err != nil {
		return ""
	}
	err = w.Close()
	if err != nil {
		return ""
	}
	return in.String()
}

// DoZlibUnCompress 进行zlib解压缩
func DoZlibUnCompress(compressSrc string) string {
	if len(compressSrc) == 0 {
		return ""
	}
	b := bytes.NewReader([]byte(compressSrc))
	var out bytes.Buffer
	read, _ := zlib.NewReader(b)
	_, err := io.Copy(&out, read)
	if err != nil {
		return ""
	}
	read.Close()
	return out.String()
}
