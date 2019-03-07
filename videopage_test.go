package main

import (
	"testing"
)

const url string = "https://www.vidio.com/watch/1601035-dilan-1991-sukses-pecahkan-rekor-luar-biasa-ini"

func TestGetPlaylistUrl(t *testing.T) {
	result := "https://www.vidio.com/videos/1601035/vjs_playlist.m3u8"
	page, _ := GetVideoPage(url)
	playlistUrl := page.GetPlaylistUrl()

	if playlistUrl != result {
		t.Errorf("URL playlist should be %s, got %s", result, playlistUrl)
	}
}

func TestGetStreamPlaylists(t *testing.T) {
	page, _ := GetVideoPage(url)
	playlists, _ := page.GetStreamPlaylists()

	if len(playlists) != 4 {
		t.Errorf("Stream playlists count should be %d, got %d", 4, len(playlists))
	}
}