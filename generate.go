package nebel

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"text/template"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/yosssi/gohtml"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type Header struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
}

type Post struct {
	Title         string
	Date          time.Time
	RawContent    string
	Path          string
	ParsedContent string
	NextPost      *Post
	PrevPost      *Post
	FullContent   string
	Index         bool
}

func Generate() error {
	posts, err := createPostObjects()
	if err != nil {
		return err
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.Before(posts[j].Date)
	})

	if err := processPosts(posts); err != nil {
		return err
	}

	if err := writePostFiles(posts); err != nil {
		return err
	}

	if err := generateIndexHTML(posts); err != nil {
		return err
	}

	if err := generateAtomXML(posts); err != nil {
		return err
	}

	return copyStaticFiles()
}

func processPosts(posts []*Post) error {
	count := 1
	for pos, post := range posts {
		if err := post.convertMarkdown(); err != nil {
			return err
		}

		currentDate := post.Date.Format("2006/01/02")
		if pos > 0 && posts[pos-1].Date.Format("2006/01/02") == currentDate {
			count++
		} else if pos > 0 {
			count = 1
		}

		post.Path = fmt.Sprintf("/blog/%s/%d", currentDate, count)

		if pos > 0 {
			post.PrevPost = posts[pos-1]
		}
		if pos < len(posts)-1 {
			post.NextPost = posts[pos+1]
		}
	}
	return nil
}

func writePostFiles(posts []*Post) error {
	for pos, post := range posts {
		if err := post.processLayout(); err != nil {
			return err
		}

		if pos > len(posts)-3 {
			path := filepath.Join("public", post.Path, "index.html")
			if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return err
			}
			if err := os.WriteFile(path, []byte(formatHTML(post.FullContent)), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateIndexHTML(posts []*Post) error {
	latestPost := posts[len(posts)-1]
	indexHTML, err := latestPost.processPostTemplate(true)
	if err != nil {
		return err
	}

	path := filepath.Join("public", "index.html")
	return os.WriteFile(path, []byte(formatHTML(*indexHTML)), 0644)
}

func copyStaticFiles() error {
	return filepath.Walk("static", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel("static", path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join("public", relativePath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, os.ModePerm)
		}

		input, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, input, info.Mode())
	})
}

func createPostObjects() ([]*Post, error) {
	var posts []*Post

	files, err := os.ReadDir("posts")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		post, err := createPostObject(file)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func generateAtomXML(posts []*Post) error {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	posts = posts[0:9]

	tmpl, err := template.New("atom.xml").Funcs(template.FuncMap{
		"formatDate": formatDate,
	}).ParseFiles("layouts/atom.xml")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, posts)
	if err != nil {
		return err
	}

	path := filepath.Join("public", "atom.xml")
	return os.WriteFile(path, buf.Bytes(), 0644)
}

func createPostObject(file os.DirEntry) (*Post, error) {
	post := &Post{}

	path := filepath.Join("posts", file.Name())

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	inHeader := false
	inBody := false

	headerString := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if !inHeader && !inBody && line == "---" {
			inHeader = true
			continue
		}

		if inHeader && !inBody && line == "---" {
			inHeader = false
			inBody = true
			continue
		}

		if inHeader {
			headerString += line + "\n"
		}

		if inBody {
			post.RawContent += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	header := &Header{}
	if err := yaml.Unmarshal([]byte(headerString), header); err != nil {
		return nil, err
	}

	post.Title = header.Title
	post.Date, _ = time.ParseInLocation("2006-01-02 15:04:05 +0900", header.Date, time.FixedZone("Asia/Tokyo", 9*60*60))

	if post.Date.IsZero() {
		post.Date, _ = time.ParseInLocation("2006-01-02 15:04:05", header.Date, time.FixedZone("Asia/Tokyo", 9*60*60))
	}

	if post.Date.IsZero() {
		post.Date, _ = time.ParseInLocation("2006-01-02 15:04", header.Date, time.FixedZone("Asia/Tokyo", 9*60*60))
	}

	return post, nil
}

func (p *Post) convertMarkdown() error {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		))

	var buf bytes.Buffer
	err := md.Convert([]byte(p.RawContent), &buf)
	if err != nil {
		return err
	}

	p.ParsedContent = buf.String()

	return nil
}

func (p *Post) processLayout() error {
	content, err := p.processPostTemplate(false)
	if err != nil {
		return err
	}

	p.FullContent = *content

	return nil
}

func (p *Post) processPostTemplate(index bool) (*string, error) {
	p.Index = index

	tmpl, err := template.New("post.html").Funcs(template.FuncMap{
		"formatDate": formatDate,
	}).ParseFiles("layouts/post.html")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, p)
	if err != nil {
		return nil, err
	}

	content := buf.String()

	return &content, nil
}

func formatDate(t time.Time, layout string) string {
	return t.Format(layout)
}

// formatHTML formats HTML and cleans up whitespace inside inline <code> tags
// that gohtml.Format() incorrectly adds
func formatHTML(html string) string {
	formatted := gohtml.Format(html)

	// Remove whitespace inside inline <code> tags
	// Pattern matches <code> followed by whitespace, content, whitespace, </code>
	// but excludes <pre><code> blocks (which are already handled correctly by gohtml)
	re := regexp.MustCompile(`(?s)<code>\s*(.*?)\s*</code>`)
	return re.ReplaceAllString(formatted, "<code>$1</code>")
}
