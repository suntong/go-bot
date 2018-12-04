package asynchronous

import (
	"errors"
	"fmt"
	"go-bot/utils"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/lexkong/log"
)

var huyaRE = regexp.MustCompile(`(?m)img\s*class\s*=\s*"\s*pic-con\s*[",']\s*src\s*=\s*[",'](.+?)[",']\s*[\s\S]+?liveRoomName\s*=\s*[",'](.+?)[",']\s*;[\s\S]+?ISLIVE\s*=\s*([a-z]+)\s*;[\s\S]+?TOPSID\s*=\s*[",'](.+?)[",']\s*;[\s\S]+?ANTHOR_NICK\s*=\s*'(.+?)[",']\s*;`)

func huyaOnline(addr string) interface{} {
	roomid := addr[strings.Index(addr, "-")+1:]
	req, err := http.NewRequest("GET", fmt.Sprintf("https://m.huya.com/%s", roomid), nil)
	if err != nil {
		log.Error("虎牙", err)
		return err
	}
	req.Header.Set("If-Modified-Since", "0")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", strings.Join([]string{
		"Mozilla/5.0",
		"(Linux; Android 6.0; Nexus 5 Build/MRA58N)",
		"AppleWebKit/537.36 (KHTML, like Gecko)",
		"Chrome/64.0.3282.140 Mobile Safari/537.36"}, ""))
	resp, err := client.Do(req)
	if err != nil {
		log.Error("虎牙", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("network error")
	}

	// 顺序 1 图片 2 title 3 状态 4 主播id 5 名字
	var out = make([]string, 6)
	r, _ := ioutil.ReadAll(resp.Body)
	for i, match := range huyaRE.FindStringSubmatch(string(r)) {
		out[i] = match
	}

	resp.Body.Close()
	// 开播
	if out[3] == "true" {
		m := utils.NewMessage()
		m.AddMsg(utils.CQshare(fmt.Sprintf("http://www.huya.com/%s", roomid),
			out[5], out[2], out[1]))
		return m.Message()
	} else {
		return nil
	}
}
