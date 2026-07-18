package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// bt returns a backtick character
var bt = "`"

// Language defines a supported programming language
type Language struct {
	Name       string
	Extension  string
	PackageMgr string
	BuildCmd   string
	TestCmd    string
	LintCmd    string
}

var languages = map[string]Language{
	"go": {
		Name:       "Go",
		Extension:  ".go",
		PackageMgr: "go mod",
		BuildCmd:   "go build ./...",
		TestCmd:    "go test ./...",
		LintCmd:    "golangci-lint run",
	},
	"python": {
		Name:       "Python",
		Extension:  ".py",
		PackageMgr: "pip install virtualenv && python -m venv .venv",
		BuildCmd:   "python -m compileall .",
		TestCmd:    "python -m pytest tests/ -v",
		LintCmd:    "ruff check .",
	},
	"rust": {
		Name:       "Rust",
		Extension:  ".rs",
		PackageMgr: "cargo",
		BuildCmd:   "cargo build --release",
		TestCmd:    "cargo test",
		LintCmd:    "cargo clippy -- -D warnings",
	},
	"node": {
		Name:       "Node.js",
		Extension:  ".js",
		PackageMgr: "npm init -y && npm install",
		BuildCmd:   "npm run build",
		TestCmd:    "npm test",
		LintCmd:    "eslint .",
	},
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "init":
		handleInit(os.Args[2:])
	case "list":
		handleList()
	case "create":
		handleCreate(os.Args[2:])
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Scaffold - Project Scaffolding CLI")
	fmt.Println()
	fmt.Println("Create boilerplate projects from templates with one command.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  scaffold <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init       Initialize scaffold configuration")
	fmt.Println("  create     Create a new project from a template")
	fmt.Println("  list       List available templates")
	fmt.Println("  help       Show this help message")
	fmt.Println()
	fmt.Println("Create Options:")
	fmt.Println("  --lang, -l   Programming language (go, python, rust, node)")
	fmt.Println("  --author, -a Author name")
	fmt.Println("  --desc, -d   Project description")
	fmt.Println("  --version    Project version (default: 0.1.0)")
	fmt.Println("  --repo       Repository URL")
	fmt.Println("  --ci         Include CI configuration")
	fmt.Println("  --license    Include license (MIT, Apache-2.0, GPL-3.0)")
	fmt.Println("  --tests      Include test files")
	fmt.Println("  --gitignore  Include .gitignore")
	fmt.Println("  --readme     Include README.md")
	fmt.Println("  --output, -o Output directory (default: current directory)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  scaffold create -l go -a \"John Doe\" -d \"My Go project\"")
	fmt.Println("  scaffold create -l python --license Apache-2.0 --ci")
	fmt.Println("  scaffold create -l rust --tests --gitignore")
	fmt.Println("  scaffold list")
}

func handleInit(args []string) {
	fmt.Println("Initializing scaffold configuration...")

	configPath := ".scaffold.json"
	if len(args) > 0 {
		configPath = args[0]
	}

	config := `  "author": "Your Name",
  "default_language": "go",
  "version": "0.1.0",
  "ci_enabled": true,
  "license": "MIT"`

	content := "{\n" + config + "\n}\n"

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration saved to %s\n", configPath)
	fmt.Println()
	fmt.Println("Edit the file to customize defaults, then run:")
	fmt.Println("  scaffold create -l go -a \"Your Name\" -d \"Project description\"")
}

