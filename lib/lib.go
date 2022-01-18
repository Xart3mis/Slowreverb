package lib

import (
	"SlowReverb/lib/Secrets"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type document struct {
	Etag  string `json:"etag"`
	Items []struct {
		Etag string `json:"etag"`
		ID   struct {
			Kind       string  `json:"kind"`
			PlaylistID *string `json:"playlistId,omitempty"`
			VideoID    *string `json:"videoId,omitempty"`
		} `json:"id"`
		Kind    string `json:"kind"`
		Snippet struct {
			ChannelID            string    `json:"channelId"`
			ChannelTitle         string    `json:"channelTitle"`
			Description          string    `json:"description"`
			LiveBroadcastContent string    `json:"liveBroadcastContent"`
			PublishTime          time.Time `json:"publishTime"`
			PublishedAt          time.Time `json:"publishedAt"`
			Thumbnails           struct {
				Default heightURLWidth `json:"default"`
				High    heightURLWidth `json:"high"`
				Medium  heightURLWidth `json:"medium"`
			} `json:"thumbnails"`
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
	Kind          string `json:"kind"`
	NextPageToken string `json:"nextPageToken"`
	PageInfo      struct {
		ResultsPerPage int `json:"resultsPerPage"`
		TotalResults   int `json:"totalResults"`
	} `json:"pageInfo"`
	RegionCode string `json:"regionCode"`
}
type Release struct {
	URL             string    `json:"url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	HTMLURL         string    `json:"html_url"`
	ID              int       `json:"id"`
	Author          Author    `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []Assets  `json:"assets"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
}
type Author struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}
type Uploader struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}
type Assets struct {
	URL                string    `json:"url"`
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	Uploader           Uploader  `json:"uploader"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
}
type heightURLWidth struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}
type Result struct {
	Response  *document
	SourceURL *string
	Filename  *string
}

type rvrbtypes struct {
	Hall struct {
		Large_Hall  string
		Medium_Hall string
		Small_Hall  string
	}
	Chamber struct {
		Large_Chamber  string
		Medium_Chamber string
		Small_Chamber  string
		Vocal_Chamber  string
	}
}

func ReverbTypes() rvrbtypes {
	return rvrbtypes{
		Hall: struct {
			Large_Hall  string
			Medium_Hall string
			Small_Hall  string
		}{
			Large_Hall:  "1 Halls 01 Large Hall  M-to-S.wav",
			Medium_Hall: "1 Halls 02 Medium Hall  M-to-S.wav",
			Small_Hall:  "1 Halls 03 Small Hall  M-to-S.wav",
		},
		Chamber: struct {
			Large_Chamber  string
			Medium_Chamber string
			Small_Chamber  string
			Vocal_Chamber  string
		}{
			Large_Chamber:  "4 Chambers 01 Large Chamber  M-to-S.wav",
			Medium_Chamber: "4 Chambers 02 Medium Chamber  M-to-S.wav",
			Small_Chamber:  "4 Chambers 03 Small Chamber  M-to-S.wav",
			Vocal_Chamber:  "4 Chambers 10 Vocal Chamber  M-to-S.wav",
		},
	}
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}

}

const (
	STDERRCH = iota
	STDOUTCH = iota
)

const (
	BACKSPACE  = "\r"
	YELLOW_FG  = "\033[33m"
	GREEN_FG   = "\033[32m"
	MAGENTA_FG = "\033[35m"
	CYAN_FG    = "\033[36m"
	RESET      = "\033[0m"
)

func Init(timeout int) *http.Client {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.MaxIdleConns = 100
	customTransport.MaxConnsPerHost = 100
	customTransport.MaxIdleConnsPerHost = 100
	client := http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: customTransport,
	}
	CheckAllDeps(&client)
	return &client
}

func GetSong(title string, artist string, client *http.Client, verbose ...bool) *Result {
	vb := false
	if len(verbose) > 0 {
		vb = verbose[0]
	}

	doc := document{}
	search := Secrets.SearchUrl
	search = fmt.Sprintf(search, url.QueryEscape(fmt.Sprintf("%s - %s", title, artist)), Secrets.Key)
	req, err := http.NewRequest(http.MethodGet, search, nil)
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

	fname := fmt.Sprintf("%s.%s", fmt.Sprintf("%s - %s", title, artist), "m4a")

	os.Chdir(DependencyPath)

	dir, _ := os.Getwd()
	dT := strings.Split(dir, "\\")

	dT = dT[:len(dT)-2]
	dir = ""

	for i := range dT {
		dir += dT[i] + "/"
	}

	dir += "Temp/"
	fname = dir + fname

	ch := make(chan string)
	go func() {
		RunCommandCh(ch, "\r\n", STDOUTCH, "youtube-dl.exe", *doc.Items[0].ID.VideoID, "-x", "--audio-quality", "0",
			"--audio-format", "m4a", "--ffmpeg-location", "ffmpeg.exe", "-o", fname)
	}()

	downloadRGX := regexp.MustCompile(`\[download\]\s+\d`)
	for v := range ch {
		match := downloadRGX.MatchString(v)
		if vb && !match {
			fmt.Print(BACKSPACE + YELLOW_FG + v + RESET)
		}
		if match {
			fmt.Print(BACKSPACE + MAGENTA_FG + v + RESET)
		}
	}
	fmt.Println()

	return &Result{Response: &doc, SourceURL: &search, Filename: &fname}
}

