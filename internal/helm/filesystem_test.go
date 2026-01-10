package helm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyChart(t *testing.T) {
	tests := []struct {
		setup   func(t *testing.T) (srcPath, dstPath string)
		name    string
		errMsg  string
		wantErr bool
	}{
		{
			name: "positive case - successful copy",
			setup: func(t *testing.T) (string, string) {
				t.Helper()

				// Create a temporary source file with content
				tmpDir := t.TempDir()
				srcPath := filepath.Join(tmpDir, "source-chart.tgz")
				dstPath := filepath.Join(tmpDir, "dest-chart.tgz")

				content := []byte("fake chart content for testing")
				err := os.WriteFile(srcPath, content, 0o644)
				require.NoError(t, err)

				return srcPath, dstPath
			},
			wantErr: false,
		},
		{
			name: "negative case - source file does not exist",
			setup: func(t *testing.T) (string, string) {
				t.Helper()

				tmpDir := t.TempDir()
				srcPath := filepath.Join(tmpDir, "nonexistent.tgz")
				dstPath := filepath.Join(tmpDir, "dest.tgz")

				return srcPath, dstPath
			},
			wantErr: true,
			errMsg:  "cannot open chart in repository cache",
		},
		{
			name: "negative case - destination directory does not exist",
			setup: func(t *testing.T) (string, string) {
				t.Helper()

				tmpDir := t.TempDir()
				srcPath := filepath.Join(tmpDir, "source.tgz")
				dstPath := filepath.Join(tmpDir, "nonexistent", "dest.tgz")

				content := []byte("fake chart content")
				err := os.WriteFile(srcPath, content, 0o644)
				require.NoError(t, err)

				return srcPath, dstPath
			},
			wantErr: true,
			errMsg:  "cannot create chart in target path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcPath, dstPath := tt.setup(t)

			err := copyChart(srcPath, dstPath)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)

				return
			}

			require.NoError(t, err)

			// Verify destination file exists
			_, err = os.Stat(dstPath)
			require.NoError(t, err, "destination file should exist")

			// Verify contents match
			srcContent, err := os.ReadFile(srcPath)
			require.NoError(t, err)

			dstContent, err := os.ReadFile(dstPath)
			require.NoError(t, err)

			require.Equal(t, srcContent, dstContent, "source and destination content should match")
		})
	}
}

// createTestTarGz creates a test tar.gz archive with the specified files.
// Each file in the files map has a path (key) and content (value).
// The chartName is used as the top-level directory name (simulating Helm chart structure).
func createTestTarGz(t *testing.T, chartName string, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer

	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	for path, content := range files {
		// Add chartName as prefix (Helm charts have a top-level directory)
		fullPath := chartName + "/" + path

		header := &tar.Header{
			Name: fullPath,
			Mode: 0o644,
			Size: int64(len(content)),
		}

		err := tw.WriteHeader(header)
		require.NoError(t, err)

		_, err = tw.Write([]byte(content))
		require.NoError(t, err)
	}

	err := tw.Close()
	require.NoError(t, err)

	err = gzw.Close()
	require.NoError(t, err)

	return buf.Bytes()
}

