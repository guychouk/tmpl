.PHONY: build clean serve

# Build the project
build:
	go build -o tmpl

# Clean build artifacts
clean:
	rm -f tmpl
	rm -rf public/*

# Serve the generated files
serve:
	cd public && python3 -m http.server 8000

# Build and serve
all: build serve 