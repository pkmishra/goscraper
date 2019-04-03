package goscraper

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	url2 "net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/willf/bloom"

	"golang.org/x/net/html"
)

const (
	MAX_RATE = 3
)

type deepUrl struct {
	depth int
	urls  []string
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type HttpFetcher struct {
}

var mux sync.Mutex
var outstanding int

func (f HttpFetcher) Fetch(url string) (body string, urls []string, err error) {
	u, err := url2.Parse(url)
	if err != nil || u.Scheme == "" || u.Host == "" || u.Path == "" {
		log.Println("invalid input url :", url, err)
		return
	}
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Get(u.String())
	if err != nil {
		log.Println("Error occurred while fetching the url :", url, err)
		return
	}
	defer res.Body.Close()

	urls = resolveRelative(parseBaseURL(url), extractUrl(res.Body))
	b, err := ioutil.ReadAll(res.Body)
	body = string(b)
	return body, urls, err
}

var rateSemaphore chan struct{}

func Run(fetcher Fetcher, url string, depth int, pattern string, rate int) {
	if rate > MAX_RATE {
		rate = MAX_RATE
	}
	//semaphore channel
	rateSemaphore = make(chan struct{}, rate)
	for i := 0; i < rate; i++ {
		rateSemaphore <- struct{}{}
	}
	var wg sync.WaitGroup
	urlQ := make(chan deepUrl)
	seenUrl := bloom.New(100000, 5)
	wg.Add(1)
	go scrape(fetcher, url, depth, &urlQ, &wg)

	seenUrl.AddString(url)
	//Doing Breadth First Scraping therefore outstanding denotes the height of the tree

	for outstanding := 1; outstanding > 0; decreaseCounter() {
		fmt.Println("1.outstanding value is ", outstanding)

		u := <-urlQ
		if u.depth <= 1 {
			continue
		}
		for _, url := range u.urls {
			if isUrlPatternValid(url, pattern) && !seenUrl.TestString(url) {
				increaseCounter()
				seenUrl.AddString(url)
				wg.Add(1)
				scrape(fetcher, url, depth, &urlQ, &wg)

			}
		}

	}
	wg.Wait()
}

func increaseCounter() {
	mux.Lock()
	outstanding++
	mux.Unlock()
}

func decreaseCounter() {
	mux.Lock()
	outstanding--
	mux.Unlock()
}

func scrape(fetcher Fetcher, url string, depth int, q *chan deepUrl, wg *sync.WaitGroup) {
	defer wg.Done()

	//honor rate limit
	<-rateSemaphore
	_, urls, err := fetcher.Fetch(url)
	rateSemaphore <- struct{}{}

	if err != nil {
		log.Println("Error occurred while fetching the url :", url, err)
		return
	}

	if len(urls) > 0 {
		*q <- deepUrl{depth - 1, urls}
		//increaseCounter()
		fmt.Println("2.outstanding value is ", outstanding)
		fmt.Println("Picking up at depth and url and # of urls", url, depth, urls)
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
	var internalUrls []string
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
	urls := make([]string, 1)
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

				//log.Printf("found link %s", url)
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
