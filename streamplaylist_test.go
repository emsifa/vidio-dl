package main

import (
	"testing"
)

func TestGetStreamUrls(t *testing.T) {
	sp := StreamPlaylist{
		Name: "640p",
		Resolution: "640x360",
		Url: "https://cdn0-a.production.vidio.static6.com/uploads/1601035/ets-dilan-201991-20sukses-20pecahkan-20rekor-20luar-20biasa-20ini-1229-b900.mp4.m3u8",
		Bandwidth: "900000",
	}

	urls, _ := sp.GetStreamUrls()

	if len(urls) != 12 {
		t.Errorf("Stream urls count should be %d, got %d", 12, len(urls))
	}

	for i, url := range urls {
		if IsUrl(url) == false {
			t.Errorf("Stream URL at index %d is not URL, %s", i, url)
		}
	}
}