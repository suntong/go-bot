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

func xionmaoOnline(addr string) interface{} {
	var result struct {
		Data struct {
			Info struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Status   string `json:"status"`
				Pictures struct {
					Img string `json:"cover_img"`
				} `json:"pictures"`
			} `json:"roominfo"`
			User struct {
				Name string `json:"name"`
			} `json:"hostinfo"`
		} `json:"data"`
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://www.panda.tv/api_room_v2?roomid=%s", addr), nil)
	if err != nil {
		log.Error("熊猫", err)
		return err
	}
	req.Header.Set("If-Modified-Since", "0")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)
	if err != nil {
		log.Error("熊猫", err)
		return err
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("熊猫", err)
		return err
	}
	resp.Body.Close()
	json.Unmarshal(r, &result)
	// 开播
	if result.Data.Info.Status == "2" {
		m := utils.NewMessage()
		// m.AddMsg(utils.CQimage(result.Data.Info.Pictures.Img))
		m.AddMsg(utils.CQtext(fmt.Sprintf("【%s】开播了\n", result.Data.User.Name)))
		m.AddMsg(utils.CQtext(fmt.Sprintf("%s", result.Data.Info.Name)))
		m.AddMsg(utils.CQtext(fmt.Sprintf("http://www.panda.tv/%s", result.Data.Info.ID)))
		return m.Message()
	} else if result.Data.Info.Status != "" {
		return nil
	}
	return errors.New("network error")
}
