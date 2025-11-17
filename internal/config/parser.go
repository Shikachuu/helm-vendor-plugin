package config

import (
	"errors"
	"fmt"

	"github.com/kaptinlin/jsonschema"

	_ "embed"
)

// This hack is required for go embed to work, however we want to keep the schema.json in the root, for easier yaml-ls import.
//
//go:generate cp ../../schema.json schema.json
//go:embed schema.json
var jsonSchema []byte

// Parser describes a struct that should be able to parse a vendor-charts config file.
type Parser interface {
	// Validate will check the structural content of the given byte array.
	// Returns an error if any.
	Validate([]byte) error

	// Unmarshall will parse the given vendor-charts config file, after checking it's structural integrity.
	// Returns the unmarshalled config or an error if any.
	Unmarshall([]byte) ([]VendorChart, error)
}

// JSONConfigParser implements the Parser interface using JSON schema validation.
// It validates and unmarshalls vendor-charts configuration files against an embedded JSON schema.
type JSONConfigParser struct {
	schema *jsonschema.Schema
}

// Unmarshall parses the given vendor-charts configuration after validating it against the JSON schema.
// It returns the parsed VendorChart slice or an error if validation or unmarshalling fails.
func (j *JSONConfigParser) Unmarshall(cfg []byte) ([]VendorChart, error) {
	var vcs []VendorChart

	err := j.Validate(cfg)
	if err != nil {
		return nil, err
	}

	err = j.schema.Unmarshal(vcs, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall configuration: %w", err)
	}

	return vcs, nil
}

// Validate checks the structural integrity of the configuration file against the JSON schema.
// It returns a detailed error message listing all validation failures, or nil if valid.
func (j *JSONConfigParser) Validate(cfg []byte) error {
	r := j.schema.ValidateJSON(cfg)
	errMsg := "invalid configuration file:"

	if !r.IsValid() {
		for field, err := range r.Errors {
			errMsg = fmt.Sprintf("%s\n- %s: %s", errMsg, field, err.Message)
		}

		return errors.New(errMsg) //nolint:err113 // We want a dynamic error here to preserve the keys for user exp
	}

	return nil
}

// NewJSONConfigParser creates a new JSONConfigParser with the embedded JSON schema compiled and ready for use.
// It returns an error if the embedded schema cannot be compiled.
func NewJSONConfigParser() (*JSONConfigParser, error) {
	compiler := jsonschema.NewCompiler()

	s, err := compiler.Compile(jsonSchema)
	if err != nil {
		return nil, fmt.Errorf("unable to compile json schema: %w", err)
	}

	return &JSONConfigParser{schema: s}, nil
}

var _ Parser = (*JSONConfigParser)(nil)