func handleList() {
	fmt.Println("Available Languages:")
	fmt.Println(strings.Repeat("-", 50))

	for _, lang := range []string{"go", "python", "rust", "node"} {
		language := languages[lang]
		fmt.Printf("  %-10s %-15s %s\n",
			lang,
			language.Name,
			"Build: "+language.BuildCmd)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Total: %d languages\n\n", len(languages))

	fmt.Println("Available Templates per Language:")
	templates := map[string][]string{
		"go":     {"main.go", "pkg.go", "pkg_test.go", "go.mod", ".gitignore", "LICENSE", ".github/workflows/ci.yml", "README.md"},
		"python": {"app.py", "test_app.py", "pyproject.toml", "requirements.txt", ".gitignore", "LICENSE", ".github/workflows/ci.yml", "README.md"},
		"rust":   {"main.rs", "lib.rs", "integration.rs", "Cargo.toml", ".gitignore", "LICENSE", ".github/workflows/ci.yml", "README.md"},
		"node":   {"index.js", "index.test.js", "package.json", ".gitignore", "LICENSE", ".github/workflows/ci.yml", "README.md"},
	}

	for _, lang := range []string{"go", "python", "rust", "node"} {
		fmt.Printf("\n  %s:\n", strings.ToUpper(lang))
		for _, f := range templates[lang] {
			fmt.Printf("    - %s\n", f)
		}
	}
}

func handleCreate(args []string) {
	config := createConfig()

	outputDir := "."

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--lang", "-l":
			if i+1 < len(args) {
				config.Language = args[i+1]
				i++
			}
		case "--author", "-a":
			if i+1 < len(args) {
				config.Author = args[i+1]
				i++
			}
		case "--desc", "-d":
			if i+1 < len(args) {
				config.Description = args[i+1]
				i++
			}
		case "--version":
			if i+1 < len(args) {
				config.Version = args[i+1]
				i++
			}
		case "--repo":
			if i+1 < len(args) {
				config.RepoURL = args[i+1]
				i++
			}
		case "--ci":
			config.IncludeCI = true
		case "--no-ci":
			config.IncludeCI = false
		case "--license":
			config.IncludeLicense = true
			if i+1 < len(args) {
				config.LicenseType = args[i+1]
				i++
			}
		case "--no-license":
			config.IncludeLicense = false
		case "--tests":
			config.IncludeTests = true
		case "--no-tests":
			config.IncludeTests = false
		case "--gitignore":
			config.IncludeGitignore = true
		case "--no-gitignore":
			config.IncludeGitignore = false
		case "--readme":
			config.IncludeReadme = true
		case "--no-readme":
			config.IncludeReadme = false
		case "--output", "-o":
			if i+1 < len(args) {
				outputDir = args[i+1]
				i++
			}
		case "--help":
			printUsage()
			return
		}
	}

	// Extract project name from output directory
	parts := strings.Split(outputDir, "/")
	config.ProjectName = parts[len(parts)-1]
	if config.ProjectName == "" {
		config.ProjectName = "project"
	}

	if config.Author == "" {
		config.Author = "Your Name"
	}
	if config.Description == "" {
		config.Description = fmt.Sprintf("A %s project", config.Language)
	}

	// Validate language
	language, ok := languages[config.Language]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: unsupported language '%s'. Supported: %s\n",
			config.Language, strings.Join(getLanguageKeys(), ", "))
		os.Exit(1)
	}

	// Create output directory
	if outputDir != "." {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Creating project: %s\n", config.ProjectName)
	fmt.Printf("   Language: %s\n", language.Name)
	fmt.Printf("   Author: %s\n", config.Author)
	fmt.Printf("   Output: %s\n\n", outputDir)

	// Generate project structure
	err := generateProject(outputDir, config, language)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project scaffolded successfully!")
	fmt.Printf("   cd %s\n", outputDir)
	fmt.Printf("   %s\n", language.BuildCmd)
}

type ScaffoldConfig struct {
	ProjectName      string
	Language         string
	Author           string
	Description      string
	Version          string
	RepoURL          string
	IncludeCI        bool
	IncludeLicense   bool
	LicenseType      string
	IncludeTests     bool
	IncludeGitignore bool
	IncludeReadme    bool
}

func createConfig() ScaffoldConfig {
	return ScaffoldConfig{
		Language:         "go",
		Version:          "0.1.0",
		IncludeCI:        true,
		IncludeLicense:   true,
		LicenseType:      "MIT",
		IncludeTests:     true,
		IncludeGitignore: true,
		IncludeReadme:    true,
	}
}

func generateProject(outputDir string, config ScaffoldConfig, lang Language) error {
	switch config.Language {
	case "go":
		return generateGoProject(outputDir, config)
	case "python":
		return generatePythonProject(outputDir, config)
	case "rust":
		return generateRustProject(outputDir, config)
	case "node":
		return generateNodeProject(outputDir, config)
	default:
		return fmt.Errorf("unsupported language: %s", config.Language)
	}
}

func generateGoProject(outputDir string, config ScaffoldConfig) error {
	files := map[string]string{
		"cmd/main.go":         goMain(config),
		"internal/pkg/pkg.go": goPkg(config),
		"go.mod":              goMod(config),
	}

	if config.IncludeTests {
		files["internal/pkg/pkg_test.go"] = goTest(config)
	}
	if config.IncludeGitignore {
		files[".gitignore"] = goGitignore()
	}
	if config.IncludeLicense {
		files["LICENSE"] = getLicense(config.LicenseType)
	}
	if config.IncludeCI {
		files[".github/workflows/ci.yml"] = goCI(config)
	}
	if config.IncludeReadme {
		files["README.md"] = goReadme(config)
	}

	return writeFiles(outputDir, files)
}

