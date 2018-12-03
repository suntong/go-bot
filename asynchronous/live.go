package asynchronous

import (
	"go-bot/pkg/memory"
	"go-bot/pkg/message"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lexkong/log"
)

var (
	registerID = sync.Map{}
)

func notification() {
	sometion, _ := memory.GetLive("inform").Range()
	for i := range sometion {
		item := sometion[i]
		db := memory.GetKV(item)
		keys, _ := db.GetKey()
		for i := range keys {
			result, _ := db.Get(keys[i])
			out, ok := registerID.Load(keys[i])
			if !ok {
				continue
			}

			if out == nil {
				db.Set(keys[i], "false")
				continue
			}

			if out != nil && result == "false" {
				tmp := strings.Index(item, "-")
				id, err := strconv.ParseInt(item[tmp+1:], 10, 64)
				if err != nil {
					log.Error("监控", err)
					continue
				}
				if item[:tmp] == "group" {
					memory.DefaultMes.Push(
						message.SendMsg(message.MSG_GROUP, id, out, false, ""),
					)
				} else {
					memory.DefaultMes.Push(
						message.SendMsg(message.MSG_PRIVATE, id, out, false, message.MSG_GROUP),
					)
				}
				db.Set(keys[i], "true")
			}

		}

	}
}

func cyclic() {
	roomID, err := memory.GetLive("liveRoom").Range()
	if err != nil {
		log.Fatal("cyclic", err)
	}
	for i := range roomID {
		item := strings.Split(roomID[i], "-")
		var out interface{}
		switch item[0] {
		case "斗鱼":
			out = douyuOnline(item[1])
		case "熊猫":
			out = xionmaoOnline(item[1])
		case "虎牙":
			out = huyaOnline(item[1])
		case "B站":
			out = bilibiliOnline(item[1])
		}
		if err, ok := out.(error); ok {
			log.Errorf(err, "cyclic:%s[%s]", item[0], item[1])
			continue
		}
		if out != nil {
			registerID.Store(roomID[i], out)
			notification()
		} else {
			registerID.Store(roomID[i], nil)
		}
	}
}

func Run() {
	go func() {
		for range time.NewTicker(30 * time.Second).C {
			cyclic()
		}
	}()
}
