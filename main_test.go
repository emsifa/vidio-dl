package main

import (
	"os"
	"testing"
	"time"
)

func TestGetFilename(t *testing.T) {
	url := "https://www.vidio.com/watch/lorem-ipsum"
	filename := getFilename(url)
	expect := "lorem-ipsum"

	if filename != expect {
		t.Fatalf("Filename should be %s, got %s", expect, filename)
	}
}

func TestMakeTempDir(t *testing.T) {
	suffix := time.Now().Format("060102150405")
	dir := makeTempDir()
	path := Concat("vidio-dl-tmp-", suffix)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Directory %s is not created", path)
	}
	os.RemoveAll(dir)
}
