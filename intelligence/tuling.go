package intelligence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-bot/utils"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"

	"github.com/lexkong/log"
)

func getLocation() interface{} {
	return nil
}

func getUserInfo(userId, groupId, userIdName string) interface{} {
	return struct {
		ApiKey     string `json:"apiKey"`
		UserId     string `json:"userId"`
		GroupId    string `json:"groupId"`
		UserIdName string `json:"userIdName"`
	}{ApiKey: viper.GetString("tuling.api_key"), UserId: userId, GroupId: groupId, UserIdName: userIdName}
}

func GetTulingChat(reqType int, text, userId, groupId, userIdName string, at bool) interface{} {
	var perception interface{}
	switch reqType {
	case 0:
		perception = struct {
			Data struct {
				Text string `json:"text"`
			} `json:"inputText"`
		}{Data: struct {
			Text string `json:"text"`
		}{Text: text}}
	case 1:
		perception = struct {
			Data struct {
				Text string `json:"url"`
			} `json:"inputImage"`
		}{Data: struct {
			Text string `json:"url"`
		}{Text: text}}
	case 2:
		perception = struct {
			Data struct {
				Text string `json:"url"`
			} `json:"inputMedia"`
		}{Data: struct {
			Text string `json:"url"`
		}{Text: text}}
	}

	result := struct {
		ReqType    int         `json:"reqType"`
		Perception interface{} `json:"perception"`
		UserInfo   interface{} `json:"userInfo"`
	}{ReqType: reqType, Perception: perception, UserInfo: getUserInfo(userId, groupId, userIdName)}
	jsonBytes, _ := json.Marshal(&result)
	fmt.Println(string(jsonBytes))
	resp, err := http.Post("http://openapi.tuling123.com/openapi/api/v2", "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		log.Error("tuling", err)
		return ""
	}

	var out struct {
		Intent struct {
			Code       int    `json:"code"`
			IntentName string `json:"intentName"`
			ActionName string `json:"actionName"`
			Parameters struct {
				NearbyPlace string `json:"nearby_place"`
			} `json:"parameters"`
		} `json:"intent"`
		Results []struct {
			GroupType  int         `json:"groupType"`
			ResultType string      `json:"resultType"`
			Values     interface{} `json:"values"`
		} `json:"results"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("tuling", err)
		return nil
	}

	json.Unmarshal(body, &out)

	m := utils.NewMessage()
	if at {
		m.AddMsg(utils.CQat(userId))
	}
	for _, v := range out.Results {
		valus, ok := v.Values.(map[string]interface{})
		if !ok {
			return nil
		}
		switch v.ResultType {
		case "text":
			m.AddMsg(utils.CQtext(valus["text"].(string)))
		case "url":
			m.AddMsg(utils.CQtext(valus["url"].(string)))
		case "voice":
			m.AddMsg(utils.CQrecord(valus["voice"].(string), false))
		case "image":
			m.AddMsg(utils.CQimage(valus["image"].(string)))
		default:
		}
	}
	return m.Message()
}
