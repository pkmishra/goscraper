package main

import (
	"flag"
	"gowatch/crawler"
	"os"
)

func main() {

	url := flag.String("url", "", "example https://en.wikipedia.org/wiki/Slope_One")
	depth := flag.Int("depth", 1, "integer value e.g. 1")
	pattern := flag.String("pattern", "", "regex pattern to extract link e.g. ")
	flag.Parse()
	if *url == "" || *depth < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	crawler.Run(*url, *depth, *pattern)
}
