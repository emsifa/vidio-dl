package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	s "strings"
	"syscall"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "vidio-dl"
	app.Usage = "Download video from vidio.com"
	app.Version = "0.1.0"
	app.UsageText = "vidio-dl [url-video]"
	app.Action = downloadCmd

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func downloadCmd(cmd *cli.Context) error {
	url := cmd.Args().Get(0)
	if url == "" {
		return cli.NewExitError("URL is required", 0)
	}

	page, err := GetVideoPage(url)
	if err != nil {
		return cli.NewExitError("Failed to fetch page information", 1)
	}

	playlist, err := askResolution(page)
	if err != nil {
		return cli.NewExitError("Failed to fetch resolutions", 2)
	}

	output := Concat(getFilename(url), " (", playlist.Name, ").mp4")
	replace := false
	if _, err := os.Stat(output); err == nil {
		replace = askReplace(output)
		if !replace {
			os.Exit(0)
		}
	}

	tmpDir := makeTempDir()
	files, err := downloadStreamFiles(playlist, tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		return cli.NewExitError("Something wrong while downloading file", 3)
	}

	if replace {
		os.Remove(output)
	}

	err = combineFiles(files, output)
	if err != nil {
		os.RemoveAll(tmpDir)
		return cli.NewExitError("Failed to combine files", 4)
	}

	fmt.Printf("Saved: %s", output)

	os.RemoveAll(tmpDir)
	return nil
}

func askReplace(output string) bool {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U000027A4 {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "\U000027A4 {{ . | blue | cyan }}",
	}

	prompt := promptui.Select{
		Label:     fmt.Sprintf("File %s already exists. Do you want to replace it", output),
		Items:     []string{"Yes", "No"},
		Templates: templates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return false
	}

	return i == 0
}

func askResolution(page VideoPage) (StreamPlaylist, error) {
	playlists, err := page.GetStreamPlaylists()
	if err != nil {
		return StreamPlaylist{}, err
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U000027A4 {{ .Name | cyan }} ({{ .Resolution | blue }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Resolution | blue }})",
		Selected: "\U000027A4 {{ .Name | blue | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Resolution",
		Items:     playlists,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return StreamPlaylist{}, err
	}

	return playlists[i], nil
}

func downloadStreamFiles(playlist StreamPlaylist, tmpDir string) ([]string, error) {
	urls, err := playlist.GetStreamUrls()
	if err != nil {
		return make([]string, 0), err
	}

	files := make([]string, 0)

	// Remove downloaded files on Ctrl+C
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.RemoveAll(tmpDir)
		fmt.Println("aborted")
		os.Exit(0)
	}()

	// Start download files
	fmt.Println("")
	fmt.Println("Downloading chunk files")
	bar := pb.StartNew(len(urls))
	for _, url := range urls {
		filename := getFilename(url)
		filepath := Concat(tmpDir, "/", filename)
		files = append(files, filepath)
		DownloadFile(url, filepath, func(curr int64, total int64) {})
		bar.Increment()
	}
	bar.FinishPrint("")

	return files, nil
}

func combineFiles(files []string, output string) error {
	out, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, file := range files {
		in, err := os.Open(file)
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
		in.Close()
	}

	out.Close()

	return nil
}

func getFilename(url string) string {
	split := s.Split(url, "/")
	last := split[len(split)-1]

	return last
}

func makeTempDir() string {
	dir := Concat("vidio-dl-tmp-", time.Now().Format("060102150405"))

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}

	return dir
}
