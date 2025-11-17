# Helm Vendor Plugin

A Helm plugin for downloading and vendoring Helm charts from remote repositories.

## Installation

```bash
helm plugin install https://github.com/Shikachuu/helm-vendor-plugin
```

## Usage

### Verify Configuration

Verify your vendor-charts configuration file:

```bash
helm vendor verify -f .vendor-charts.yaml
```

## Configuration

Create a `.vendor-charts.yaml` file in your project root. See `schema.json` for the configuration schema.

## Development

Requirements:
- Go 1.25.3 or later
- [mise](https://mise.jdx.dev/) (optional, for managing tools)

### Build

```bash
go build -o bin/vendor-plugin
```

Or with mise:

```bash
mise run build
```

### Test

```bash
go test ./...
```

Or with mise:

```bash
mise run test
```

### Lint

```bash
mise run lint
```

## License

See LICENSE file for details.
