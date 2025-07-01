
```markdown
# Go Concurrent Web Crawler

A blazing-fast, concurrent CLI tool built in Go that crawls multiple websites in parallel, extracts page titles, measures response times, and logs results to both the terminal and a file — with color-coded output for quick readability.

---

## Features

- Concurrent crawling using goroutines and channels
- Extracts `<title>` tag from each webpage
- Measures per-URL and total execution time
- Outputs results to both terminal and `results.txt`
- Color-coded terminal output (green = success, red = error)
- Lightweight and fast — only one external package

---

## Technologies Used

- [Go](https://golang.org/) — Language & concurrency
- `net/http` — Fetching pages
- `golang.org/x/net/html` — HTML parsing
- ANSI escape codes — Terminal colors
- Goroutines + `sync.WaitGroup` — Concurrency control
- Channels — Output coordination

---

## Installation

1. Clone this repository:
```bash
git clone https://github.com/KeerthanaDevarasu/go-web-crawler.git
cd go-web-crawler
```

2. Install dependencies:

```bash
go mod tidy
```

---

## Usage

To run the crawler, use the following command with one or more URLs:

```bash
go run main.go https://example.com https://github.com https://golang.org
```

You can provide any number of URLs as arguments.

---

## Example Input/Output

**Input:**

```bash
go run main.go https://github.com https://notarealwebsite.bogus
```

**Output:**

```
✓ https://github.com — Title: "GitHub · Build and ship software on a single, collaborative platform · GitHub" (⏱ 481.48ms)
✗ https://notarealwebsite.bogus — ERROR: Get "https://notarealwebsite.bogus": dial tcp: lookup notarealwebsite.bogus: no such host (⏱ 25.43ms)

All done in 507.12ms
```

* Successful URLs are shown in green
* Failed URLs (timeouts, 404s, etc.) are shown in red

---

## Output File: `results.txt`

In addition to terminal output, a file named `results.txt` is automatically created in your project directory after each run.

It contains the same crawl results in plain text, for example:

```
✓ https://github.com — Title: "GitHub · Build and ship software..." (⏱ 481.48ms)
✗ https://notarealwebsite.bogus — ERROR: no such host (⏱ 25.43ms)

All done in 507.12ms
```

---

## Author

[Keerthana Devarasu](https://github.com/KeerthanaDevarasu)

```
