package huya

import (
	"io/ioutil"
	"net/http"
	"time"
)

var (
	client = &http.Client{Timeout: 5 * time.Second}
)

func getPage(u string) (page string, e error) {
	req, e := http.NewRequest("GET", u, nil)
	if e != nil {
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36 Edg/87.0.664.47")

	resp, e := client.Do(req)
	if e != nil {
		return
	}

	body, e := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	page = string(body)

	return
}
