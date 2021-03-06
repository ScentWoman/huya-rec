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
	Tags   []Tag     `json:"vMultiStreamInfo"`
}

// InfoAPI from hyPlayerConfig.stream.data .
type InfoAPI struct {
	LiveInfo struct {
		Nick string `json:"nick"`
		Room string `json:"roomName"`
	} `json:"gameLiveInfo"`
	StreamList []Stream `json:"gameStreamInfoList"`
}

// Tag ...
type Tag struct {
	Name string `json:"sDisplayName"`
}

// Stream ...
type Stream struct {
	Name      string `json:"sStreamName"`
	URL       string `json:"sFlvUrl"`
	Suffix    string `json:"sFlvUrlSuffix"`
	AntiCode  string `json:"sFlvAntiCode"`
	HURL      string `json:"sHlsUrl"`
	HSuffix   string `json:"sHlsUrlSuffix"`
	HAntiCode string `json:"sHlsAntiCode"`
	PURL      string `json:"sP2pUrl"`
	PSuffix   string `json:"sP2pUrlSuffix"`
	PAntiCode string `json:"newCFlvAntiCode"`
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

	// tag := "Tags:"
	// for k := range si.Tags {
	// 	tag = tag + " " + si.Tags[k].Name
	// }
	// log.Println(tag)

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

	// s.PURL = strings.TrimPrefix(s.PURL, "http://")
	// s.PURL = strings.TrimPrefix(s.PURL, "https://")

	// _, _ = sb.WriteString("https://")
	// _, _ = sb.WriteString(s.PURL)
	// _, _ = sb.WriteString("/")
	// _, _ = sb.WriteString(s.Name)
	// _, _ = sb.WriteString("_0_0_66.")
	// _, _ = sb.WriteString(s.PSuffix)
	// _, _ = sb.WriteString("?")
	// _, _ = sb.WriteString(html.UnescapeString(s.PAntiCode))
	// _, _ = sb.WriteString("&ex1=0&baseIndex=0&quickTime=2000&timeStamp=")
	// _, _ = sb.WriteString(time.Now().Format("2006-01-02_15:04:05.000"))
	// _, _ = sb.WriteString("&u=0&t=100&sv=2011191002")

	// tsb := strings.Builder{}
	// s.HURL = strings.TrimPrefix(s.HURL, "http://")
	// s.HURL = strings.TrimPrefix(s.HURL, "https://")

	// _, _ = tsb.WriteString("https://")
	// _, _ = tsb.WriteString(s.HURL)
	// _, _ = tsb.WriteString("/")
	// _, _ = tsb.WriteString(s.Name)
	// _, _ = tsb.WriteString(".")
	// _, _ = tsb.WriteString(s.HSuffix)
	// _, _ = tsb.WriteString("?")
	// _, _ = tsb.WriteString(html.UnescapeString(s.HAntiCode))
	// _, _ = tsb.WriteString("&ex1=0&baseIndex=0&quickTime=2000&timeStamp=")
	// _, _ = tsb.WriteString(time.Now().Format("2006-01-02_15:04:05.000"))
	// _, _ = tsb.WriteString("&u=0&t=100&sv=2011191002")
	// log.Println(tsb.String())

	return sb.String()
}
