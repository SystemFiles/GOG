package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"runtime"
	"strings"
)

func UnarchiveBinary(r io.Reader, binaryName string) (*bytes.Reader, error) {
	switch strings.ToLower(runtime.GOOS) {
	case "windows":
		return unzipBinary()
	case "linux", "darwin":
		b, err := untarBinary(r, binaryName)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(b), nil
	}
}

func unzipBinary(r io.Reader, binaryName string) ([]byte, error) {

}

func untarBinary(r io.Reader, binaryName string) ([]byte, error) {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		// if no gog binary is found return with error
		case err == io.EOF:
			return nil, errors.New("failed to locate required binary")
		case err != nil:
			return nil, err
		case header == nil:
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if header.Name == binaryName {
				var binary []byte
				if _, err := tr.Read(binary); err != nil {
					return nil, err
				}

				return binary, nil
			}
		}
	}
}