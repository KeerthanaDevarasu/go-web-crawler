package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const (
	Green = "\033[32m"
	Red   = "\033[31m"
	Reset = "\033[0m"
)


func main() {
	// Get URLs from command-line arguments
	urls := os.Args[1:]

	if len(urls) == 0 {
		fmt.Println("Please provide at least one URL:")
		fmt.Println("Usage: go run main.go https://example.com https://google.com")
		return
	}

	start := time.Now()

	outputFile, err := os.Create("results.txt")
	if err != nil {
		fmt.Println("❌ Could not create results.txt:", err)
		return
	}
	defer outputFile.Close()


	var wg sync.WaitGroup
	results := make(chan string)

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			startTime := time.Now()
			title, err := fetchTitle(u)
			duration := time.Since(startTime)

			if err != nil {
				results <- fmt.Sprintf("%s✗ %s — ERROR: %v (⏱ %s)%s", Red, u, err, duration, Reset)
				return
			}
			results <- fmt.Sprintf("%s✓ %s — Title: \"%s\" (⏱ %s)%s", Green, u, title, duration, Reset)


		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println(res)
		outputFile.WriteString(res + "\n")
	}


	elapsed := time.Since(start)
	summary := fmt.Sprintf("\nAll done in %s\n", elapsed)

	fmt.Print(summary)
	outputFile.WriteString(summary)

}


//helper function

func fetchTitle(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return "", fmt.Errorf("no title found")
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "title" {
				z.Next()
				return z.Token().Data, nil
			}
		}
	}
}
