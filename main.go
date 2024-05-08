package main

import (
  "os"
  "fmt"
  "log"
  "flag"
  "time"
  "sort"
  "math"
  "sync"
  "bufio"
  "bytes"
  "strings"
  "path/filepath"
  "html/template"

  "github.com/yuin/goldmark"
  "github.com/yuin/goldmark/parser"
  "github.com/yuin/goldmark-meta"

  highlighting "github.com/yuin/goldmark-highlighting/v2"
  chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
)

type Note struct {
  Name string
  Content template.HTML
  Title string
  Date string
  ReadingTime int
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

func RemoveDuplicateLinesInPlace(input *bytes.Buffer) error {
  scanner := bufio.NewScanner(input)
  seen := make(map[string]bool)
  var output bytes.Buffer
  for scanner.Scan() {
    line := scanner.Text()
    if _, exists := seen[line]; !exists {
      seen[line] = true
      output.WriteString(line + "\n")
    }
  }
  if err := scanner.Err(); err != nil {
    return err
  }
  input.Reset()
  input.Write(output.Bytes())
  return nil
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

func main() {
  var outputDir string
  var templatesFile string
  flag.StringVar(&outputDir, "output", "./public", "Output directory for the HTML files")
  flag.StringVar(&templatesFile, "templates", "./templates.html", "A file that contains all of the templates")
  flag.Parse()
  args := flag.Args()
  if len(args) < 1 {
    log.Fatal("Error: Please provide a directory of notes")
  }
  err := EnsureDir(outputDir)
  if err != nil {
    fmt.Println("Error creating directory:", err)
  }
  notesDir := args[0]
  files, err := os.ReadDir(notesDir)
  if err != nil {
    log.Fatal(err)
  }
  templates, err := template.ParseFiles(templatesFile)
  if err != nil {
    log.Fatal("Error parsing templates: ", err)
  }
  var notes []Note
  var css bytes.Buffer
  var wg sync.WaitGroup
  var mutex sync.Mutex
  var notesMutex sync.Mutex
  markdown := goldmark.New(
    goldmark.WithExtensions(
      meta.Meta,
      highlighting.NewHighlighting(
        highlighting.WithStyle("onedark"),
        highlighting.WithCSSWriter(&css),
        highlighting.WithFormatOptions(
          chromahtml.WithClasses(true),
          chromahtml.WithLineNumbers(false),
        ),
      ),
    ),
  )
  for _, file := range files {
    wg.Add(1)
    go func(file os.DirEntry) {
      defer wg.Done()
      var buf bytes.Buffer
      noteMd, err := os.ReadFile(filepath.Join(notesDir, file.Name()))
      if err != nil {
        log.Fatal(err)
      }
      context := parser.NewContext()
      if err := markdown.Convert([]byte(noteMd), &buf, parser.WithContext(context)); err != nil {
        panic(err)
      }
      metaData := meta.Get(context)
      noteDate, err := FormatDate(metaData["date"].(string))
      if err != nil {
        panic(err)
      }
      mutex.Lock()
      err = RemoveDuplicateLinesInPlace(&css)
      if err != nil {
        log.Fatal("Error removing duplicates:", err)
      }
      mutex.Unlock()
      note := Note{
        Name: strings.TrimSuffix(file.Name(), ".md"),
        Content: template.HTML(buf.String()),
        Title: metaData["title"].(string),
        Date: noteDate,
        ReadingTime: CalculateReadingTime(&buf),
      }
      notesMutex.Lock()
      notes = append(notes, note)
      notesMutex.Unlock()
      noteHtmlFile, err := os.Create(filepath.Join(outputDir, note.Name + ".html"))
      if err != nil {
        log.Fatal(err)
      }
      defer noteHtmlFile.Close()
      err = templates.ExecuteTemplate(noteHtmlFile, "NotePage", note)
      if err != nil {
        log.Fatal("Error executing template: ", err)
      }
    }(file)
  }
  wg.Wait()
  sort.Slice(notes, func(i, j int) bool {
    dateI, _ := time.Parse("01/02/06", notes[i].Date)
    dateJ, _ := time.Parse("01/02/06", notes[j].Date)
    return dateI.After(dateJ)
  })
  indexHtmlFile, err := os.Create(filepath.Join(outputDir, "index.html"))
  if err != nil {
    log.Fatal(err)
  }
  defer indexHtmlFile.Close()
  err = templates.ExecuteTemplate(indexHtmlFile, "Index", notes)
  if err != nil {
    log.Fatal(err)
  }
  aboutHtmlFile, err := os.Create(filepath.Join(outputDir, "about.html"))
  if err != nil {
    log.Fatal(err)
  }
  defer aboutHtmlFile.Close()
  err = templates.ExecuteTemplate(aboutHtmlFile, "About", notes)
  if err != nil {
    log.Fatal(err)
  }
  err = os.WriteFile(filepath.Join(outputDir, "highlight.css"), css.Bytes(), 0666)
  if err != nil {
    log.Fatal(err)
  }
}
