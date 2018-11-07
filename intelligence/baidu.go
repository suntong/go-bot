package intelligence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lexkong/log"

	"github.com/spf13/viper"
)

const (
	baiduUrl = "http://tsn.baidu.com/text2audio"
)

func getBaiduToken() string {
	u := url.Values{}
	u.Set("grant_type", "client_credentials")
	u.Set("client_id", viper.GetString("baidu.client_id"))
	u.Set("client_secret", viper.GetString("baidu.client_secret"))
	uu, _ := url.Parse("https://openapi.baidu.com/oauth/2.0/token")
	uu.RawQuery = u.Encode()
	fmt.Println(uu.String())
	resp, err := http.DefaultClient.Get(uu.String())
	if err != nil {
		log.Error("baiduToken", err)
		return ""
	}
	defer resp.Body.Close()
	fmt.Println(resp)
	var j struct {
		Token string `json:"access_token"`
	}
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
	json.Unmarshal(data, &j)
	return j.Token
}

func createBaiduAudio(text string) string {
	u := url.Values{}
	u.Set("tex", text)

	u.Set("tok", getBaiduToken())
	u.Set("cuid", viper.GetString("baidu.cuid"))
	u.Set("ctp", "1")
	u.Set("lan", "zh")
	u.Add("per", "0")
	return u.Encode()
}

func GetBaiduAudio(text string) []byte {
	req, _ := url.Parse(baiduUrl)
	req.RawQuery = createBaiduAudio(text)
	resp, err := http.DefaultClient.Get(req.String())
	if err != nil {
		log.Error("BaiduAudio", err)
		return []byte{}
	}
	defer resp.Body.Close()
	d, _ := ioutil.ReadAll(resp.Body)
	return d
}
