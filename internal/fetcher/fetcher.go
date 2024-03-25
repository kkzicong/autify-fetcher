package fetcher

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

type UrlFetcher struct {
	outputDir string
	meta      *bool
}

// InitFetcher initializes a fetcher
func InitFetcher(outputDir string, meta *bool) *UrlFetcher {
	return &UrlFetcher{outputDir, meta}
}

// Fetch downloads the HTML content of the given url into a file in the directory specified in config.yml
func (f *UrlFetcher) Fetch(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	filename := fmt.Sprintf("%s/%s.html", f.outputDir, path.Base(url))

	// Open file in truncate mode to overwrite if it exists
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	htmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body %s: %w", url, err)
	}
	file.Write(htmlData)

	if *f.meta {
		displayMetaData(url, htmlData)
	}

	return nil
}

// displayMetaData counts the links and images of a HTML body, and print out the info
func displayMetaData(url string, data []byte) {
	linkCount := bytes.Count(data, []byte("</a>"))
	imgCount := bytes.Count(data, []byte("<img "))
	t := time.Now()

	fmt.Printf("site: %s\n", url)
	fmt.Printf("num_links: %d\n", linkCount)
	fmt.Printf("images: %d\n", imgCount)
	fmt.Printf("last_fetch: %s\n\n", t.String())
}
