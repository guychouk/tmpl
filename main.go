package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type Config struct {
	OutputDir      string
	TemplatesFile  string
	NotesDir       string
	WordsPerMinute int
	MaxWorkers     int
}

type TemplateError struct {
	Op  string
	Err error
}

type Note struct {
	Name        string
	Content     template.HTML
	Title       string
	Date        string
	ReadingTime int
}

func (e *TemplateError) Error() string {
	return fmt.Sprintf("template operation %s failed: %v", e.Op, e.Err)
}

func FormatDate(dateString string) (string, error) {
	inputFormat := "2006-01-02"
	outputFormat := "01/02/06"
	parsedTime, err := time.Parse(inputFormat, dateString)
	if err != nil {
		return "", err
	}
	return parsedTime.Format(outputFormat), nil
}

func CalculateReadingTime(buf *bytes.Buffer) int {
	wordsPerMinute := 200
	words := len(strings.Fields(buf.String()))
	readingTime := float64(words) / float64(wordsPerMinute)
	return int(math.Ceil(readingTime))
}

type CSSCollector struct {
	mu    sync.Mutex
	seen  map[string]struct{}
	rules []string
}

func NewCSSCollector() *CSSCollector {
	return &CSSCollector{
		seen: make(map[string]struct{}),
	}
}

// Write implements io.Writer and collects unique CSS blocks
func (c *CSSCollector) Write(p []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	css := string(p)
	// Split by closing brace, which ends a CSS rule
	blocks := strings.Split(css, "}")
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		block = block + "}" // Add the closing brace back
		if _, exists := c.seen[block]; !exists {
			c.seen[block] = struct{}{}
			c.rules = append(c.rules, block)
		}
	}
	return len(p), nil
}

// CSS returns the deduplicated CSS as a string, preserving order
func (c *CSSCollector) CSS() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return strings.Join(c.rules, "\n\n")
}

func (c *CSSCollector) WriteToFile(outputDir string) error {
	cssPath := filepath.Join(outputDir, "highlight.css")
	return os.WriteFile(cssPath, []byte(c.CSS()), 0666)
}

func EnsureDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return err
		}
		fmt.Println("Directory created:", dirName)
	} else {
		fmt.Println("Directory already exists:", dirName)
	}
	return nil
}

func DefaultConfig() Config {
	return Config{
		OutputDir:     "./public",
		TemplatesFile: "./templates.html",
		MaxWorkers:    5,
	}
}

func (c Config) Validate() error {
	if c.OutputDir == "" {
		return fmt.Errorf("output directory cannot be empty")
	}
	if c.TemplatesFile == "" {
		return fmt.Errorf("templates file cannot be empty")
	}
	if c.NotesDir == "" {
		return fmt.Errorf("notes directory cannot be empty")
	}
	if c.MaxWorkers < 1 {
		return fmt.Errorf("workers must be greater than 0")
	}
	return nil
}

func ParseConfig() (Config, error) {
	cfg := DefaultConfig()
	flag.StringVar(&cfg.OutputDir, "output", cfg.OutputDir, "Output directory for the HTML files")
	flag.StringVar(&cfg.TemplatesFile, "templates", cfg.TemplatesFile, "A file that contains all of the templates")
	flag.IntVar(&cfg.MaxWorkers, "workers", cfg.MaxWorkers, "Maximum number of concurrent workers")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return cfg, fmt.Errorf("missing required notes directory argument")
	}
	cfg.NotesDir = args[0]
	if err := cfg.Validate(); err != nil {
		return cfg, fmt.Errorf("invalid configuration: %w", err)
	}
	return cfg, nil
}

func ProcessNote(file os.DirEntry, notesDir string, markdown goldmark.Markdown) (Note, error) {
	var buf bytes.Buffer
	noteMd, err := os.ReadFile(filepath.Join(notesDir, file.Name()))
	if err != nil {
		return Note{}, fmt.Errorf("reading file: %w", err)
	}
	context := parser.NewContext()
	if err := markdown.Convert(noteMd, &buf, parser.WithContext(context)); err != nil {
		return Note{}, fmt.Errorf("converting markdown: %w", err)
	}
	metaData := meta.Get(context)
	noteDate, err := FormatDate(metaData["date"].(string))
	if err != nil {
		return Note{}, fmt.Errorf("formatting date: %w", err)
	}
	note := Note{
		Name:        strings.TrimSuffix(file.Name(), ".md"),
		Content:     template.HTML(buf.String()),
		Title:       metaData["title"].(string),
		Date:        noteDate,
		ReadingTime: CalculateReadingTime(&buf),
	}
	return note, nil
}

func main() {
	cfg, err := ParseConfig()

	if err != nil {
		log.Fatalf("Error parsing configuration: %v", err)
	}
	if err := EnsureDir(cfg.OutputDir); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}
	files, err := os.ReadDir(cfg.NotesDir)
	if err != nil {
		log.Fatalf("Error reading notes directory: %v", err)
	}
	templates, err := template.ParseFiles(cfg.TemplatesFile)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	cssCollector := NewCSSCollector()

	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("onedark"),
				highlighting.WithCSSWriter(cssCollector),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
					chromahtml.WithLineNumbers(false),
				),
			),
		),
	)

	var notes []Note
	var wg sync.WaitGroup
	var notesMutex sync.Mutex

	for _, file := range files {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()
			note, err := ProcessNote(file, cfg.NotesDir, md)
			if err != nil {
				log.Printf("Error processing %s: %v", file.Name(), err)
				return
			}
			notesMutex.Lock()
			notes = append(notes, note)
			notesMutex.Unlock()
		}(file)
	}
	wg.Wait()

	sort.Slice(notes, func(i, j int) bool {
		dateI, _ := time.Parse("01/02/06", notes[i].Date)
		dateJ, _ := time.Parse("01/02/06", notes[j].Date)
		return dateI.After(dateJ)
	})
	for _, note := range notes {
		noteHtmlFile, err := os.Create(filepath.Join(cfg.OutputDir, note.Name + ".html"))
		if err != nil {
			log.Printf("Error creating HTML file for %s: %v", note.Name, err)
			continue
		}
		if err := templates.ExecuteTemplate(noteHtmlFile, "NotePage", note); err != nil {
			log.Printf("Error executing template for %s: %v", note.Name, err)
		}
		noteHtmlFile.Close()
	}

	indexHtmlFile, err := os.Create(filepath.Join(cfg.OutputDir, "index.html"))
	if err != nil {
		log.Fatal(err)
	}
	defer indexHtmlFile.Close()
	err = templates.ExecuteTemplate(indexHtmlFile, "Index", notes)
	if err != nil {
		log.Fatal(err)
	}

	aboutHtmlFile, err := os.Create(filepath.Join(cfg.OutputDir, "about.html"))
	if err != nil {
		log.Fatal(err)
	}
	defer aboutHtmlFile.Close()
	err = templates.ExecuteTemplate(aboutHtmlFile, "About", notes)
	if err != nil {
		log.Fatal(err)
	}

	if err := cssCollector.WriteToFile(cfg.OutputDir); err != nil {
		log.Fatalf("Error writing highlight.css: %v", err)
	}
}