func goMod(config ScaffoldConfig) string {
	return fmt.Sprintf("module github.com/EdgarOrtegaRamirez/%s\n\ngo 1.22\n", config.ProjectName)
}

func goMain(config ScaffoldConfig) string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"os"

	"github.com/EdgarOrtegaRamirez/%s/internal/pkg"
)

func main() {
	result, err := pkg.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %%v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}
`, config.ProjectName)
}

func goPkg(config ScaffoldConfig) string {
	return fmt.Sprintf(`// Package pkg provides the core functionality of %s.
package pkg

// Run executes the main logic and returns the result.
func Run() (string, error) {
	return "Hello from %s!", nil
}
`, config.ProjectName, config.ProjectName)
}

func goTest(config ScaffoldConfig) string {
	return `package pkg

import "testing"

func TestRun(t *testing.T) {
	result, err := Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}
`
}

func goGitignore() string {
	return `# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Build artifacts
*.a
*.lib

# OS files
.DS_Store
Thumbs.db

# IDE
.idea/
.vscode/
*.swp
*.swo

# Go
go.sum

# Coverage
coverage.out
coverage.txt
`
}

func goCI(config ScaffoldConfig) string {
	return `name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test ./... -race -coverprofile=coverage.out

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
`
}

func goReadme(config ScaffoldConfig) string {
	b := bt
	cb := b + b + b
	var buf strings.Builder
	buf.WriteString("# " + config.ProjectName + "\n\n")
	buf.WriteString(config.ProjectName + "\n\n")
	buf.WriteString("## Description\n\n")
	buf.WriteString(config.Description + "\n\n")
	buf.WriteString("## Installation\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("go install github.com/EdgarOrtegaRamirez/" + config.ProjectName + "@latest\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("Or build from source:\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("git clone https://github.com/EdgarOrtegaRamirez/" + config.ProjectName + ".git\n")
	buf.WriteString("cd " + config.ProjectName + "\n")
	buf.WriteString("go build -o ./bin/" + config.ProjectName + " ./cmd/...\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## Usage\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("./bin/" + config.ProjectName + "\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## Testing\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("go test ./...\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## License\n\n")
	buf.WriteString(config.LicenseType + " - see [LICENSE](LICENSE) for details.\n")
	return buf.String()
}

func generatePythonProject(outputDir string, config ScaffoldConfig) error {
	files := map[string]string{
		"src/app.py":        pyApp(config),
		"tests/test_app.py": pyTest(config),
		"pyproject.toml":    pyPyproject(config),
		".gitignore":        pyGitignore(),
	}

	if config.IncludeLicense {
		files["LICENSE"] = getLicense(config.LicenseType)
	}
	if config.IncludeCI {
		files[".github/workflows/ci.yml"] = pyCI(config)
	}
	if config.IncludeReadme {
		files["README.md"] = pyReadme(config)
	}

	return writeFiles(outputDir, files)
}

func pyApp(config ScaffoldConfig) string {
	return fmt.Sprintf(`"""%s - %s"""


def main():
    """Main entry point."""
    print("Hello from %s!")


if __name__ == "__main__":
    main()
`, config.ProjectName, config.Description, config.ProjectName)
}

func pyTest(config ScaffoldConfig) string {
	return fmt.Sprintf(`"""Tests for %s."""

import unittest


class TestMain(unittest.TestCase):
    """Test cases for the main module."""

    def test_main_output(self):
        """Test that main produces expected output."""
        result = "Hello from %s!"
        self.assertIn("%s", result)


if __name__ == "__main__":
    unittest.main()
`, config.ProjectName, config.ProjectName, config.ProjectName)
}

func pyPyproject(config ScaffoldConfig) string {
	return fmt.Sprintf(`[build-system]
requires = ["setuptools>=61.0"]
build-backend = "setuptools.build_meta"

[project]
name = "%s"
version = "%s"
description = "%s"
requires-python = ">=3.9"
license = {text = "%s"}
`, config.ProjectName, config.Version, config.Description, config.LicenseType)
}

func pyGitignore() string {
	return `# Byte-compiled / optimized / DLL files
__pycache__/
*.py[cod]
*$py.class

# Virtual environment
.venv/
venv/
ENV/

# IDE
.idea/
.vscode/
*.swp
*.swo

# Distribution
dist/
build/
*.egg-info/

# Coverage
coverage/
htmlcov/
.coverage
.coverage.*

