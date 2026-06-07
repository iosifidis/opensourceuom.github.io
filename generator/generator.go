package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// Generator handles parsing input files and generating the static site
type Generator struct {
	InputDir  string
	OutputDir string
	Templates *template.Template
	Posts     []*Post
	Pages     []*Post
}

// NewGenerator creates a new Generator instance
func NewGenerator(inputDir, outputDir, templatesDir string) (*Generator, error) {
	// Parse all HTML templates
	tmplPattern := filepath.Join(templatesDir, "*.html")
	tmpl, err := template.ParseGlob(tmplPattern)
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %v", err)
	}

	return &Generator{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Templates: tmpl,
		Posts:     []*Post{},
		Pages:     []*Post{},
	}, nil
}

// Scan reads the content directory and parses all articles and pages
func (g *Generator) Scan() error {
	// Parse blog posts
	blogDir := filepath.Join(g.InputDir, "blog")
	blogEntries, err := os.ReadDir(blogDir)
	if err == nil {
		for _, entry := range blogEntries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				slug := strings.TrimSuffix(entry.Name(), ".md")
				mdPath := filepath.Join(blogDir, entry.Name())
				
				post, err := ParseMarkdown(mdPath, slug, false)
				if err != nil {
					fmt.Printf("Warning: failed to parse %s: %v\n", slug, err)
					continue
				}
				g.Posts = append(g.Posts, post)
				fmt.Printf("Parsed %s\n", slug)
			}
		}
	}

	// Parse pages
	pagesDir := filepath.Join(g.InputDir, "pages")
	pageEntries, err := os.ReadDir(pagesDir)
	if err == nil {
		for _, entry := range pageEntries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				slug := strings.TrimSuffix(entry.Name(), ".md")
				mdPath := filepath.Join(pagesDir, entry.Name())
				
				post, err := ParseMarkdown(mdPath, slug, true)
				if err != nil {
					fmt.Printf("Warning: failed to parse %s: %v\n", slug, err)
					continue
				}
				g.Pages = append(g.Pages, post)
				fmt.Printf("Parsed %s\n", slug)
			}
		}
	}

	// Sort posts by date descending
	sort.Slice(g.Posts, func(i, j int) bool {
		return g.Posts[i].Date.After(g.Posts[j].Date)
	})

	return nil
}

// Generate builds the static site
func (g *Generator) Generate() error {
	// 1. Create output directory structure
	if err := os.MkdirAll(g.OutputDir, 0755); err != nil {
		return err
	}

	// 2. Generate Posts (in /blog/slug/index.html)
	for _, post := range g.Posts {
		postDir := filepath.Join(g.OutputDir, "blog", post.Slug)
		if err := os.MkdirAll(postDir, 0755); err != nil {
			return err
		}

		outFile := filepath.Join(postDir, "index.html")
		if err := g.renderTemplate(outFile, "post.html", map[string]interface{}{
			"Post": post,
			"Site": map[string]string{"Title": "Open Source UoM"},
		}); err != nil {
			return err
		}
	}

	// 3. Generate Pages (in /slug/index.html)
	for _, page := range g.Pages {
		// special mapping for blog index
		if page.Slug == "blog" {
			pageDir := filepath.Join(g.OutputDir, "blog")
			if err := os.MkdirAll(pageDir, 0755); err != nil {
				return err
			}
			outFile := filepath.Join(pageDir, "index.html")
			if err := g.renderTemplate(outFile, "blog-index.html", map[string]interface{}{
				"Page":  page,
				"Posts": g.Posts, // Pass posts for the listing
				"Site":  map[string]string{"Title": "Open Source UoM"},
			}); err != nil {
				return err
			}
		} else {
			pageDir := filepath.Join(g.OutputDir, page.Slug)
			if err := os.MkdirAll(pageDir, 0755); err != nil {
				return err
			}
			outFile := filepath.Join(pageDir, "index.html")
			if err := g.renderTemplate(outFile, "page.html", map[string]interface{}{
				"Page": page,
				"Site": map[string]string{"Title": "Open Source UoM"},
			}); err != nil {
				return err
			}
		}
	}

	// 4. Generate Main Index (Home page)
	homePost := &Post{
		Title: "Home",
		IsPage: true,
	}

	// Only show 4 recent posts on home page
	recentPosts := g.Posts
	if len(recentPosts) > 4 {
		recentPosts = recentPosts[:4]
	}

	if err := g.renderTemplate(filepath.Join(g.OutputDir, "index.html"), "index.html", map[string]interface{}{
		"Page":        homePost,
		"RecentPosts": recentPosts,
		"Site":        map[string]string{"Title": "Open Source UoM"},
	}); err != nil {
		return err
	}

	return nil
}

func (g *Generator) renderTemplate(filePath string, tmplName string, data interface{}) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer f.Close()

	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", tmplName),
	)
	if err != nil {
		return fmt.Errorf("error parsing template %s for %s: %v", tmplName, filePath, err)
	}

	if err := tmpl.ExecuteTemplate(f, "base", data); err != nil {
		return fmt.Errorf("error executing template %s for %s: %v", tmplName, filePath, err)
	}
	return nil
}
