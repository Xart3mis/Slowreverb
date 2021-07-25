package lib

import (
	"SlowReverb/lib/Secrets"
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

const DependencyPath string = `dependencies/bin/`
const ytdl = "youtube-dl.exe"
const ffmpeg = "ffmpeg.zip"

func CheckAllDeps(client *http.Client) {
	checkFolders()
	checkFFmpeg(client)
	checkYoutubeDl(client)
	fmt.Println("Dependencies Installed")
}

type writeCounter struct {
	Total uint64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc writeCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func checkFolders() {
	if !exists(`dependencies/bin/`) {
		os.MkdirAll("dependencies/bin/", os.ModePerm)
	}
	if !exists(`Output/`) {
		os.MkdirAll("Output/", os.ModePerm)
	}
}

func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	counter := &writeCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	fmt.Print("\n")

	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

func getYoutubeDl(client *http.Client) {
	url := Secrets.YoutubeDlRelease
	doc := Release{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	CheckErr(err)

	req.Header.Set("User-Agent", "SlowReverb-App")

	resp, getErr := client.Do(req)
	CheckErr(getErr)

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	CheckErr(readErr)

	jsonErr := json.Unmarshal([]byte(body), &doc)
	CheckErr(jsonErr)

	ch := make(chan string)
	go func() {
		for _, asset := range doc.Assets {
			if asset.Name == "youtube-dl.exe" {
				ch <- asset.BrowserDownloadURL
			}
		}
	}()

	dlErr := downloadFile(DependencyPath+ytdl, <-ch)
	fmt.Println("Download Finished")
	CheckErr(dlErr)
}

func checkYoutubeDl(client *http.Client) {
	ch := make(chan bool)
	go func(name string) {
		if _, err := os.Stat(name); err != nil {
			if os.IsNotExist(err) {
				ch <- false
			}
		}
		ch <- true
	}(DependencyPath + ytdl)

	if !<-ch {
		fmt.Println("Youtube-dl dependency not found.")
		getYoutubeDl(client)
	}
}

func getFFmpeg(client *http.Client) {
	url := Secrets.FFmpegRelease
	doc := Release{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	CheckErr(err)

	req.Header.Set("User-Agent", "SlowReverb-App")

	resp, getErr := client.Do(req)
	CheckErr(getErr)

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	CheckErr(readErr)

	jsonErr := json.Unmarshal([]byte(body), &doc)
	CheckErr(jsonErr)

	ch := make(chan string)
	go func() {
		for _, asset := range doc.Assets {
			sl := strings.Split(asset.Name, "-")
			buildType := sl[len(sl)-1]
			if buildType == "essentials_build.zip" {
				ch <- asset.BrowserDownloadURL
			}
		}
	}()

	dlErr := downloadFile(`dependencies/`+ffmpeg, <-ch)
	fmt.Println("Download Finished")
	CheckErr(dlErr)
	extractFFmpeg(`dependencies/` + ffmpeg)
}

func checkFFmpeg(client *http.Client) {
	ch := make(chan bool)
	go func(name string) {
		if _, err := os.Stat(name); err != nil {
			if os.IsNotExist(err) {
				ch <- false
			}
		}
		ch <- true
	}(`dependencies/bin/` + "ffmpeg.exe")

	if !<-ch {
		fmt.Println("FFmpeg dependency not found.")
		getFFmpeg(client)
	}
}
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
func extractFFmpeg(yourZipFile string) error {
	tmpDir, err := ioutil.TempDir(os.Getenv("TEMP"), "SR-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	if filenames, err := Unzip(yourZipFile, tmpDir); err != nil {
		return err
	} else {
		for _, path := range filenames {
			fmt.Println("Extracted: " + path)
			sl := strings.Split(path, `\`)
			file := sl[len(sl)-1]
			copied := 0
			if (file == "ffmpeg.exe" || file == "ffplay.exe" || file == "ffprobe.exe") && copied <= 3 {
				copy(path, `dependencies/bin/`+file)
				fmt.Println("Copied " + file)
				copied += 1
			} else if copied > 3 {
				break
			}
		}
	}
	return nil
}

func Unzip(src string, destination string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)

	if err != nil {

		return filenames, err
	}

	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destination, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s is an illegal filepath", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			f.Mode())

		if err != nil {

			return filenames, err
		}

		rc, err := f.Open()

		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}

	return filenames, nil
}