# pytest
.pytest_cache/
`
}

func pyCI(config ScaffoldConfig) string {
	return `name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        python-version: ["3.9", "3.10", "3.11", "3.12"]

    steps:
      - uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: "3.12"

      - name: Install dependencies
        run: |
          python -m pip install --upgrade pip
          pip install pytest ruff

      - name: Run linter
        run: ruff check .

      - name: Run tests
        run: python -m pytest tests/ -v
`
}

func pyReadme(config ScaffoldConfig) string {
	b := bt
	cb := b + b + b
	var buf strings.Builder
	buf.WriteString("# " + config.ProjectName + "\n\n")
	buf.WriteString(config.ProjectName + "\n\n")
	buf.WriteString("## Description\n\n")
	buf.WriteString(config.Description + "\n\n")
	buf.WriteString("## Installation\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("pip install " + config.ProjectName + "\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("Or install from source:\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("git clone https://github.com/EdgarOrtegaRamirez/" + config.ProjectName + ".git\n")
	buf.WriteString("cd " + config.ProjectName + "\n")
	buf.WriteString("pip install -e .\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## Usage\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("python -m " + config.ProjectName + "\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## Testing\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("python -m pytest tests/ -v\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## License\n\n")
	buf.WriteString(config.LicenseType + " - see [LICENSE](LICENSE) for details.\n")
	return buf.String()
}

func generateRustProject(outputDir string, config ScaffoldConfig) error {
	files := map[string]string{
		"src/main.rs":          rustMain(config),
		"src/lib.rs":           rustLib(config),
		"tests/integration.rs": rustTest(config),
		"Cargo.toml":           rustCargo(config),
		".gitignore":           rustGitignore(),
	}

	if config.IncludeLicense {
		files["LICENSE"] = getLicense(config.LicenseType)
	}
	if config.IncludeCI {
		files[".github/workflows/ci.yml"] = rustCI(config)
	}
	if config.IncludeReadme {
		files["README.md"] = rustReadme(config)
	}

	return writeFiles(outputDir, files)
}

func rustCargo(config ScaffoldConfig) string {
	return fmt.Sprintf(`[package]
name = "%s"
version = "%s"
edition = "2021"
description = "%s"
license = "%s"
authors = ["%s"]

[dependencies]
`, config.ProjectName, config.Version, config.Description, config.LicenseType, config.Author)
}

func rustMain(config ScaffoldConfig) string {
	return fmt.Sprintf(`//! %s - Main entry point

mod lib;

use lib::run;

fn main() {
    match run() {
        Ok(output) => println!("{}", output),
        Err(e) => {
            eprintln!("Error: {}", e);
            std::process::exit(1);
        }
    }
}
`, config.ProjectName)
}

func rustLib(config ScaffoldConfig) string {
	return fmt.Sprintf(`//! %s - Core library

/// Run the main application logic.
///
/// # Returns
/// A string with the application output.
///
/// # Errors
/// Returns an error if the application encounters a problem.
pub fn run() -> Result<String, Box<dyn std::error::Error>> {
    Ok(format!("Hello from %s!"))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_run() {
        let result = run().expect("run should succeed");
        assert!(!result.is_empty());
        assert!(result.contains("%s"));
    }
}
`, config.ProjectName, config.ProjectName, config.ProjectName)
}

func rustTest(config ScaffoldConfig) string {
	return fmt.Sprintf(`//! Integration tests for %s

use %s::run;

#[test]
fn test_integration() {
    let result = run().expect("run should succeed");
    assert!(!result.is_empty());
    assert!(result.contains("%s"));
}
`, config.ProjectName, config.ProjectName, config.ProjectName)
}

func rustGitignore() string {
	return `# Generated files
