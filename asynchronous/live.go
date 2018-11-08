package asynchronous

import (
	"encoding/json"
	"fmt"
	"go-bot/pkg/memory"
	"go-bot/pkg/message"
	"go-bot/utils"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lexkong/log"
)

// 修改为并发安全
var (
	registerID = make(map[string]string)
	blibilire  = regexp.MustCompile(`(?m)room_id"\s*:\s*(\d+)\s*,[\s\S]+?user_cover\s*"\s*:\s*"([\s\S]+?)"\s*,[\s\S]+?uname\s*"\s*:\s*"([\s\S]+?)"\s*,[\s\S]+?live_status\s*"\s*:\s*(\d+)\s*[\s\S]+?title\s*"\s*:\s*"([\s\S]+?)"\s*,`)
	re         = regexp.MustCompile(`(?m)picURL\s*=\s*'(.+?)'\s*[\s\S]+?liveRoomName\s*=\s*'(.+?)'\s*;[\s\S]+?ISLIVE\s*=\s*([a-z]+)\s*;[\s\S]+?TOPSID\s*=\s*'(.+?)'\s*;[\s\S]+?ANTHOR_NICK\s*=\s*'(.+?)'\s*;`)
)

func Live() {
	roomID, _ := memory.GetLive("liveRoom").Range()
	client := &http.Client{}
	for i := range roomID {
		l := strings.Split(roomID[i], "-")
		switch l[0] {
		case "斗鱼":
			var result struct {
				Data struct {
					Roomid     string `json:"room_id"`
					RoomName   string `json:"room_name"`
					RoomStatus string `json:"room_status"`
					Avatar     string `json:"avatar"`
					Owner      string `json:"owner_name"`
				} `json:"data"`
			}

			req, err := http.NewRequest("GET", fmt.Sprintf("http://open.douyucdn.cn/api/RoomApi/room/%s", l[1]), nil)
			if err != nil {
				log.Error("douyu", err)
				break
			}
			req.Header.Set("If-Modified-Since", "0")
			req.Header.Set("Cache-Control", "no-cache")
			resp, err := client.Do(req)
			if err != nil {
				log.Error("douyu", err)
				break
			}
			r, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(r, &result)
			// 开播
			if result.Data.RoomStatus == "1" {
				m := utils.NewMessage()
				m.AddMsg(utils.CQshare(fmt.Sprintf("https://www.douyu.com/%s", result.Data.Roomid),
					result.Data.Owner, result.Data.RoomName, result.Data.Avatar))
				group, err := memory.GetLive(roomID[i]).Range()
				if err != nil {
					log.Error("douyu", err)
				}
				for _, v := range group {
					if registerID[strings.Join([]string{"斗鱼", l[1], v}, "-")] != result.Data.RoomStatus {
						tmp, _ := strconv.ParseInt(v, 10, 64)
						memory.DefaultMes.Push(
							message.SendMsg(message.MSG_GROUP, tmp, utils.NewMessage().
								AddMsg(utils.CQat("all")).
								AddMsg(utils.CQtext(fmt.Sprintf("[%s]开播了！！！", result.Data.Owner))).
								Message(), false, ""),
						)

						memory.DefaultMes.Push(
							message.SendMsg(message.MSG_GROUP, tmp, m.Message(), false, ""),
						)
						registerID[strings.Join([]string{"斗鱼", l[1], v}, "-")] = result.Data.RoomStatus
					}
				}
			}
		case "熊猫":
			var result struct {
				Data struct {
					Info struct {
						ID     string `json:"id"`
						Name   string `json:"name"`
						Status string `json:"status"`
					} `json:"roominfo"`
					User struct {
						Name   string `json:"name"`
						Avatar string `json:"avatar"`
					} `json:"hostinfo"`
				} `json:"data"`
			}
			req, err := http.NewRequest("GET", fmt.Sprintf("http://www.panda.tv/api_room_v2?roomid=%s", l[1]), nil)
			if err != nil {
				log.Error("xionmao", err)
				break
			}
			req.Header.Set("If-Modified-Since", "0")
			req.Header.Set("Cache-Control", "no-cache")
			resp, err := client.Do(req)
			if err != nil {
				log.Error("xionmao", err)
				break
			}
			r, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(r, &result)
			// 开播
			if result.Data.Info.Status == "2" {
				m := utils.NewMessage()
				m.AddMsg(utils.CQshare(fmt.Sprintf("http://www.panda.tv/%s", result.Data.Info.ID),
					result.Data.User.Name, result.Data.Info.Name, result.Data.User.Avatar))
				group, err := memory.GetLive(roomID[i]).Range()
				if err != nil {
					log.Error("xionmao", err)
				}
				for _, v := range group {
					tmp, _ := strconv.ParseInt(v, 10, 64)
					if registerID[strings.Join([]string{"熊猫", l[1], v}, "-")] != result.Data.Info.Status {
						memory.DefaultMes.Push(
							message.SendMsg(message.MSG_GROUP, tmp, utils.NewMessage().
								AddMsg(utils.CQat("all")).
								AddMsg(utils.CQtext(fmt.Sprintf("[%s]开播了！！！", result.Data.User.Name))).
								Message(), false, ""),
						)
						memory.DefaultMes.Push(
							message.SendMsg(message.MSG_GROUP, tmp, m.Message(), false, ""),
						)
						registerID[strings.Join([]string{"熊猫", l[1], v}, "-")] = result.Data.Info.Status
					}
				}
			}
		case "B站":
			req, err := http.NewRequest("GET", fmt.Sprintf("https://live.bilibili.com/%s", l[1]), nil)
			if err != nil {
				log.Error("bilibili", err)
				break
			}
			req.Header.Set("If-Modified-Since", "0")
			req.Header.Set("Cache-Control", "no-cache")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36")
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
			req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
			resp, err := client.Do(req)
			if err != nil {
				log.Error("bilibili", err)
				break
			}
			// 顺序 1 房间号 2 图片 3 名字 4 状态 5 标题
			var out = make([]string, 6)
			r, _ := ioutil.ReadAll(resp.Body)
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
				group, err := memory.GetLive(roomID[i]).Range()
				if err != nil {
					log.Error("bilibili", err)
				}
				for _, v := range group {
					tmp, _ := strconv.ParseInt(v, 10, 64)
					if registerID[strings.Join([]string{"bilibili", l[1], v}, "-")] != out[4] {
						memory.DefaultMes.Push(
							message.SendMsg(message.MSG_GROUP, tmp, utils.NewMessage().
								AddMsg(utils.CQat("all")).
								AddMsg(utils.CQtext(fmt.Sprintf("[%s]开播了！！！", out[3]))).
								Message(), false, ""),
						)
						memory.DefaultMes.Push(
							message.SendMsg(message.MSG_GROUP, tmp, m.Message(), false, ""),
						)
						registerID[strings.Join([]string{"bilibili", l[1], v}, "-")] = out[4]
					}
				}
			}
		case "虎牙":
			req, err := http.NewRequest("GET", fmt.Sprintf("https://m.huya.com/%s", l[1]), nil)
			if err != nil {
				log.Error("huya", err)
				break
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
				log.Error("huya", err)
				break
			}
			// 顺序 1 图片 2 title 3 状态 4 房间 5 名字
			var out = make([]string, 6)
			r, _ := ioutil.ReadAll(resp.Body)
			for i, match := range re.FindStringSubmatch(string(r)) {
				out[i] = match
			}

			resp.Body.Close()
			// 开播
			if out[3] == "true" {
				m := utils.NewMessage()
				m.AddMsg(utils.CQshare(fmt.Sprintf("http://www.huya.com/%s", out[4]),
					out[5], out[2], out[1]))
				group, err := memory.GetLive(roomID[i]).Range()
				if err != nil {
					log.Error("huya", err)
				}
				for _, v := range group {
					tmp, _ := strconv.ParseInt(v, 10, 64)
					if registerID[strings.Join([]string{"huya", l[1], v}, "-")] != out[3] {
						memory.DefaultMes.Push(
							message.SendMsg(message.MSG_GROUP, tmp, utils.NewMessage().
								AddMsg(utils.CQat("all")).
								AddMsg(utils.CQtext(fmt.Sprintf("[%s]开播了！！！", out[5]))).
								Message(), false, ""),
						)
						memory.DefaultMes.Push(
							message.SendMsg(message.MSG_GROUP, tmp, m.Message(), false, ""),
						)
						registerID[strings.Join([]string{"huya", l[1], v}, "-")] = out[3]
					}
				}
			}
		}
	}
	time.Sleep(10 * time.Second)
}
