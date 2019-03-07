package main

import (
	"regexp"
	s "strings"
)

type VideoPage struct {
	url         string
	playlistUrl string
}

func GetVideoPage(url string) (VideoPage, error) {
	playlistUrl, err := getPlaylistUrl(url)
	if err != nil {
		return VideoPage{url: "", playlistUrl: ""}, err
	}

	vp := VideoPage{url, playlistUrl}

	return vp, nil
}

func (vp VideoPage) GetPlaylistUrl() string {
	return vp.playlistUrl
}

func (vp VideoPage) GetStreamPlaylists() ([]StreamPlaylist, error) {
	res, err := GetResponseBody(vp.playlistUrl)
	if err != nil {
		return make([]StreamPlaylist, 0), err
	}

	lines := s.Split(res, "\n")

	playlists := make([]StreamPlaylist, 0)
	isPlaylistInfo := false
	for _, line := range lines {
		if !isPlaylistInfo {
			isPlaylistInfo, _ = regexp.MatchString("^#EXT-X-STREAM-INF", line)

			if isPlaylistInfo {
				playlists = append(playlists, StreamPlaylist{
					Name:       extractLineInfo(line, "NAME=\"(\\d+p)"),
					Bandwidth:  extractLineInfo(line, "BANDWIDTH=(\\d+)"),
					Resolution: extractLineInfo(line, "RESOLUTION=(\\d+x\\d+)"),
					Url:        "",
				})
			}
		} else {
			playlists[len(playlists)-1].Url = line
			isPlaylistInfo = false
		}
	}

	return playlists, nil
}

func extractLineInfo(line string, regex string) string {
	rg, _ := regexp.Compile(regex)
	match := rg.FindStringSubmatch(line)

	if len(match) == 2 {
		return match[1]
	} else {
		return ""
	}
}

func getPlaylistUrl(url string) (string, error) {
	res, err := GetResponseBody(url)
	if err != nil {
		return "", err
	}

	r, _ := regexp.Compile("https\\:\\/\\/www.vidio.com\\/videos\\/\\d+/vjs_playlist\\.m3u8")

	playlistUrl := r.FindString(res)

	return playlistUrl, nil
}
