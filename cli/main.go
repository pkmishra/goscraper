package main

import (
	"flag"
	"goscraper"
	"os"
)

func main() {

	url := flag.String("url", "", "example https://en.wikipedia.org/wiki/Slope_One")
	depth := flag.Int("depth", 1, "integer value e.g. 1")
	pattern := flag.String("pattern", "", "regex pattern to extract link e.g. ")
	rate := flag.Int("rate", 1, "requests per seconds allowed;default is 1 and max is 5")
	flag.Parse()
	if *url == "" || *depth < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *rate < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	goscraper.Run(goscraper.HttpFetcher{}, *url, *depth, *pattern, *rate)
}
