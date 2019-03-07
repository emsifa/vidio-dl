package main

import (
	"regexp"
	s "strings"
)

type StreamPlaylist struct {
	Name 		string
	Resolution 	string
	Bandwidth 	string
	Url 		string
}

func (sp StreamPlaylist) GetStreamUrls() ([]string, error) {
	res, err := GetResponseBody(sp.Url)
	if err != nil {
		return make([]string, 0), err
	}

	lines := s.Split(res, "\n")

	urls := make([]string, 0)
	isPlaylistUrl := false
	for _, line := range lines {
		if !isPlaylistUrl {
			isPlaylistUrl, _ = regexp.MatchString("^#EXTINF:", line)
		} else {
			url := resolveUrl(line, sp.Url)
			urls = append(urls, url)
			isPlaylistUrl = false
		}
	}

	return urls, nil
}

func resolveUrl(line string, playlistUrl string) string {
	isUrl := IsUrl(line)
	if isUrl {
		return line
	}

	split := s.Split(playlistUrl, "/")
	split = split[:len(split)-1]
	split = append(split, line)

	return s.Join(split, "/")
}