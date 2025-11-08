package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/time/rate"
)

type Config struct {
	BaseURL     string
	MaxDepth    int
	RateLimit   float64 // requests per second
	UserAgent   string
	OutputDir   string
	URLListFile string
	AcceptRegex string
	RejectRegex string
	Timeout     time.Duration
	Verbose     bool
	// Special rate limiting for servers that block after X requests in Y seconds
	SpecialRate   bool          // Enable special rate limiting mode
	MaxRequests   int           // Max requests in time window (default: 2)
	TimeWindow    time.Duration // Time window (default: 5s)
	BlockDuration time.Duration // How long server blocks access (default: 10s)
}

type Spider struct {
	config      Config
	limiter     *rate.Limiter
	client      *http.Client
	visited     map[string]bool
	discovered  []string
	mu          sync.RWMutex
	acceptRegex *regexp.Regexp
	rejectRegex *regexp.Regexp
	// Special rate limiting tracking
	requestTimes  []time.Time // Track recent request times
	lastBlockTime time.Time   // When we were last blocked
	requestTimeMu sync.Mutex  // Mutex for request time tracking
}

type DiscoveredFile struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Size     int64  `json:"size,omitempty"`
	Modified string `json:"modified,omitempty"`
	Type     string `json:"type"`
}

func NewSpider(config Config) (*Spider, error) {
	s := &Spider{
		config:       config,
		limiter:      rate.NewLimiter(rate.Limit(config.RateLimit), 1),
		visited:      make(map[string]bool),
		discovered:   make([]string, 0),
		requestTimes: make([]time.Time, 0),
	}

	// Setup HTTP client with timeout
	s.client = &http.Client{
		Timeout: config.Timeout,
	}

	// Compile regex patterns
	if config.AcceptRegex != "" {
		var err error
		s.acceptRegex, err = regexp.Compile(config.AcceptRegex)
		if err != nil {
			return nil, fmt.Errorf("invalid accept regex: %w", err)
		}
	}

	if config.RejectRegex != "" {
		var err error
		s.rejectRegex, err = regexp.Compile(config.RejectRegex)
		if err != nil {
			return nil, fmt.Errorf("invalid reject regex: %w", err)
		}
	}

	return s, nil
}

func (s *Spider) makeRequest(ctx context.Context, method, url string) (*http.Response, error) {
	if s.config.SpecialRate {
		// Handle special rate limiting (e.g., 2 requests per 5 seconds)
		if err := s.handleSpecialRateLimit(ctx); err != nil {
			return nil, err
		}
	} else {
		// Regular rate limiting
		if err := s.limiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", s.config.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	if s.config.Verbose {
		log.Printf("Requesting: %s %s", method, url)
	}

	resp, err := s.client.Do(req)

	// If using special rate limiting, track this request and check for blocks
	if s.config.SpecialRate {
		s.trackRequest()

		// Check if we got blocked (common HTTP status codes for rate limiting)
		if resp != nil && (resp.StatusCode == 429 || resp.StatusCode == 503 || resp.StatusCode == 502) {
			if s.config.Verbose {
				log.Printf("Detected rate limiting block (HTTP %d), waiting %v", resp.StatusCode, s.config.BlockDuration)
			}
			s.handleBlock()
			// Wait for the block duration before continuing
			time.Sleep(s.config.BlockDuration)
		}
	}

	return resp, err
}

func (s *Spider) handleSpecialRateLimit(ctx context.Context) error {
	s.requestTimeMu.Lock()
	defer s.requestTimeMu.Unlock()

	now := time.Now()

	// If we were recently blocked, wait for the block to expire
	if !s.lastBlockTime.IsZero() {
		timeSinceBlock := now.Sub(s.lastBlockTime)
		if timeSinceBlock < s.config.BlockDuration {
			remainingWait := s.config.BlockDuration - timeSinceBlock
			if s.config.Verbose {
				log.Printf("Still in block period, waiting %v", remainingWait)
			}
			s.requestTimeMu.Unlock()
			time.Sleep(remainingWait)
			s.requestTimeMu.Lock()
			now = time.Now()
		}
	}

	// Clean old request times (outside the window)
	cutoff := now.Add(-s.config.TimeWindow)
	var validTimes []time.Time
	for _, t := range s.requestTimes {
		if t.After(cutoff) {
			validTimes = append(validTimes, t)
		}
	}
	s.requestTimes = validTimes

	// Check if we've hit the request limit
	if len(s.requestTimes) >= s.config.MaxRequests {
		// Calculate how long to wait
		oldestRequest := s.requestTimes[0]
		waitTime := s.config.TimeWindow - now.Sub(oldestRequest)

		if waitTime > 0 {
			if s.config.Verbose {
				log.Printf("Rate limit reached (%d requests in %v), waiting %v",
					len(s.requestTimes), s.config.TimeWindow, waitTime)
			}
			s.requestTimeMu.Unlock()

			// Use context-aware sleep
			select {
			case <-ctx.Done():
				s.requestTimeMu.Lock()
				return ctx.Err()
			case <-time.After(waitTime):
				s.requestTimeMu.Lock()
			}
		}
	}

	return nil
}

