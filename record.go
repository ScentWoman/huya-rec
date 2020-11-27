package rec

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ScentWoman/huya-rec/huya"
)

var (
	session = &http.Client{}
)

// Record ...
func Record(room, src string, split, retry time.Duration, output string) {
	for {
		e := recOnce(room, src, split, output)
		if e != nil {
			log.Println(e)
		}
		time.Sleep(retry)
	}
}

func recOnce(room, src string, split time.Duration, output string) (e error) {
	info, e := huya.GetInfo(room)
	if e != nil {
		return
	}
	if !info.On {
		return
	}

	log.Printf("\"%v\" is on live: \"%v\"\n", info.Name, info.Title)

	ctx, cancel := context.WithTimeout(context.Background(), split)
	defer cancel()

	var resp *http.Response
	respChan := make(chan *http.Response)

	for _, v := range info.Stream {
		go func(u string) {
			defer func() {
				if p := recover(); p != nil {
					log.Println(p.(error))
				}
			}()

			if !strings.Contains(u, src+".flv.huya.com") {
				return
			}

			req, e := http.NewRequestWithContext(ctx, "GET", u, nil)
			if e != nil {
				panic(e)
			}
			req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36 Edg/87.0.664.47")

			resp, e := session.Do(req)
			if e != nil {
				panic(e)
			}

			select {
			case respChan <- resp:
			default:
				resp.Body.Close()
			}
		}(v)
	}

	select {
	case resp = <-respChan:
	case <-time.Tick(5 * time.Second):
		log.Println("Failed to establish connection in 5 seconds...")
		return
	}

	defer resp.Body.Close()

	fw, e := os.Create(filepath.Join(output, legalFilename(info.Title+"_"+time.Now().Format("15_04_05"))+".flv"))
	if e != nil {
		return
	}
	defer fw.Close()

	_, e = io.Copy(fw, resp.Body)

	return
}