func TestExtractChartTgz(t *testing.T) {
	tests := []struct {
		setup   func(t *testing.T) (srcPath, dstPath string)
		verify  func(t *testing.T, dstPath string)
		name    string
		errMsg  string
		wantErr bool
	}{
		{
			name: "positive case - successful extraction",
			setup: func(t *testing.T) (string, string) {
				t.Helper()

				tmpDir := t.TempDir()
				srcPath := filepath.Join(tmpDir, "chart.tgz")
				dstPath := filepath.Join(tmpDir, "extracted")

				// Create test tar.gz with sample files
				files := map[string]string{
					"Chart.yaml":            "name: test-chart\nversion: 1.0.0",
					"values.yaml":           "replicas: 3",
					"templates/deploy.yaml": "apiVersion: v1\nkind: Deployment",
				}
				tgzContent := createTestTarGz(t, "test-chart", files)

				err := os.WriteFile(srcPath, tgzContent, 0o644)
				require.NoError(t, err)

				return srcPath, dstPath
			},
			wantErr: false,
			verify: func(t *testing.T, dstPath string) {
				t.Helper()

				// Verify extracted files exist with correct content
				chartYaml, err := os.ReadFile(filepath.Join(dstPath, "Chart.yaml"))
				require.NoError(t, err)
				require.Contains(t, string(chartYaml), "name: test-chart")

				valuesYaml, err := os.ReadFile(filepath.Join(dstPath, "values.yaml"))
				require.NoError(t, err)
				require.Contains(t, string(valuesYaml), "replicas: 3")

				deployYaml, err := os.ReadFile(filepath.Join(dstPath, "templates", "deploy.yaml"))
				require.NoError(t, err)
				require.Contains(t, string(deployYaml), "kind: Deployment")
			},
		},
		{
			name: "negative case - source file does not exist",
			setup: func(t *testing.T) (string, string) {
				t.Helper()

				tmpDir := t.TempDir()
				srcPath := filepath.Join(tmpDir, "nonexistent.tgz")
				dstPath := filepath.Join(tmpDir, "extracted")

				return srcPath, dstPath
			},
			wantErr: true,
			errMsg:  "cannot open chart in repository cache",
		},
		{
			name: "negative case - invalid gzip data",
			setup: func(t *testing.T) (string, string) {
				t.Helper()

				tmpDir := t.TempDir()
				srcPath := filepath.Join(tmpDir, "invalid.tgz")
				dstPath := filepath.Join(tmpDir, "extracted")

				// Write invalid gzip data
				err := os.WriteFile(srcPath, []byte("not a valid gzip file"), 0o644)
				require.NoError(t, err)

				return srcPath, dstPath
			},
			wantErr: true,
			errMsg:  "unable to read gzip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcPath, dstPath := tt.setup(t)

			err := extractChartTgz(srcPath, dstPath)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)

				return
			}

			require.NoError(t, err)

			if tt.verify != nil {
				tt.verify(t, dstPath)
			}
		})
	}
}

func TestExtractTarGz(t *testing.T) {
	tests := []struct {
		setup   func(t *testing.T) (reader io.Reader, dstPath string)
		verify  func(t *testing.T, dstPath string)
		name    string
		errMsg  string
		wantErr bool
	}{
		{
			name: "positive case - successful extraction",
			setup: func(t *testing.T) (io.Reader, string) {
				t.Helper()

				tmpDir := t.TempDir()
				dstPath := filepath.Join(tmpDir, "extracted")

				files := map[string]string{
					"file1.txt": "content1",
					"file2.txt": "content2",
				}
				tgzContent := createTestTarGz(t, "mychart", files)

				return bytes.NewReader(tgzContent), dstPath
			},
			wantErr: false,
			verify: func(t *testing.T, dstPath string) {
				t.Helper()

				content1, err := os.ReadFile(filepath.Join(dstPath, "file1.txt"))
				require.NoError(t, err)
				require.Equal(t, "content1", string(content1))

				content2, err := os.ReadFile(filepath.Join(dstPath, "file2.txt"))
				require.NoError(t, err)
				require.Equal(t, "content2", string(content2))
			},
		},
		{
			name: "negative case - invalid gzip",
			setup: func(t *testing.T) (io.Reader, string) {
				t.Helper()

				tmpDir := t.TempDir()
				dstPath := filepath.Join(tmpDir, "extracted")

				return bytes.NewReader([]byte("invalid gzip")), dstPath
			},
			wantErr: true,
			errMsg:  "unable to read gzip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, dstPath := tt.setup(t)

			err := extractTarGz(reader, dstPath)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)

				return
			}

			require.NoError(t, err)

			if tt.verify != nil {
				tt.verify(t, dstPath)
			}
		})
	}
}

