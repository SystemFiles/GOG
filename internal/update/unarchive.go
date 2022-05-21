package update

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
)

func UntarBinary(r io.Reader, binaryName string) (*tar.Reader, error) {
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
				return tr, nil
			}
		}
	}
}