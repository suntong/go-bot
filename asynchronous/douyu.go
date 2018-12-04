package asynchronous

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-bot/utils"
	"io/ioutil"
	"net/http"

	"github.com/lexkong/log"
)

func douyuOnline(addr string) interface{} {
	var result struct {
		Data struct {
			Roomid     string `json:"room_id"`
			RoomName   string `json:"room_name"`
			RoomStatus string `json:"room_status"`
			Thumb      string `json:"room_thumb"`
			Owner      string `json:"owner_name"`
		} `json:"data"`
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://open.douyucdn.cn/api/RoomApi/room/%s", addr), nil)
	if err != nil {
		log.Error("douyu", err)
		return err
	}
	req.Header.Set("If-Modified-Since", "0")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)
	if err != nil {
		log.Error("douyu", err)
		return err
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("斗鱼", err)
		return err
	}
	resp.Body.Close()
	json.Unmarshal(r, &result)
	// 开播
	if result.Data.RoomStatus == "1" {
		m := utils.NewMessage()
		m.AddMsg(utils.CQshare(fmt.Sprintf("https://www.douyu.com/%s", result.Data.Roomid),
			result.Data.Owner, result.Data.RoomName, result.Data.Thumb))
		return m.Message()
	} else if result.Data.RoomStatus != "" {
		return nil
	}
	return errors.New("network error")
}