func TestExtractTar(t *testing.T) {
	tests := []struct {
		setup   func(t *testing.T) (reader io.Reader, dstPath string)
		verify  func(t *testing.T, dstPath string)
		name    string
		errMsg  string
		wantErr bool
	}{
		{
			name: "positive case - extract files and directories",
			setup: func(t *testing.T) (io.Reader, string) {
				t.Helper()

				tmpDir := t.TempDir()
				dstPath := filepath.Join(tmpDir, "extracted")

				// Create a tar archive (not gzipped)
				var buf bytes.Buffer

				tw := tar.NewWriter(&buf)

				files := map[string]struct {
					content string
					mode    int64
					isDir   bool
				}{
					"chart/":                      {isDir: true, mode: 0o755},
					"chart/Chart.yaml":            {content: "name: test", mode: 0o644},
					"chart/templates/":            {isDir: true, mode: 0o755},
					"chart/templates/deploy.yaml": {content: "kind: Deployment", mode: 0o644},
				}

				for name, file := range files {
					var header *tar.Header
					if file.isDir {
						header = &tar.Header{
							Name:     name,
							Mode:     file.mode,
							Typeflag: tar.TypeDir,
						}
					} else {
						header = &tar.Header{
							Name:     name,
							Mode:     file.mode,
							Size:     int64(len(file.content)),
							Typeflag: tar.TypeReg,
						}
					}

					err := tw.WriteHeader(header)
					require.NoError(t, err)

					if !file.isDir {
						_, err = tw.Write([]byte(file.content))
						require.NoError(t, err)
					}
				}

				err := tw.Close()
				require.NoError(t, err)

				return &buf, dstPath
			},
			wantErr: false,
			verify: func(t *testing.T, dstPath string) {
				t.Helper()

				// Verify files were extracted (without top-level directory)
				chartYaml, err := os.ReadFile(filepath.Join(dstPath, "Chart.yaml"))
				require.NoError(t, err)
				require.YAMLEq(t, "name: test", string(chartYaml))

				deployYaml, err := os.ReadFile(filepath.Join(dstPath, "templates", "deploy.yaml"))
				require.NoError(t, err)
				require.YAMLEq(t, "kind: Deployment", string(deployYaml))
			},
		},
		{
			name: "negative case - invalid tar format",
			setup: func(t *testing.T) (io.Reader, string) {
				t.Helper()

				tmpDir := t.TempDir()
				dstPath := filepath.Join(tmpDir, "extracted")

				return bytes.NewReader([]byte("not a valid tar")), dstPath
			},
			wantErr: true,
			errMsg:  "unable to read tar content",
		},
		{
			name: "positive case - skip files without subdirectory",
			setup: func(t *testing.T) (io.Reader, string) {
				t.Helper()

				tmpDir := t.TempDir()
				dstPath := filepath.Join(tmpDir, "extracted")

				// Create tar with file at root (no subdirectory)
				var buf bytes.Buffer

				tw := tar.NewWriter(&buf)

				// This file should be skipped (no subdirectory)
				header := &tar.Header{
					Name:     "rootfile.txt",
					Mode:     0o644,
					Size:     4,
					Typeflag: tar.TypeReg,
				}
				err := tw.WriteHeader(header)
				require.NoError(t, err)
				_, err = tw.Write([]byte("test"))
				require.NoError(t, err)

				// This file should be extracted
				header = &tar.Header{
					Name:     "chart/file.txt",
					Mode:     0o644,
					Size:     7,
					Typeflag: tar.TypeReg,
				}
				err = tw.WriteHeader(header)
				require.NoError(t, err)
				_, err = tw.Write([]byte("content"))
				require.NoError(t, err)

				err = tw.Close()
				require.NoError(t, err)

				return &buf, dstPath
			},
			wantErr: false,
			verify: func(t *testing.T, dstPath string) {
				t.Helper()

				// Root file should not exist
				_, err := os.Stat(filepath.Join(dstPath, "rootfile.txt"))
				require.Error(t, err)
				require.True(t, os.IsNotExist(err))

				// Chart file should exist
				content, err := os.ReadFile(filepath.Join(dstPath, "file.txt"))
				require.NoError(t, err)
				require.Equal(t, "content", string(content))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, dstPath := tt.setup(t)

			err := extractTar(reader, dstPath)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)

				return
			}

			require.NoError(t, err)

			if tt.verify != nil {
				tt.verify(t, dstPath)
			}
		})
	}
}
