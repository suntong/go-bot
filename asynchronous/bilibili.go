package asynchronous

import (
	"errors"
	"fmt"
	"go-bot/utils"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/lexkong/log"
)

var blibilire = regexp.MustCompile(`(?m)room_id"\s*:\s*(\d+)\s*,[\s\S]+?user_cover\s*"\s*:\s*"([\s\S]+?)"\s*,[\s\S]+?uname\s*"\s*:\s*"([\s\S]+?)"\s*,[\s\S]+?live_status\s*"\s*:\s*(\d+)\s*[\s\S]+?title\s*"\s*:\s*"([\s\S]+?)"\s*,`)

func bilibiliOnline(addr string) interface{} {
	roomid := addr[strings.Index(addr, "-")+1:]
	req, err := http.NewRequest("GET", fmt.Sprintf("https://live.bilibili.com/%s", roomid), nil)
	if err != nil {
		log.Error("B站", err)
		return err
	}
	req.Header.Set("If-Modified-Since", "0")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		log.Error("B站", err)
		return err
	}
	// 顺序 1 房间号 2 图片 3 名字 4 状态 5 标题
	var out = make([]string, 6)
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("B站", err)
		return err
	}
	for i, match := range blibilire.FindStringSubmatch(string(r)) {
		out[i] = match
	}
	out[2], _ = strconv.Unquote(`"` + out[2] + `"`)
	// url解析
	resp.Body.Close()
	// 开播
	if out[4] == "1" {
		m := utils.NewMessage()
		m.AddMsg(utils.CQshare(fmt.Sprintf("https://live.bilibili.com/%s", out[1]),
			out[3], out[5], out[2]))
		return m.Message()
	} else {
		return nil
	}
	return errors.New("network error")
}
