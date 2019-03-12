package goscraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	url2 "net/url"
	"regexp"
	"strings"

	"github.com/willf/bloom"

	"golang.org/x/net/html"
)

type deepUrl struct {
	depth int
	urls  []string
}

func Run(url string, depth int, pattern string) {
	urlQueue := make(chan deepUrl, 0)
	seenUrl := bloom.New(100, 5)
	baseUrl := parseBaseURL(url)
	go scrape(baseUrl, url, depth, urlQueue)
	seenUrl.AddString(url)

	outstanding := 1

	for outstanding > 0 {

		u := <-urlQueue
		outstanding--

		if depth <= 0 {
			continue
		}
		for _, url := range u.urls {
			if isUrlPatternValid(url, pattern) && !seenUrl.TestString(url) {
				log.Println("going to scrape url :", url)
				go scrape(baseUrl, url, depth, urlQueue)
				seenUrl.AddString(url)
			}
		}
	}
}

func scrape(baseUrl string, url string, depth int, q chan deepUrl) {

	u, err := url2.Parse(url)
	if err != nil || u.Scheme == "" || u.Host == "" || u.Path == "" {
		log.Println("invalid input url :", url)
		return
	}
	res, err := http.Get(u.String())

	if err != nil {
		log.Println("Error occurred while fetching the url :", url, err)
		return
	}
	defer res.Body.Close()

	urls := extractUrl(res.Body)

	if len(urls) > 0 {
		q <- deepUrl{depth - 1, resolveRelative(baseUrl, urls)}
	}
}

func parseBaseURL(u string) string {
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
}

func resolveRelative(baseURL string, hrefs []string) []string {
	internalUrls := []string{}

	for _, href := range hrefs {
		u, _ := url2.Parse(href)

		//if url is already has baseUrl or of different domain?
		if strings.HasPrefix(href, baseURL) || u.Host != "" {
			internalUrls = append(internalUrls, href)
			continue
		}
		var resolvedURL string

		if href == "" {
			resolvedURL = baseURL
		} else if strings.HasPrefix(href, "/") {
			resolvedURL = fmt.Sprintf("%s%s", baseURL, href)

		} else {
			resolvedURL = fmt.Sprintf("%s/%s", baseURL, href)
		}
		internalUrls = append(internalUrls, resolvedURL)
	}

	return internalUrls
}
func extractUrl(body io.ReadCloser) []string {
	urls := make([]string, 1, 1)
	z := html.NewTokenizer(body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return urls
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			ok, url := getHref(t)

			if ok {

				log.Printf("found link %s", url)
				u, err := url2.Parse(url)
				if err == nil {
					urls = append(urls, u.String())
				}

			}
		}
	}

}

func isUrlPatternValid(url string, pattern string) bool {
	if pattern == "" {
		return true
	}
	r, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return url != "" && r.MatchString(url)

}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}
