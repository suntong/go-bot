package asynchronous

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-bot/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/lexkong/log"
)

var serverState sync.Map

func Jw3Server() {
	// GMT
	l, err := time.LoadLocation("GMT")
	if err != nil {
		log.Error("localtime", err)
		return
	}
	p := struct {
		Ts string `json:"ts"`
	}{Ts: strings.Replace(time.Now().In(l).Format("20060102030405.000"), ".", "", 1)}
	bodyBytes, err := json.Marshal(p)
	if err != nil {
		log.Error("json", err)
		return
	}
	body := bytes.NewBuffer(bodyBytes)
	req, err := http.NewRequest("POST", "https://m.pvp.xoyo.com/tools/zone/check", body)
	if err != nil {
		log.Error("create post", err)
		return
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("clientkey", "1")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", "okhttp/3.11.0")
	result, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("network", err)
		return
	}
	rebytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Error("bytes", err)
		return
	}
	defer result.Body.Close()

	var reJSON struct {
		Data []struct {
			Zone       string `json:"zoneName"`
			ServerName string `json:"serverName"`
			State      bool   `json:"connectState"`
			MainName   string `json:"mainServer"`
		} `json:"data"`
	}
	err = json.Unmarshal(rebytes, &reJSON)
	if err != nil {
		log.Error("json unmarshal", err)
		return
	}
	for i := range reJSON.Data {
		tmp := []string{fmt.Sprint(reJSON.Data[i].State), reJSON.Data[i].Zone, reJSON.Data[i].MainName}
		serverState.Store(reJSON.Data[i].ServerName, tmp)
	}
}

func IsOnline(name string) interface{} {
	if out, ok := serverState.Load(name); ok {
		outTmp, ok := out.([]string)
		if !ok {
			return nil
		}
		if outTmp[0] == "true" {
			// 当前开服
			m := utils.NewMessage()
			m.AddMsg(utils.CQimage("https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1544692640&di=442ab4434e7b703ec890b582c231670f&imgtype=jpg&er=1&src=http%3A%2F%2Fb-ssl.duitang.com%2Fuploads%2Fitem%2F201709%2F14%2F20170914233246_FZJHA.jpeg"))
			m.AddMsg(utils.CQtext(fmt.Sprintf("\n【%s】开服了\n", name)))
			m.AddMsg(utils.CQtext(fmt.Sprintf("%s-%s", outTmp[1], name)))
			return m.Message()
		}
	}
	return nil
}
