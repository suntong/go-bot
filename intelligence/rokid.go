package intelligence

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/lexkong/log"
)

// version={version};time={time};sign={sign};key={key};device_type_id={device_type_id};device_id={device_id};service={service}
// key={key}&device_type_id={device_type_id}&device_id={device_id}&service={service}&version={version}&time={time}&secret={secret}

func generateAuthorization(version, secret, key, deviceTypeID, deviceID, service string) string {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	return fmt.Sprintf("version=%s;time=%s;sign=%s;key=%s;device_type_id=%s;device_id=%s;service=%s",
		version, now, generateSign(now, version, secret, key, deviceTypeID, deviceID, service),
		key, deviceTypeID, deviceID, service)
}

func generateSign(now, version, secret, key, deviceTypeID, deviceID, service string) string {
	src := "key=" + key + "&device_type_id=" + deviceTypeID + "&device_id=" + deviceID + "&service=" + service + "&version=" + version + "&time=" + now + "&secret=" + secret
	sign := strings.ToUpper(makeMD5(src))
	return sign
}

func makeMD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

const (
	ttsRokidApi = "https://apigwrest.open.rokid.com/api/v1/tts/TtsProxy/Tts"
)

func GetRokidAudio(text string) []byte {
	client := &http.Client{}
	out := generateAuthorization(viper.GetString("rokid.version"),
		viper.GetString("rokid.secret"),
		viper.GetString("rokid.key"),
		viper.GetString("rokid.device_type_id"),
		viper.GetString("rokid.device_id"),
		viper.GetString("name"))
	j := struct {
		Text      string `json:"text"`
		Declaimer string `json:"declaimer"`
		Codec     string `json:"codec"`
	}{Text: text, Declaimer: "xmly", Codec: "mp3"}
	jsonBytes, _ := json.Marshal(&j)
	req, err := http.NewRequest("POST", ttsRokidApi, bytes.NewReader(jsonBytes))
	if err != nil {
		log.Error("rokid", err)
		return []byte{}
	}
	req.Header.Set("Authorization", out)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error("rokid", err)
		return []byte{}
	}

	var result struct {
		Data string `json:"voice"`
	}

	voiceBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("rokid", err)
		return []byte{}
	}
	json.Unmarshal(voiceBytes, &result)
	str := result.Data
	str, _ = url.PathUnescape(str)
	data, _ := base64.StdEncoding.DecodeString(str)
	return data
}