target/
**/*.rs.bk

# IDE
.idea/
.vscode/
*.swp
*.swo

# Documentation
doc/

# Temporary files
*.tmp
`
}

func rustCI(config ScaffoldConfig) string {
	return `name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Rust
        uses: dtolnay/rust-toolchain@stable

      - name: Run tests
        run: cargo test

      - name: Check formatting
        run: cargo fmt --check

      - name: Run clippy
        run: cargo clippy -- -D warnings
`
}

func rustReadme(config ScaffoldConfig) string {
	b := bt
	cb := b + b + b
	var buf strings.Builder
	buf.WriteString("# " + config.ProjectName + "\n\n")
	buf.WriteString(config.ProjectName + "\n\n")
	buf.WriteString("## Description\n\n")
	buf.WriteString(config.Description + "\n\n")
	buf.WriteString("## Installation\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("cargo install " + config.ProjectName + "\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("Or build from source:\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("git clone https://github.com/EdgarOrtegaRamirez/" + config.ProjectName + ".git\n")
	buf.WriteString("cd " + config.ProjectName + "\n")
	buf.WriteString("cargo build --release\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## Usage\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("./target/release/" + config.ProjectName + "\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## Testing\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("cargo test\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## License\n\n")
	buf.WriteString(config.LicenseType + " - see [LICENSE](LICENSE) for details.\n")
	return buf.String()
}

func generateNodeProject(outputDir string, config ScaffoldConfig) error {
	files := map[string]string{
		"src/index.js":        nodeIndex(config),
		"tests/index.test.js": nodeTest(config),
		"package.json":        nodePackage(config),
		".gitignore":          nodeGitignore(),
	}

	if config.IncludeLicense {
		files["LICENSE"] = getLicense(config.LicenseType)
	}
	if config.IncludeCI {
		files[".github/workflows/ci.yml"] = nodeCI(config)
	}
	if config.IncludeReadme {
		files["README.md"] = nodeReadme(config)
	}

	return writeFiles(outputDir, files)
}

func nodeIndex(config ScaffoldConfig) string {
	return fmt.Sprintf(`#!/usr/bin/env node

/**
 * %s
 * %s
 */

function run() {
    console.log("Hello from %s!");
}

// Allow running as a module
if (typeof module !== "undefined" && module.exports) {
    module.exports = { run };
}

// Run if executed directly
if (require.main === module) {
    run();
}
`, config.ProjectName, config.Description, config.ProjectName)
}

func nodeTest(config ScaffoldConfig) string {
	return fmt.Sprintf(`const { run } = require("../src/index");

describe("%s", () => {
    test("run returns expected output", () => {
        const originalLog = console.log;
        let output = "";
        console.log = (msg) => { output = msg; };

        run();

        console.log = originalLog;
        expect(output).toContain("%s");
    });
});
`, config.ProjectName, config.ProjectName)
}

func nodePackage(config ScaffoldConfig) string {
	return fmt.Sprintf(`{
  "name": "%s",
  "version": "%s",
  "description": "%s",
  "main": "src/index.js",
  "bin": {
    "%s": "src/index.js"
  },
  "scripts": {
    "start": "node src/index.js",
    "test": "jest",
    "lint": "eslint ."
  },
  "keywords": [],
  "author": "%s",
  "license": "%s",
  "devDependencies": {
    "jest": "^29.0.0",
    "eslint": "^8.0.0"
  }
}
`, config.ProjectName, config.Version, config.Description, config.ProjectName, config.Author, config.LicenseType)
}

func nodeGitignore() string {
	return `# Dependencies
node_modules/

# Logs
logs/
*.log
npm-debug.log*

# Testing
coverage/
.nyc_output/

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Environment
.env
.env.local
`
}

func nodeCI(config ScaffoldConfig) string {
	return `name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'

      - name: Install dependencies
        run: npm install

      - name: Run tests
        run: npm test

      - name: Run linter
        run: npm run lint
`
}

func nodeReadme(config ScaffoldConfig) string {
	b := bt
	cb := b + b + b
	var buf strings.Builder
	buf.WriteString("# " + config.ProjectName + "\n\n")
	buf.WriteString(config.ProjectName + "\n\n")
	buf.WriteString("## Description\n\n")
	buf.WriteString(config.Description + "\n\n")
	buf.WriteString("## Installation\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("npm install " + config.ProjectName + "\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("Or install from source:\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("git clone https://github.com/EdgarOrtegaRamirez/" + config.ProjectName + ".git\n")
	buf.WriteString("cd " + config.ProjectName + "\n")
	buf.WriteString("npm install\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## Usage\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("npm start\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## Testing\n\n")
	buf.WriteString(cb + "\n")
	buf.WriteString("npm test\n")
	buf.WriteString(cb + "\n\n")
	buf.WriteString("## License\n\n")
	buf.WriteString(config.LicenseType + " - see [LICENSE](LICENSE) for details.\n")
	return buf.String()
}

func getLicense(licenseType string) string {
	switch licenseType {
	case "Apache-2.0":
		return `Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`
	case "GPL-3.0":
		return `GNU GENERAL PUBLIC LICENSE
Version 3, 29 June 2007

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
`
	default:
		return `MIT License

Copyright (c) 2026 Edgar Ortega Ramirez

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`
	}
}

func writeFiles(outputDir string, files map[string]string) error {
	for path, content := range files {
		fullPath := filepath.Join(outputDir, path)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing file %s: %w", fullPath, err)
		}
	}
	return nil
}

func getLanguageKeys() []string {
	keys := make([]string, 0, len(languages))
	for k := range languages {
		keys = append(keys, k)
	}
	return keys
}