func ModifySpeed(filename string, factor float64, verbose ...bool) string {
	vb := false
	if len(verbose) > 0 {
		vb = verbose[0]
	}

	os.Chdir(DependencyPath)
	dir, _ := os.Getwd()
	dT := strings.Split(dir, "\\")

	dT = dT[:len(dT)-2]
	dir = ""

	for i := range dT {
		dir += dT[i] + "/"
	}

	dir += "Output/"

	endFnameA := strings.Split(filename, "/")
	endFname := endFnameA[len(endFnameA)-1]
	endFname = "slwd_" + endFname

	ch := make(chan string)
	go func() {
		RunCommandCh(ch, "\r\n", STDERRCH, "ffmpeg.exe", "-y", "-i", filename, "-filter:a", fmt.Sprintf("atempo=%e", factor), "-vn", dir+endFname)
	}()

	timeRGX := regexp.MustCompile(`size=[\s\d]+.B\stime=[\d:\.]+\sbitrate=.+\/s\sspeed=.+x`)
	for v := range ch {
		match := timeRGX.MatchString(v)
		if vb && !match {
			fmt.Print(BACKSPACE + YELLOW_FG + v + RESET)
		}
		if match {
			fmt.Print(BACKSPACE + GREEN_FG + v + RESET)
		}
	}
	fmt.Println()

	return dir + endFname
}

func RunCommandCh(stdoutCh chan<- string, cutset string, stdwch int, command string, flags ...string) error {
	cmd := exec.Command(command, flags...)
	cmd.Stdin = os.Stdin
	var output io.ReadCloser
	var err error

	switch stdwch {
	case STDERRCH:
		output, err = cmd.StderrPipe()

	case STDOUTCH:
		output, err = cmd.StdoutPipe()

	default:
		output, err = cmd.StdoutPipe()
	}
	if err != nil {
		return fmt.Errorf("RunCommand: cmd.StdoutPipe(): %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("RunCommand: cmd.Start(): %v", err)
	}

	go func() {
		defer close(stdoutCh)
		for {
			buf := make([]byte, 1024)
			n, err := output.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Fatal(err)
				}
				if n == 0 {
					break
				}
			}
			text := strings.TrimSpace(string(buf[:n]))
			for {
				n := strings.IndexAny(text, cutset)
				if n == -1 {
					if len(text) > 0 {
						stdoutCh <- text
					}
					break
				}
				stdoutCh <- text[:n]
				if n == len(text) {
					break
				}
				text = text[n+1:]
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("RunCommand: cmd.Wait(): %v", err)
	}
	return nil
}

func Reverberize(filename string, dryness int, wetness int, mix_ratio int, reverb_type string, verbose ...bool) string {
	vb := false
	if len(verbose) > 0 {
		vb = verbose[0]
	}

	if dryness > 10 {
		panic("value for dryness must be 0-10")
	}
	if wetness > 10 {
		panic("value for wetness must be 0-10")
	}
	os.Chdir(DependencyPath)
	dir, _ := os.Getwd()
	dT := strings.Split(dir, "\\")

	dT = dT[:len(dT)-2]
	dir = ""

	for i := range dT {
		dir += dT[i] + "/"
	}

	dir += "Output/"

	endFnameA := strings.Split(filename, "/")
	endFname := endFnameA[len(endFnameA)-1]
	endFname = "rvrbrzd_" + endFname

	ch := make(chan string)
	go func() {
		RunCommandCh(ch, "\r\n", STDERRCH, "ffmpeg.exe", "-y", "-i", filename, "-i", fmt.Sprintf("../IRAF/%s", reverb_type), "-filter_complex",
			fmt.Sprintf("[0] [1] afir=dry=%d:wet=%d [reverb]; [0] [reverb] amix=inputs=2:weights=%d 1", dryness, wetness, mix_ratio), dir+endFname)
	}()

	timeRGX := regexp.MustCompile(`size=[\s\d]+.B\stime=[\d:\.]+\sbitrate=.+\/s\sspeed=.+x`)
	for v := range ch {
		match := timeRGX.MatchString(v)
		if vb && !match {
			fmt.Print(BACKSPACE + YELLOW_FG + v + RESET)
		}
		if match {
			fmt.Print(BACKSPACE + GREEN_FG + v + RESET)
		}
	}
	fmt.Println()

	return dir + endFname
}

func Play(filename string, finished chan<- bool, verbose ...bool) {
	vb := false
	if len(verbose) > 0 {
		vb = verbose[0]
	}

	os.Chdir(DependencyPath)

	ch := make(chan string)
	go func() {
		RunCommandCh(ch, "\r\n", STDERRCH, "ffplay.exe", "-nodisp", "-autoexit", filename)
	}()

	timeRGX := regexp.MustCompile(`\d+\.\d+\sM-A`)
	for v := range ch {
		match := timeRGX.MatchString(v)
		if vb && !match {
			fmt.Print(BACKSPACE + YELLOW_FG + v + RESET)
		}
		if match {
			fmt.Print(BACKSPACE + CYAN_FG + v + RESET)
		}
	}
	fmt.Println()
	finished <- true
}
