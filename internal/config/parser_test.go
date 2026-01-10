package config

import (
	"io"
	"os"
	"testing"

	"github.com/kaptinlin/jsonschema"
	"github.com/stretchr/testify/require"
)

// readTestFile is a test helper function that helps us load json schemas from the `testdata` folder
// Returns the content of the file as a `[]byte`, fails the test in case of error.
func readTestFile(t *testing.T, f string) []byte {
	t.Helper()

	testSchema, err := os.Open(f)
	require.NoError(t, err)

	s, err := io.ReadAll(testSchema)
	require.NoError(t, err)

	return s
}

// More like and integration test, this depends on the `jsonSchema` global
func TestNewJSONConfigParser(t *testing.T) {
	tests := []struct {
		name           string
		schemaFileName string
		wantErr        bool
	}{
		{
			name:           "positive case",
			schemaFileName: "schema.json",
			wantErr:        false,
		},
		{
			name:           "negative case",
			schemaFileName: "",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		jsonSchema = []byte{}
		t.Run(tt.name, func(t *testing.T) {
			if tt.schemaFileName != "" {
				jsonSchema = readTestFile(t, tt.schemaFileName)
			}

			_, gotErr := NewJSONConfigParser()
			if tt.wantErr {
				require.Error(t, gotErr, "expected error, got nil")
				return
			}

			require.NoError(t, gotErr)

		})
	}
}

func TestJSONConfigParser_Unmarshall(t *testing.T) {
	tests := []struct {
		name    string
		cfg     []byte
		want    []VendorChart
		wantErr bool
	}{
		{
			name: "positive case",
			cfg:  []byte(`{"charts": [{"name": "traefik","repository": "oci://ghcr.io/traefik/helm","version": "37.4.0","destination": "artifacts/traefik"}]}`),
			want: []VendorChart{
				{
					Name:        "traefik",
					Repository:  "oci://ghcr.io/traefik/helm",
					Version:     "37.4.0",
					Destination: "traefik",
					Insecure:    false,
					Verify:      false,
					Extract:     false,
				},
			},
			wantErr: false,
		},
		{
			name:    "negative case",
			cfg:     []byte(`{}`),
			want:    []VendorChart{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			schema := readTestFile(t, "schema.json")

			compiler := jsonschema.NewCompiler()

			s, err := compiler.Compile(schema)
			require.NoError(t, err)

			j := JSONConfigParser{schema: s}
			gotErr := j.Validate(tt.cfg)

			if tt.wantErr {
				require.Error(t, gotErr)
				return
			}

			require.NoError(t, gotErr)
		})
	}
}
