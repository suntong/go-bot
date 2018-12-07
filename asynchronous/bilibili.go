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

func bilibiliOnline(addr string) interface{} {
	var result struct {
		Data struct {
			Roomid     int    `json:"room_id"`
			RoomName   string `json:"title"`
			RoomStatus string `json:"status"`
			Avatar     string `json:"cover"`
			Owner      string `json:"uname"`
		} `json:"data"`
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://api.live.bilibili.com/AppRoom/index?room_id=%s&platform=android", addr), nil)
	if err != nil {
		log.Error("B站", err)
		return err
	}
	req.Header.Set("If-Modified-Since", "0")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)
	if err != nil {
		log.Error("B站", err)
		return err
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("B站", err)
		return err
	}
	resp.Body.Close()
	json.Unmarshal(r, &result)
	// 开播
	if result.Data.RoomStatus == "LIVE" {
		m := utils.NewMessage()
		m.AddMsg(utils.CQimage(result.Data.Avatar))
		m.AddMsg(utils.CQtext(fmt.Sprintf("\n【%s】开播了", result.Data.Owner)))
		m.AddMsg(utils.CQtext(fmt.Sprintf("%s", result.Data.RoomName)))
		m.AddMsg(utils.CQtext(fmt.Sprintf("\n直播间地址：【%s】", fmt.Sprintf("https://live.bilibili.com/%d", result.Data.Roomid))))

		return m.Message()
	} else if result.Data.RoomStatus != "" {
		return nil
	}
	return errors.New("network error")
}
