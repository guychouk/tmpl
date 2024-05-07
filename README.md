# Temple: my templating "engine"

This is a program that I use as a templating engine for my website: [https://guycho.uk](https://guycho.uk).

It reads a directory with Markdown files (which may include a [YAML frontmatter](https://pandoc.org/MANUAL.html#extension-yaml_metadata_block)), parses these `.md` files, reads a bunch of templates defined in a `templates.html` file, and generates HTML along with a CSS file for syntax highlighting in code blocks, all done using [goldmark](https://github.com/yuin/goldmark) and [chroma](https://github.com/alecthomas/chroma).

Here's how I use it to generate the HTML for my `notes`:

```shell
./temple --output ./public --templates ./src/templates.html notes
```

This is a minimal replacement for static site generators such as [hugo](https://github.com/gohugoio/hugo).

I mainly wrote this to serve my needs and to see if I can do it. Not only was it easy, it was quite fun as well ¯\\_(ツ)_/¯.
