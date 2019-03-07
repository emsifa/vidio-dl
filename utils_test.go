package main

import (
	"testing"
)

func TestIsUrl(t *testing.T) {
	validUrls := []string{
		"https://www.google.com",
		"http://www.google.com",
		"http://localhost",
		"https://192.168.100.1/foo/bar/?baz=qux",
	}

	for _, url := range validUrls {
		if IsUrl(url) == false {
			t.Errorf("%s should be valid URL", url)
		}
	}

	invalidUrls := []string{
		"lorem-ipsum",
		"mailto:foo@mail.com",
		"http:/invalid.com",
		"http:invalid.com",
	}
	
	for _, url := range invalidUrls {
		if IsUrl(url) == true {
			t.Errorf("%s should be invalid URL", url)
		}
	}
}

func TestConcat(t *testing.T) {
	var a, b string

	a = Concat("foo", "bar")
	b = "foobar"
	if a != b {
		t.Errorf("Concat(\"foo\", \"bar\") should be \"%s\", got \"%s\"", b, a)
	}

	a = Concat("foo", "bar", "baz", "qux")
	b = "foobarbazqux"
	if a != b {
		t.Errorf("Concat(\"foo\", \"bar\", \"baz\", \"qux\") should be \"%s\", got \"%s\"", b, a)
	}
}