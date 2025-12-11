package helm

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	securejoin "github.com/cyphar/filepath-securejoin"
)

var errUnknownHeaderType = errors.New("unknown filesystem header")

// copyChart copies the chart archive to the destination directory
func copyChart(srcPath, dstPath string) error {
	src, err := os.Open(filepath.Clean(srcPath))
	if err != nil {
		return fmt.Errorf("cannot open chart in repository cache: %w", err)
	}

	defer func() {
		_ = src.Close()
	}()

	dst, err := os.Create(filepath.Clean(dstPath))
	if err != nil {
		return fmt.Errorf("cannot create chart in target path: %w", err)
	}

	defer func() {
		_ = dst.Close()
	}()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("copy source chart to destination: %w", err)
	}

	return nil
}

// extractChartTgz decompress the source gzip archive, then copy the files from the tar archive to the destination
func extractChartTgz(src, dst string) error {
	f, err := os.Open(filepath.Clean(src))
	if err != nil {
		return fmt.Errorf("cannot open chart in repository cache: %w", err)
	}

	err = extractTarGz(f, dst)
	if err != nil {
		return fmt.Errorf("extracting tgz: %w", err)
	}

	return nil
}

// extractTarGz extracts a gzipped tar archive to a directory.
func extractTarGz(r io.Reader, dst string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("unable to read gzip: %w", err)
	}

	defer func() {
		_ = gzr.Close()
	}()

	return extractTar(gzr, dst)
}

// extractTar extracts a tar archive to a directory.
func extractTar(r io.Reader, dst string) error {
	tarReader := tar.NewReader(r)

	for {
		header, hErr := tarReader.Next()
		if errors.Is(hErr, io.EOF) {
			break
		}

		if hErr != nil {
			return fmt.Errorf("unable to read tar content: %w", hErr)
		}

		// Strip the top-level directory from the path
		parts := strings.SplitN(header.Name, "/", 2)
		if len(parts) < 2 {
			continue
		}

		p, pErr := securejoin.SecureJoin(dst, parts[1])
		if pErr != nil {
			return fmt.Errorf("path contains invalid segments: %w", pErr)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(p, 0o750); err != nil {
				return fmt.Errorf("create directory: %w", err)
			}
		case tar.TypeReg:
			dir := filepath.Dir(p)
			if err := os.MkdirAll(dir, 0o750); err != nil {
				return fmt.Errorf("create parent folders for file: %w", err)
			}

			//nolint:gosec // G115 is a false alert, since it's all file modes it's safe
			outFile, err := os.OpenFile(filepath.Clean(p), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("read file content: %w", err)
			}

			_, err = io.CopyN(outFile, tarReader, header.Size)
			if err != nil {
				_ = outFile.Close()
				return fmt.Errorf("cannot copy file: %w", err)
			}

			_ = outFile.Close()
		case tar.TypeXGlobalHeader, tar.TypeXHeader:
			continue
		default:
			return fmt.Errorf("%w: %b in %s", errUnknownHeaderType, header.Typeflag, header.Name)
		}
	}

	return nil
}