func (s *Spider) trackRequest() {
	s.requestTimeMu.Lock()
	defer s.requestTimeMu.Unlock()

	s.requestTimes = append(s.requestTimes, time.Now())
}

func (s *Spider) handleBlock() {
	s.requestTimeMu.Lock()
	defer s.requestTimeMu.Unlock()

	s.lastBlockTime = time.Now()
	// Clear request times since we got blocked
	s.requestTimes = nil
}

func (s *Spider) shouldAcceptURL(rawURL string) bool {
	if s.rejectRegex != nil && s.rejectRegex.MatchString(rawURL) {
		return false
	}

	if s.acceptRegex != nil && !s.acceptRegex.MatchString(rawURL) {
		return false
	}

	return true
}

func (s *Spider) isFile(rawURL string) bool {
	// Simple heuristic: if it doesn't end with / and has an extension or no extension
	return !strings.HasSuffix(rawURL, "/")
}

func (s *Spider) discover(ctx context.Context, targetURL string, depth int) error {
	if depth > s.config.MaxDepth {
		return nil
	}

	s.mu.Lock()
	if s.visited[targetURL] {
		s.mu.Unlock()
		return nil
	}
	s.visited[targetURL] = true
	s.mu.Unlock()

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return err
	}

	// Check if it's a file or directory
	if s.isFile(targetURL) {
		// For files, just do a HEAD request to check if they exist
		resp, err := s.makeRequest(ctx, "HEAD", targetURL)
		if err != nil {
			if s.config.Verbose {
				log.Printf("Error checking file %s: %v", targetURL, err)
			}
			return nil
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 && s.shouldAcceptURL(targetURL) {
			s.mu.Lock()
			s.discovered = append(s.discovered, targetURL)
			s.mu.Unlock()

			if s.config.Verbose {
				log.Printf("Found file: %s", targetURL)
			}
		}
		return nil
	}

	// For directories, get the content and parse for links
	resp, err := s.makeRequest(ctx, "GET", targetURL)
	if err != nil {
		if s.config.Verbose {
			log.Printf("Error fetching directory %s: %v", targetURL, err)
		}
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit concurrent requests

	doc.Find("a[href]").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}

		// Resolve relative URLs
		linkURL, err := parsedURL.Parse(href)
		if err != nil {
			return
		}

		// Skip external links
		if linkURL.Host != parsedURL.Host {
			return
		}

		// Skip parent directory links
		if strings.Contains(href, "..") {
			return
		}

		fullURL := linkURL.String()

		// Skip if already visited
		s.mu.RLock()
		alreadyVisited := s.visited[fullURL]
		s.mu.RUnlock()

		if alreadyVisited {
			return
		}

		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := s.discover(ctx, url, depth+1); err != nil && s.config.Verbose {
				log.Printf("Error discovering %s: %v", url, err)
			}
		}(fullURL)
	})

	wg.Wait()
	return nil
}

func (s *Spider) DiscoverStructure(ctx context.Context) ([]string, error) {
	if s.config.Verbose {
		log.Printf("Starting discovery from: %s", s.config.BaseURL)
	}

	err := s.discover(ctx, s.config.BaseURL, 0)
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	result := make([]string, len(s.discovered))
	copy(result, s.discovered)
	s.mu.RUnlock()

	sort.Strings(result)
	return result, nil
}

