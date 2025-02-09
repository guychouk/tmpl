# guychouk's templating "engine"

This is a program that I use as a templating "engine" for the notes on [my website](https://guycho.uk).

All `tmpl` does is it reads a directory with Markdown files (which may include a [YAML frontmatter](https://pandoc.org/MANUAL.html#extension-yaml_metadata_block)), maps each `.md` file to a `Note` struct, reads the HTML templates defined in the provided file, and finally generates HTML files along with a CSS file for syntax highlighting in code blocks.

Here's how I use it to generate the HTML for my `notes`:

```shell
./tmpl --output ./public --templates ./src/templates.html notes
```

This is a minimal replacement for static site generators such as [hugo](https://github.com/gohugoio/hugo).

I sat down to write this thinking it's gonna be difficult, but thanks to go's excellent [html/template pkg](https://pkg.go.dev/html/template), and packages like [goldmark](https://github.com/yuin/goldmark) and [chroma](https://github.com/alecthomas/chroma), it was too damn easy ¯\\_(ツ)_/¯.

read the (ugly) code and see for yourself.
