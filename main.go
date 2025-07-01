package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const (
	Green = "\033[32m"
	Red   = "\033[31m"
	Reset = "\033[0m"
)

type PageInfo struct {
	URL         string
	Title       string
	Description string
	Links       []string
	Duration    time.Duration
	Err         error
}

var (
	visited   = make(map[string]bool)
	mu        sync.Mutex
	wg        sync.WaitGroup
	maxDepth  = flag.Int("depth", 2, "Depth of recursive crawling")
	maxPages  = flag.Int("max-pages", 100, "Maximum number of pages to crawl")
	pageCount = 0
)

func main() {
	flag.Parse()
	urls := flag.Args()

	if len(urls) == 0 {
		fmt.Println("Please provide at least one URL:")
		fmt.Println("Usage: go run main.go --depth=2 --max-pages=100 https://example.com")
		return
	}

	start := time.Now()

	outputFile, err := os.Create("results.txt")
	if err != nil {
		fmt.Println("❌ Could not create results.txt:", err)
		return
	}
	defer outputFile.Close()

	results := make(chan PageInfo)

	for _, u := range urls {
		wg.Add(1)
		go crawlWorker(u, 0, getDomain(u), results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for info := range results {
		printAndSave(info, outputFile)
	}

	summary := fmt.Sprintf("\nAll done in %s\n", time.Since(start))
	fmt.Print(summary)
	outputFile.WriteString(summary)
}

func getDomain(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return parsed.Host
}

func crawlWorker(link string, depth int, domain string, results chan<- PageInfo) {
	defer wg.Done()

	if depth > *maxDepth {
		return
	}

	mu.Lock()
	if visited[link] || pageCount >= *maxPages {
		mu.Unlock()
		return
	}
	visited[link] = true
	pageCount++
	mu.Unlock()

	start := time.Now()
	title, desc, links, err := fetchPage(link)
	duration := time.Since(start)

	info := PageInfo{
		URL:         link,
		Title:       title,
		Description: desc,
		Links:       links,
		Duration:    duration,
		Err:         err,
	}

	results <- info

	if err == nil {
		for _, href := range links {
			if strings.HasPrefix(href, "/") {
				href = "https://" + domain + href
			}
			if u, err := url.Parse(href); err == nil && u.Host == domain {
				wg.Add(1)
				go crawlWorker(href, depth+1, domain, results)
			}
		}
	}
}

func fetchPage(link string) (string, string, []string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return "", "", nil, err
	}
	defer resp.Body.Close()

	var title, desc string
	links := []string{}
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// end of document
			goto Dedup
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			switch t.Data {
			case "title":
				z.Next()
				title = z.Token().Data
			case "meta":
				attrMap := getAttrMap(t.Attr)
				if strings.ToLower(attrMap["name"]) == "description" {
					desc = attrMap["content"]
				}
			case "a":
				for _, a := range t.Attr {
					if a.Key == "href" {
						href := strings.TrimSpace(a.Val)
						if href != "" {
							links = append(links, href)
						}
					}
				}
			}
		}
	}

Dedup:
	unique := make(map[string]struct{})
	filtered := []string{}
	for _, l := range links {
		if _, seen := unique[l]; !seen {
			unique[l] = struct{}{}
			filtered = append(filtered, l)
		}
	}

	return title, desc, filtered, nil
}

func getAttrMap(attrs []html.Attribute) map[string]string {
	m := make(map[string]string)
	for _, a := range attrs {
		m[a.Key] = a.Val
	}
	return m
}

func printAndSave(info PageInfo, output *os.File) {
	if info.Err != nil {
		msg := fmt.Sprintf("%s✗ %s — ERROR: %v (⏱ %v)%s", Red, info.URL, info.Err, info.Duration, Reset)
		fmt.Println(msg)
		output.WriteString(msg + "\n")
		return
	}

	msg := fmt.Sprintf("%s✓ %s —\nTitle: \"%s\"\nDescription: \"%s\"\nLinks found: %d\n(⏱ %v)%s",
		Green, info.URL, info.Title, info.Description, len(info.Links), info.Duration, Reset)
	fmt.Println(msg)
	output.WriteString(msg + "\n")
}