func (s *Spider) SaveDiscoveredURLs(urls []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, url := range urls {
		_, err := writer.WriteString(url + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Spider) DownloadFiles(ctx context.Context, urls []string) error {
	for _, rawURL := range urls {
		if err := s.downloadFile(ctx, rawURL); err != nil {
			log.Printf("Error downloading %s: %v", rawURL, err)
		}
	}
	return nil
}

func (s *Spider) downloadFile(ctx context.Context, rawURL string) error {
	resp, err := s.makeRequest(ctx, "GET", rawURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d for %s", resp.StatusCode, rawURL)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	// Create local file path
	localPath := filepath.Join(s.config.OutputDir, parsedURL.Host, parsedURL.Path)

	// Ensure directory exists
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create file
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	if s.config.Verbose {
		log.Printf("Downloaded: %s -> %s", rawURL, localPath)
	}

	return nil
}

func loadURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}

	return urls, scanner.Err()
}

func main() {
	var (
		baseURL      = flag.String("url", "", "Base URL to spider")
		maxDepth     = flag.Int("depth", 3, "Maximum crawling depth")
		rateLimit    = flag.Float64("rate", 1.0, "Rate limit (requests per second)")
		userAgent    = flag.String("user-agent", "WebSpider/1.0 (Compatible crawler)", "User agent string")
		outputDir    = flag.String("output", "./downloads", "Output directory for downloads")
		urlListFile  = flag.String("urls", "", "File containing URLs to download (skip discovery)")
		acceptRegex  = flag.String("accept", "", "Regex pattern for URLs to accept")
		rejectRegex  = flag.String("reject", "", "Regex pattern for URLs to reject")
		timeout      = flag.Duration("timeout", 30*time.Second, "HTTP request timeout")
		verbose      = flag.Bool("verbose", false, "Verbose output")
		saveList     = flag.String("save-list", "discovered-urls.txt", "File to save discovered URLs")
		discoverOnly = flag.Bool("discover-only", false, "Only discover URLs, don't download")
		// Special rate limiting flags
		specialRate   = flag.Bool("special-rate", false, "Enable special rate limiting mode")
		maxRequests   = flag.Int("max-requests", 2, "Max requests in time window (for special rate limiting)")
		timeWindow    = flag.Duration("time-window", 5*time.Second, "Time window for request limiting")
		blockDuration = flag.Duration("block-duration", 10*time.Second, "Duration server blocks access after rate limit")
	)

	flag.Parse()

	if *baseURL == "" && *urlListFile == "" {
		fmt.Println("Usage: webspider -url <base-url> [options]")
		fmt.Println("   or: webspider -urls <url-list-file> [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  # Discover directory structure:")
		fmt.Println("  webspider -url https://example.com/files/ -discover-only -depth 5")
		fmt.Println("  ")
		fmt.Println("  # Download specific files from list:")
		fmt.Println("  webspider -urls discovered-urls.txt -rate 0.5")
		fmt.Println("  ")
		fmt.Println("  # Discover with special rate limiting (2 req/5sec, 10sec block):")
		fmt.Println("  webspider -url https://sensitive.com/ -special-rate -discover-only")
		fmt.Println("  ")
		fmt.Println("  # Custom special rate limiting:")
		fmt.Println("  webspider -url https://example.com/ -special-rate -max-requests 3 -time-window 10s -block-duration 15s")
		fmt.Println("  ")
		fmt.Println("  # Discover and filter by file type:")
		fmt.Println("  webspider -url https://example.com/ -accept '\\.(pdf|zip|tar\\.gz)$' -rate 2")
		os.Exit(1)
	}

	config := Config{
		BaseURL:       *baseURL,
		MaxDepth:      *maxDepth,
		RateLimit:     *rateLimit,
		UserAgent:     *userAgent,
		OutputDir:     *outputDir,
		URLListFile:   *urlListFile,
		AcceptRegex:   *acceptRegex,
		RejectRegex:   *rejectRegex,
		Timeout:       *timeout,
		Verbose:       *verbose,
		SpecialRate:   *specialRate,
		MaxRequests:   *maxRequests,
		TimeWindow:    *timeWindow,
		BlockDuration: *blockDuration,
	}

	spider, err := NewSpider(config)
	if err != nil {
		log.Fatalf("Error creating spider: %v", err)
	}

	ctx := context.Background()

	if *urlListFile != "" {
		// Download mode: load URLs from file
		urls, err := loadURLsFromFile(*urlListFile)
		if err != nil {
			log.Fatalf("Error loading URL list: %v", err)
		}

		fmt.Printf("Loaded %d URLs from %s\n", len(urls), *urlListFile)

		if !*discoverOnly {
			fmt.Println("Starting downloads...")
			if err := spider.DownloadFiles(ctx, urls); err != nil {
				log.Fatalf("Error downloading files: %v", err)
			}
			fmt.Println("Downloads completed!")
		}
	} else {
		// Discovery mode
		fmt.Printf("Starting discovery from: %s\n", *baseURL)
		urls, err := spider.DiscoverStructure(ctx)
		if err != nil {
			log.Fatalf("Error during discovery: %v", err)
		}

		fmt.Printf("Discovered %d URLs\n", len(urls))

		// Save discovered URLs
		if err := spider.SaveDiscoveredURLs(urls, *saveList); err != nil {
			log.Fatalf("Error saving URL list: %v", err)
		}
		fmt.Printf("URL list saved to: %s\n", *saveList)

		if !*discoverOnly {
			fmt.Println("Starting downloads...")
			if err := spider.DownloadFiles(ctx, urls); err != nil {
				log.Fatalf("Error downloading files: %v", err)
			}
			fmt.Println("Downloads completed!")
		} else {
			fmt.Println("Discovery completed. Edit the URL list and run with -urls flag to download.")
		}
	}
}
