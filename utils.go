package main

import (
	"regexp"
	s "strings"

	"github.com/imroc/req"
)

func GetResponseBody(url string) (string, error) {
	r, err := req.Get(url)
	if err != nil {
		return "", err
	}

	body, err := r.ToString()
	if err != nil {
		return "", err
	}

	return body, nil
}

func DownloadFile(url string, dest string, onProgress func(curr int64, total int64)) {
	r, _ := req.Get(url, req.DownloadProgress(onProgress))
	r.ToFile(dest)
}

func IsUrl(str string) bool {
	isUrl, _ := regexp.MatchString("^https?:\\/\\/", str)
	return isUrl
}

func Concat(strs ...string) string {
	var sb s.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
