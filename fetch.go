package main

import (
	"log/slog"
	"flag"
	"net/url"
	"os"

	"autify/fetch/internal/config"
	f "autify/fetch/internal/fetcher"
)

type UrlFetcher interface {
	Fetch(string) error
}

// Setup command line flags
var meta *bool
func init() {
	meta = flag.Bool("metadata", false, "Display metadata, default to false")
}

// worker retrieves urls from the input channel and calls the fetcher method
// Then sends the result (error if any) thru the output channel
// UrlFetcher interface allows easily mocking the fetcher, and makes unit testing easier
func worker(id int, urlChan <-chan string, errChan chan<- error, fetcher UrlFetcher) {
	for url := range urlChan {
		errChan <- fetcher.Fetch(url)
    }
}

// sanitizeArgs sanitizes command line arguments and keeps only the valid URLs
func sanitizeArgs(urls []string) []string {
	i := 0
	for _, s := range urls {
		_, err := url.ParseRequestURI(s)
		if err != nil {
			slog.Warn("Invalid argument", "URL", s, "Error", err.Error())
			continue
		}
		// Move the valid URL to the front of the slice
		urls[i] = s
		i++
	}

	if i == 0 {
		slog.Error("No valid URLs to download")
		os.Exit(2)
	}

	return urls[:i]
}

func main() {
	var err error

	conf, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to read config file", "Error", err.Error())
		os.Exit(1)
	}

	flag.Parse()
	urls := sanitizeArgs(flag.Args())
	
    urlChan := make(chan string, len(urls))
    errChan := make(chan error, len(urls))

	fetcher := f.InitFetcher(conf.OutputDirectory, meta)

	// Initialize a concurrent worker for each URL
	// send URLs to fetch thru a channel
    for w := 0; w < conf.Workers; w++ {
        go worker(w, urlChan, errChan, fetcher)
    }

    for i := 0; i < len(urls); i++ {
        urlChan <- urls[i]
    }
    close(urlChan)

	// Collect results thru channel from each worker, and exit main when all URLs have been finished processing
    for a := 1; a <= len(urls); a++ {
        err = <-errChan

		if err != nil {
			slog.Error(err.Error())
		}
    }
	close(errChan)
}