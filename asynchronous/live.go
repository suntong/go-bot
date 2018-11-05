package asynchronous

import (
	"encoding/json"
	"fmt"
	"go-bot/messages"
	"go-bot/pkg/message"
	"go-bot/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	roomID = []string{"斗鱼-533640"}
)

var (
	registerID = make(map[string]string)
)

var (
	group = []int64{372197768, 913608887}
)

func Douyu() {
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

			resp, _ := http.DefaultClient.Get(fmt.Sprintf("http://open.douyucdn.cn/api/RoomApi/room/%s", l[1]))
			r, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(r, &result)
			// 开播
			if _, ok := registerID["斗鱼-"+l[1]]; !ok {
				registerID["斗鱼-"+l[1]] = "0"
			}
			if registerID["斗鱼-"+l[1]] != result.Data.RoomStatus {
				if result.Data.RoomStatus == "1" {
					m := utils.NewMessage()
					m.AddMsg(utils.CQshare(fmt.Sprintf("https://www.douyu.com/%s", result.Data.Roomid),
						result.Data.Owner, result.Data.RoomName, result.Data.Avatar))
					for g := range group {
						messages.GetDefaultMessages().Push(
							message.SendMsg(message.MSG_GROUP, group[g], m.Message(), false, ""),
						)
					}
				}
				registerID["斗鱼-"+l[1]] = result.Data.RoomStatus
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
			resp, _ := http.DefaultClient.Get(fmt.Sprintf("http://www.panda.tv/api_room_v2?roomid=%s", l[1]))
			r, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(r, &result)
			// 开播
			if _, ok := registerID["熊猫-"+l[1]]; !ok {
				registerID["熊猫-"+l[1]] = "0"
			}
			if registerID["熊猫-"+l[1]] != result.Data.Info.Status {
				if result.Data.Info.Status == "2" {
					m := utils.NewMessage()
					m.AddMsg(utils.CQshare(fmt.Sprintf("http://www.panda.tv/%s", result.Data.Info.ID),
						result.Data.User.Name, result.Data.Info.Name, result.Data.User.Avatar))
					for g := range group {
						messages.GetDefaultMessages().Push(
							message.SendMsg(message.MSG_GROUP, group[g], m.Message(), false, ""),
						)
					}
				}
				registerID["熊猫-"+l[1]] = result.Data.Info.Status
			}
		}
	}
	time.Sleep(30 * time.Second)
}
