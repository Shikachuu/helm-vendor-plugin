# Helm Vendor Plugin

A Helm plugin for downloading and vendoring Helm charts from remote repositories.

## Installation

```bash
helm plugin install https://github.com/Shikachuu/helm-vendor-plugin
```

## Usage

### Download Charts

Download helm charts defined in your vendor-charts configuration file:

```bash
helm vendor download -f .vendor-charts.yaml
```

This command will:

- Read your vendor-charts configuration
- Download each specified helm chart from OCI or Helm repositories
- Save charts to their designated destination directories

### Verify Configuration

Verify your vendor-charts configuration file:

```bash
helm vendor verify -f .vendor-charts.yaml
```

This validates the configuration file against the expected schema without downloading any charts.

### Version Information

Print version information:

```bash
helm vendor version
```

## Configuration

The plugin supports both **YAML** and **JSON** configuration formats. Create a configuration file (typically `.vendor-charts.yaml` or `.vendor-charts.json`) in your project root.

### Configuration Format

The configuration file defines a list of Helm charts to vendor. Each chart entry specifies where to download the chart from and where to save it locally.

#### YAML Example

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/Shikachuu/helm-vendor-plugin/refs/heads/main/schema.json
charts:
  - name: descheduler
    repository: https://kubernetes-sigs.github.io/descheduler
    version: 0.34.0
    destination: artifacts/descheduler
    insecure: false

  - name: prometheus
    repository: oci://ghcr.io/prometheus-community/charts
    version: 27.45.0
    destination: artifacts/prometheus

  - name: cert-manager
    repository: https://charts.jetstack.io
    version: v1.19.1
    verify: true
    destination: artifacts/cert-manager
```

#### JSON Example

```json
{
  "$schema": "https://raw.githubusercontent.com/Shikachuu/helm-vendor-plugin/refs/heads/main/schema.json",
  "charts": [
    {
      "name": "descheduler",
      "repository": "https://kubernetes-sigs.github.io/descheduler",
      "version": "0.34.0",
      "destination": "artifacts/descheduler",
      "insecure": false
    },
    {
      "name": "prometheus",
      "repository": "oci://ghcr.io/prometheus-community/charts",
      "version": "27.45.0",
      "destination": "artifacts/prometheus"
    }
  ]
}
```

### Configuration Fields

| Field         | Required | Type    | Description                                                               |
| ------------- | -------- | ------- | ------------------------------------------------------------------------- |
| `name`        | Yes      | string  | Name of the Helm chart                                                    |
| `repository`  | Yes      | string  | Chart repository URL (supports `http://`, `https://`, or `oci://`)        |
| `version`     | Yes      | string  | Chart version to vendor                                                   |
| `destination` | Yes      | string  | Local destination path for the vendored chart                             |
| `insecure`    | No       | boolean | Allow insecure (non-TLS) connections to the repository (default: `false`) |
| `verify`      | No       | boolean | Verify chart provenance (default: `false`)                                |

### JSON Schema

The configuration is validated against a [JSON Schema](schema.json) that provides:

- IDE autocompletion and validation (when using the `$schema` directive)
- Runtime validation via the `verify` command
- Documentation of all supported fields and their constraints

You can reference the schema in your editor for autocompletion:

- **YAML**: Add `# yaml-language-server: $schema=<schema-url>` at the top
- **JSON**: Add `"$schema": "<schema-url>"` to the root object

## Development

Requirements:

- Go 1.25.3 or later
- [mise](https://mise.jdx.dev/) (optional, for managing tools)

## License

See LICENSE file for details.
