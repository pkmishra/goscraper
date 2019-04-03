package goscraper

import (
	"io"
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

func Test_parseBaseURL(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"given valid url result should be valid url", args{"http://www.google.com/test"}, "http://www.google.com"},
		{"given valid long url, result should be valid url", args{"https://yahoo.co.in"}, "https://yahoo.co.in"},
		{"given invalid url, result should be empty string", args{"test"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBaseURL(tt.args.u); got != tt.want {
				t.Errorf("parseBaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resolveRelative(t *testing.T) {
	type args struct {
		baseURL string
		hrefs   []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"given array of relative urls should return full valid urls", args{"http://www.wikipedia.org", []string{"test", ""}}, []string{"http://www.wikipedia.org/test", "http://www.wikipedia.org"}},
		{"given array of relative and full url should return full valid urls", args{"http://www.wikipedia.org", []string{"http://www.wikipedia.org/test", ""}}, []string{"http://www.wikipedia.org/test", "http://www.wikipedia.org"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveRelative(tt.args.baseURL, tt.args.hrefs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveRelative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractUrl(t *testing.T) {
	type args struct {
		body io.ReadCloser
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractUrl(tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isUrlPatternValid(t *testing.T) {
	type args struct {
		url     string
		pattern string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUrlPatternValid(tt.args.url, tt.args.pattern); got != tt.want {
				t.Errorf("isUrlPatternValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goScraper(t *testing.T) {
	//Run(HttpFetcher{}, "https://pkmishra.github.io/blog/2015/04/22/leading-team-tips/", 1, `.*pkmishra.github.io.*`, 2)
	Run(HttpFetcher{}, "https://en.wikipedia.org/wiki/Free_content", 2, `.*en.wikipedia.org.*`, 2)
	//Run(fetcher, "https://golang.org/", 1, "", 2)
	//Run("https://gist.github.com/chriswhitcombe/3d5684a1eb0d9ae8adac",1,"",2)
}

func Test_getHref(t *testing.T) {
	type args struct {
		t html.Token
	}
	tests := []struct {
		name     string
		args     args
		wantOk   bool
		wantHref string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotHref := getHref(tt.args.t)
			if gotOk != tt.wantOk {
				t.Errorf("getHref() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotHref != tt.wantHref {
				t.Errorf("getHref() gotHref = %v, want %v", gotHref, tt.wantHref)
			}
		})
	}
}

// fakeFetcher is Fetcher that returns canned results.
//type fakeFetcher map[string]*fakeResult
//
//type fakeResult struct {
//	body string
//	urls []string
//}
//
//func (f fakeFetcher) Fetch(url string) (string, []string, error) {
//	if res, ok := f[url]; ok {
//		return res.body, res.urls, nil
//	}
//	return "", nil, fmt.Errorf("not found: %s", url)
//}
//
//// fetcher is a populated fakeFetcher.
//var fetcher = fakeFetcher{
//	"https://golang.org/": &fakeResult{
//		"The Go Programming Language",
//		[]string{
//			"https://golang.org/pkg/",
//			"https://golang.org/cmd/",
//		},
//	},
//	"https://golang.org/pkg/": &fakeResult{
//		"Packages",
//		[]string{
//			"https://golang.org/",
//			"https://golang.org/cmd/",
//			"https://golang.org/pkg/fmt/",
//			"https://golang.org/pkg/os/",
//		},
//	},
//	"https://golang.org/pkg/fmt/": &fakeResult{
//		"Package fmt",
//		[]string{
//			"https://golang.org/",
//			"https://golang.org/pkg/",
//		},
//	},
//	"https://golang.org/pkg/os/": &fakeResult{
//		"Package os",
//		[]string{
//			"https://golang.org/",
//			"https://golang.org/pkg/",
//		},
//	},
//}
