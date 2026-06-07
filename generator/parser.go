package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// Post represents a parsed article or page
type Post struct {
	Slug          string
	Title         string
	Content       string
	Date          time.Time
	Author        string
	Tags          []string
	FeaturedImage string
	IsPage        bool
	OriginalPath  string
}

// FrontMatter represents the YAML frontmatter of a markdown file
type FrontMatter struct {
	Title         string   `yaml:"title"`
	Date          string   `yaml:"date,omitempty"`
	Author        string   `yaml:"author,omitempty"`
	Tags          []string `yaml:"tags,omitempty"`
	FeaturedImage string   `yaml:"featured_image,omitempty"`
}

// ParseMarkdown reads a Markdown file and extracts Post data
func ParseMarkdown(filePath string, slug string, isPage bool) (*Post, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	content := string(data)
	var fm FrontMatter
	var body string

	// Extract frontmatter
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			err := yaml.Unmarshal([]byte(parts[1]), &fm)
			if err != nil {
				return nil, fmt.Errorf("error parsing frontmatter in %s: %v", filePath, err)
			}
			body = parts[2]
		} else {
			body = content
		}
	} else {
		body = content
	}

	// Parse date
	var date time.Time
	if fm.Date != "" {
		date, _ = time.Parse(time.RFC3339, fm.Date)
	}

	// Convert Markdown body to HTML
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(body), &buf); err != nil {
		return nil, fmt.Errorf("error converting markdown in %s: %v", filePath, err)
	}

	post := &Post{
		Slug:          slug,
		OriginalPath:  filePath,
		Title:         fm.Title,
		Author:        fm.Author,
		Tags:          fm.Tags,
		FeaturedImage: fm.FeaturedImage,
		Date:          date,
		Content:       buf.String(),
		IsPage:        isPage,
	}

	return post, nil
}
