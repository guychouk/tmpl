# Temple 🕍

A minimal, concurrent static site generator that transforms your Markdown notes into beautiful HTML pages. Built with Go's excellent templating and markdown packages, Temple provides a lightweight alternative to heavy static site generators. This is what I use as a templating "engine" for the notes on [my website](https://guycho.uk).

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

All `tmpl` does is it reads a directory with Markdown files (which may include a [YAML frontmatter](https://pandoc.org/MANUAL.html#extension-yaml_metadata_block)), maps each `.md` file to a `Note` struct, reads the HTML templates defined in the provided file, and finally generates HTML files along with a CSS file for syntax highlighting in code blocks.

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

---

This is a minimal replacement for static site generators such as [hugo](https://github.com/gohugoio/hugo).

I sat down to write this thinking it's gonna be difficult, but thanks to go's excellent [html/template pkg](https://pkg.go.dev/html/template), and packages like [goldmark](https://github.com/yuin/goldmark) and [chroma](https://github.com/alecthomas/chroma), it was too damn easy ¯\\_(ツ)_/¯.

Read the code and see for yourself.
