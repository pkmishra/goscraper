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
