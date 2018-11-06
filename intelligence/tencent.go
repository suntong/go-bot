package intelligence

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/teris-io/shortid"
)

const (
	apiUrl = "https://api.ai.qq.com/fcgi-bin/nlp/nlp_textchat"
	ttsUrl = "https://api.ai.qq.com/fcgi-bin/aai/aai_tts"
)

func createTencentChat(text string) string {
	u := url.Values{}
	u.Set("app_id", viper.GetString("tencent.app_id"))
	u.Set("time_stamp", fmt.Sprintf("%d", time.Now().Unix()))
	nonce, _ := shortid.Generate()
	u.Set("nonce_str", nonce)
	u.Set("question", text)
	session, _ := shortid.Generate()
	u.Set("session", session)
	m := md5.New()
	m.Write([]byte(u.Encode() + "&app_key=" + viper.GetString("tencent.app_key")))
	u.Set("sign", strings.ToUpper(hex.EncodeToString(m.Sum(nil))))
	return u.Encode()
}

func createTencentAudio(text string) string {
	u := url.Values{}
	u.Set("app_id", viper.GetString("tencent.app_id"))
	u.Set("time_stamp", fmt.Sprintf("%d", time.Now().Unix()))
	nonce, _ := shortid.Generate()
	u.Set("nonce_str", nonce)
	u.Set("text", text)
	u.Set("speaker", "6")
	u.Set("format", "3")
	u.Set("volume", "0")
	u.Set("speed", "100")
	u.Set("aht", "0")
	u.Set("apc", "58")
	m := md5.New()
	m.Write([]byte(u.Encode() + "&app_key=" + viper.GetString("tencent.app_key")))
	u.Set("sign", strings.ToUpper(hex.EncodeToString(m.Sum(nil))))
	return u.Encode()
}

func GetTencentAudio(text string) []byte {
	req, _ := url.Parse(ttsUrl)
	req.RawQuery = createTencentAudio(text)
	resp, _ := http.DefaultClient.Get(req.String())
	defer resp.Body.Close()
	d, _ := ioutil.ReadAll(resp.Body)
	var result audioResp
	json.Unmarshal(d, &result)
	str := result.Data.Speech.(string)
	str, _ = url.PathUnescape(str)
	data, _ := base64.StdEncoding.DecodeString(str)
	m := md5.New()
	m.Write(data)
	md5string := strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
	if result.Ret == 0 && md5string == result.Data.Md5 {
		return data
	}
	return nil
}

func GetTencentChat(text string) string {
	req, _ := url.Parse(apiUrl)
	req.RawQuery = createTencentChat(text)
	resp, _ := http.DefaultClient.Get(req.String())
	defer resp.Body.Close()
	d, _ := ioutil.ReadAll(resp.Body)
	var result chatResp
	json.Unmarshal(d, &result)
	if result.Ret == 0 {
		return result.Data.Answer
	}
	return ""
}

type chatResp struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data struct {
		Session string `json:"session"`
		Answer  string `json:"answer"`
	} `json:"data"`
}

type audioResp struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data struct {
		Format int         `json:"format"`
		Speech interface{} `json:"speech"`
		Md5    string      `json:"md5sum"`
	}
}
