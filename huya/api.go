package huya

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"
)

// Info illustrates a live room.
type Info struct {
	On     bool
	Name   string
	Title  string
	Stream []string
}

// StreamAPI illustrates huya hyPlayerConfig.stream .
type StreamAPI struct {
	Status int       `json:"status"`
	Msg    string    `json:"msg"`
	Data   []InfoAPI `json:"data"`
}

// InfoAPI from hyPlayerConfig.stream.data .
type InfoAPI struct {
	LiveInfo struct {
		Nick string `json:"nick"`
		Room string `json:"roomName"`
	} `json:"gameLiveInfo"`
	StreamList []Stream `json:"gameStreamInfoList"`
}

// Stream ...
type Stream struct {
	Name     string `json:"sStreamName"`
	URL      string `json:"sFlvUrl"`
	Suffix   string `json:"sFlvUrlSuffix"`
	AntiCode string `json:"sFlvAntiCode"`
}

var (
	rePC     = regexp.MustCompile(`(?s)hyPlayerConfig.*?};`)
	reStream = regexp.MustCompile(`"stream".*};`)
)

// GetInfo gets information of given url live room.
func GetInfo(u string) (i Info, e error) {
	page, e := getPage(u)
	if e != nil {
		return
	}

	pc := rePC.FindString(page)
	stream := reStream.FindString(pc)

	stream = trimStream(stream)
	if len(stream) < 10 {
		return
	}

	i.On = true
	sb, e := base64.StdEncoding.DecodeString(stream)

	if e != nil {
		return
	}

	var si StreamAPI
	e = json.Unmarshal(sb, &si)
	if e != nil {
		return
	}

	if si.Status != 200 {
		e = fmt.Errorf("huya: %v", si.Msg)
		return
	}

	if len(si.Data) == 0 {
		e = fmt.Errorf("huya: null gameStreamInfoList")
		return
	}

	data := si.Data[0]

	i.Name = data.LiveInfo.Nick
	i.Title = data.LiveInfo.Room

	for _, v := range data.StreamList {
		i.Stream = append(i.Stream, genStreamURL(v))
	}

	return
}

func trimStream(s string) string {
	for len(s) > 0 && (s[len(s)-1] != ' ') {
		s = s[:len(s)-1]
	}
	for len(s) > 0 && (s[0] != ' ') {
		s = s[1:]
	}
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")
	return s
}

// sFlvUrl + "/" + sStreamName + "." + sFlvUrlSuffix + "?" + sFlvAntiCode + "&ex1=0&baseIndex=0&quickTime=2000&timeStamp="+ {2020-11-27_13:38:59.307} + "&u=0&t=100&sv=2011191002"
func genStreamURL(s Stream) string {
	sb := strings.Builder{}
	s.URL = strings.TrimPrefix(s.URL, "http://")
	s.URL = strings.TrimPrefix(s.URL, "https://")

	_, _ = sb.WriteString("https://")
	_, _ = sb.WriteString(s.URL)
	_, _ = sb.WriteString("/")
	_, _ = sb.WriteString(s.Name)
	_, _ = sb.WriteString(".")
	_, _ = sb.WriteString(s.Suffix)
	_, _ = sb.WriteString("?")
	_, _ = sb.WriteString(html.UnescapeString(s.AntiCode))
	_, _ = sb.WriteString("&ex1=0&baseIndex=0&quickTime=2000&timeStamp=")
	_, _ = sb.WriteString(time.Now().Format("2006-01-02_15:04:05.000"))
	_, _ = sb.WriteString("&u=0&t=100&sv=2011191002")

	return sb.String()
}
