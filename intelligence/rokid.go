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

func GetRokidAudio(text string) []byte {
	client := &http.Client{}

	j := struct {
		Text      string `json:"text"`
		Declaimer string `json:"declaimer"`
		Codec     string `json:"codec"`
	}{Text: text, Declaimer: "xmly", Codec: "mp3"}
	jsonBytes, _ := json.Marshal(&j)
	req, err := http.NewRequest("POST", "https://apigwrest.open.rokid.com/api/v1/tts/TtsProxy/Tts", bytes.NewReader(jsonBytes))
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

	voiceBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(voiceBytes, &result)
	str := result.Data
	str, _ = url.PathUnescape(str)
	data, _ := base64.StdEncoding.DecodeString(str)
	return data
}
