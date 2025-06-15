# eolfmt

<img src="https://i.gyazo.com/1dfcfa8453af6b93202170afe2e63495.jpg" alt="Logo" width="500">


A high-performance tool for fixing line endings. It recursively scans directories and adds missing newlines to the end of text files.

## Features

- **Ultra-fast**: Achieves high-speed processing through parallel execution and minimal I/O
- **Smart filtering**: Automatic binary file detection based on file extensions
- **Safe**: Automatically skips binary files
- **Efficient**: Efficient parallel processing using the Sourcegraph conc library

## Installation

```bash
go install github.com/upamune/eolfmt/cmd/eolfmt@latest
```

Or build from source:

```bash
git clone https://github.com/upamune/eolfmt.git
cd eolfmt
go build -o eolfmt ./cmd/eolfmt
```

## Usage

```bash
# Process all files in the current directory
eolfmt

# Process specific directories
eolfmt /path/to/directory

# Process multiple paths
eolfmt src/ docs/ tests/

# Show version information
eolfmt -version
```

## Excluded Files and Directories

### Automatically excluded directories
- `.git`
- `node_modules`
- `.idea`
- `target`
- `build`
- `dist`
- `vendor`
- `__pycache__`

### Automatically skipped extensions (binary files)
- Executables: `.exe`, `.dll`, `.so`, `.dylib`
- Images: `.png`, `.jpg`, `.jpeg`, `.gif`
- Archives: `.pdf`, `.zip`, `.tar`, `.gz`
- Compiled: `.pyc`, `.class`, `.jar`

## Performance

| Project Size | File Count | Processing Time |
|--------------|------------|-----------------|
| Small        | 1,000      | < 0.5s          |
| Medium       | 10,000     | < 2s            |
| Large        | 100,000    | < 15s           |

## Development

```bash
# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Static analysis (revive)
make lint

# Run all checks
make check
```

## Known Issues

- **Hard link support**: The tool does not currently handle hard links specially. When processing hard-linked files, changes will affect all linked instances.

## License

MIT
