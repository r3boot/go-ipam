package storage

import (
	"bytes"
	"errors"
	"github.com/r3boot/go-ipam/models"
	"io"
	"os"
)

type NopCloser struct {
	io.Reader
}

func (NopCloser) Close() error { return nil }

func DumpAsset(fname string) (asset models.Asset, err error) {
	var (
		fs        os.FileInfo
		fd        *os.File
		num_bytes int
		raw_data  []byte
	)

	if fs, err = os.Stat(fname); err != nil {
		err = errors.New("DumpAsset: os.Stat() failed: " + err.Error())
		return
	}

	if fs.IsDir() {
		err = errors.New("DumpAsset: os.Stat(): is a directory")
		return
	}

	if fd, err = os.Open(fname); err != nil {
		err = errors.New("DumpAsset: os.Open() failed: " + err.Error())
		return
	}
	defer fd.Close()

	raw_data = make([]byte, fs.Size())

	if num_bytes, err = fd.Read(raw_data); err != nil {
		err = errors.New("DumpAsset: fd.Read() failed: " + err.Error())
		return
	}

	if num_bytes != int(fs.Size()) {
		err = errors.New("DumpAsset: fd.Read() failed: corrupted")
		return
	}

	asset = models.Asset(NopCloser{bytes.NewBuffer(raw_data)})

	return
}
