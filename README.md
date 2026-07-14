# Scaffold 🏗️

A project scaffolding CLI that creates boilerplate projects from templates with one command. Scaffold supports Go, Python, Rust, and Node.js projects out of the box.

## Features

- **4 language templates**: Go, Python, Rust, Node.js
- **Complete project structure**: Source code, tests, CI/CD, LICENSE, README
- **Configurable options**: CI, tests, gitignore, license type, and more
- **Template listing**: See all available templates per language
- **Configuration file**: Save defaults with `scaffold init`

## Installation

```bash
go install github.com/EdgarOrtegaRamirez/scaffold@latest
```

Or build from source:

```bash
git clone https://github.com/EdgarOrtegaRamirez/scaffold.git
cd scaffold
go build -o scaffold ./cmd/...
```

## Usage

```bash
# Create a new Go project
scaffold create -l go -a "John Doe" -d "My Go project"

# Create a Python project with Apache license and CI
scaffold create -l python --license Apache-2.0 --ci

# Create a Rust project with tests and gitignore
scaffold create -l rust --tests --gitignore

# List available templates
scaffold list

# Initialize configuration file
scaffold init
```

## Commands

| Command | Description |
|---------|-------------|
| `init` | Initialize scaffold configuration file |
| `create` | Create a new project from a template |
| `list` | List available templates per language |
| `help` | Show help message |

## Create Options

| Option | Description |
|--------|-------------|
| `--lang, -l` | Programming language (go, python, rust, node) |
| `--author, -a` | Author name |
| `--desc, -d` | Project description |
| `--version` | Project version (default: 0.1.0) |
| `--ci` | Include CI configuration |
| `--license` | Include license (MIT, Apache-2.0, GPL-3.0) |
| `--tests` | Include test files |
| `--gitignore` | Include .gitignore |
| `--readme` | Include README.md |
| `--output, -o` | Output directory (default: current directory) |

## Generated Project Structure

### Go
```
project/
├── cmd/
│   └── main.go
├── internal/
│   └── pkg/
│       ├── pkg.go
│       └── pkg_test.go
├── .github/workflows/ci.yml
├── .gitignore
├── go.mod
├── LICENSE
└── README.md
```

### Python
```
project/
├── src/
│   └── app.py
├── tests/
│   └── test_app.py
├── .github/workflows/ci.yml
├── .gitignore
├── LICENSE
├── pyproject.toml
└── README.md
```

### Rust
```
project/
├── src/
│   ├── main.rs
│   └── lib.rs
├── tests/
│   └── integration.rs
├── .github/workflows/ci.yml
├── .gitignore
├── Cargo.toml
├── LICENSE
└── README.md
```

### Node.js
```
project/
├── src/
│   └── index.js
├── tests/
│   └── index.test.js
├── .github/workflows/ci.yml
├── .gitignore
├── LICENSE
├── package.json
└── README.md
```

## License

MIT - see [LICENSE](LICENSE) for details.
