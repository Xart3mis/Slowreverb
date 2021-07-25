package lib

import (
	"SlowReverb/lib/Secrets"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

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

func GetSong(title string, artist string, client *http.Client) *Result {
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

	fname := fmt.Sprintf("%s.%s", fmt.Sprintf("%s - %s", doc.Items[0].Snippet.Title, artist), "m4a")

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

	cmd := exec.Command("youtube-dl.exe", *doc.Items[0].ID.VideoID, "-x", "--audio-quality", "0",
		"--audio-format", "m4a", "--ffmpeg-location", "ffmpeg.exe", "-o", fname)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()

	return &Result{Response: &doc, SourceURL: &search, Filename: &fname}
}

func ModifySpeed(filename string, factor float64) string {

	if factor < 0.5 || factor > 2.0 {
		panic("speed factor must be between 0.5 and 2.0")
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

	cmd := exec.Command("ffmpeg.exe", "-i", filename, "-filter:a", fmt.Sprintf("atempo=%e", factor), "-vn",
		dir+endFname)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()

	return dir + endFname
}
