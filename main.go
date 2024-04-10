package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const maxDepth = 6

var visited = make(map[string]bool)

func main() {
	startURL := "https://en.wikipedia.org/wiki/Israel"
	fmt.Println("Searching for the path to Hitler page...")
	path, err := findHitlerPath(startURL, 1)
	if err != nil {
		log.Fatal(err)
	}
	if path != "" {
		fmt.Println("Path to Hitler page found:", path)
	} else {
		fmt.Println("Hitler not found in 6 hops.")
	}
}

func findHitlerPath(url string, depth int) (string, error) {
	if depth > maxDepth {
		return "", nil
	}

	doc, err := fetchURL(url)
	if err != nil {
		return "", err
	}

	links := extractLinks(doc)

	for _, link := range links {
		if strings.EqualFold(link, "https://en.wikipedia.org/wiki/Adolf_Hitler") {
			fmt.Println("Hitler found", url)
			return url, nil
		}
		if !visited[link] {
			visited[link] = true
			_, err := findHitlerPath(link, depth+1)
			if err != nil {
				return "", err
			}
		}
	}
	return "", nil
}

func fetchURL(url string) (*goquery.Document, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.9999.999 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func extractLinks(doc *goquery.Document) []string {
	var links []string
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && strings.HasPrefix(link, "/wiki/") {
			links = append(links, "https://en.wikipedia.org"+link)
		}
	})
	return links
}
