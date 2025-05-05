# Temple 🕍

Temple is a minimal, concurrent static site generator that transforms Markdown notes into HTML pages. Built with Go’s excellent templating package and a couple of other packages, it offers a (limited) lightweight alternative to heavier tools like [Hugo](https://github.com/gohugoio/hugo). I use Temple as the templating engine for the notes on [my website](https://guycho.uk).

## Features

- Markdown parsing (using [goldmark](https://github.com/yuin/goldmark)) and conversion to HTML (with YAML frontmatter support)
- Syntax highlighting for code blocks using [chroma](https://github.com/alecthomas/chroma)
- Custom HTML templates using go's [html/template](https://pkg.go.dev/html/template)
- Concurrent processing of notes
- Deduplicated CSS output for the syntax highlighting

## Prerequisites

- Go 1.16 or later
- Python 3 (for serving the generated files)

## Installation

1. Clone the repository:
```shell
git clone https://github.com/guychouk/tmpl.git
cd tmpl
```

2. Build the project:
```shell
make build
```

## Usage

### Basic Usage

To generate HTML files from your Markdown notes:

```shell
./tmpl --output ./public --templates ./src/templates.html notes
```

### Using Makefile

The project includes a Makefile with several useful commands:

- `make build` - Build the project
- `make clean` - Remove build artifacts and generated files
- `make serve` - Serve the generated files using Python's built-in HTTP server
- `make all` - Build and serve in one command

### Example Workflow

1. Build the project:
```shell
make build
```

2. Generate HTML files:
```shell
./tmpl --output ./public --templates ./src/templates.html notes
```

3. Serve the generated files:
```shell
make serve
```

Then visit `http://localhost:8000` in your browser to see the generated site.

## Project Structure

```
.
├── main.go          # Main program
├── Makefile         # Build and serve commands
├── notes/           # Your Markdown notes
├── public/          # Generated HTML files
└── src/
    └── templates.html  # HTML templates
```

## How It Works

All `tmpl` does is it reads a directory with Markdown files, maps each `.md` file to a `Note` struct, reads the HTML templates defined in the provided file, and finally generates HTML files along with a CSS file for syntax highlighting in code blocks.

Each note is expected to have a [YAML frontmatter](https://pandoc.org/MANUAL.html#extension-yaml_metadata_block) with the following fields:

```md
---
title: A meaningful title
date: YYYY-MM-DD
---

your note here 
```

1. The program reads Markdown files from the specified directory
2. Each file is processed concurrently:
   - Markdown is converted to HTML
   - YAML frontmatter is extracted
   - Syntax highlighting CSS is collected
3. HTML files are generated using the provided templates
4. A single, deduplicated CSS file is created for syntax highlighting

## Dependencies

- [goldmark](https://github.com/yuin/goldmark) - Markdown parser
- [chroma](https://github.com/alecthomas/chroma) - Syntax highlighting
- [html/template](https://pkg.go.dev/html/template) - HTML templating

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
