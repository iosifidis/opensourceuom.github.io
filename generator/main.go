package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Paths
	baseDir := filepath.Join("..") // Assuming we run from generator/
	inputDir := filepath.Join(baseDir, "content")
	outputDir := filepath.Join(baseDir, "public")
	templatesDir := filepath.Join("templates")

	fmt.Println("Starting Open Source UoM Static Site Generator...")

	gen, err := NewGenerator(inputDir, outputDir, templatesDir)
	if err != nil {
		fmt.Printf("Error initializing generator: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Scanning input directory...")
	if err := gen.Scan(); err != nil {
		fmt.Printf("Error scanning files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d posts and %d pages.\n", len(gen.Posts), len(gen.Pages))

	fmt.Println("Generating static HTML...")
	if err := gen.Generate(); err != nil {
		fmt.Printf("Error generating site: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Copying static assets (storage, css, etc)...")
	// Copy storage
	if err := copyDir(filepath.Join(baseDir, "storage"), filepath.Join(outputDir, "storage")); err != nil {
		fmt.Printf("Warning: failed to copy storage dir: %v\n", err)
	}

	// Copy custom static assets (like CSS)
	if err := copyDir(filepath.Join("static"), filepath.Join(outputDir, "static")); err != nil {
		fmt.Printf("Warning: failed to copy static dir: %v\n", err)
	}

	fmt.Println("Site generation complete! Output is in public/")
}

// Helper to copy directories recursively
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// Helper to copy files
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
