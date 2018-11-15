package handle

import (
	"encoding/json"
	"fmt"
	"go-bot/intelligence"
	"go-bot/pkg/memory"
	"go-bot/pkg/message"
	"go-bot/utils"

	"github.com/lexkong/log"
)

func load(e message.EventJSON, mw ...func(message.EventJSON) interface{}) interface{} {
	var r interface{} = e
	for _, v := range mw {
		if r = v(e); r == nil {
			return nil
		}
	}
	return r
}

func Handle(data []byte) error {
	// 返回错误关闭
	var e message.EventJSON
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}
	fmt.Println(string(data))
	if r := load(e,
		Command); r != nil {
		h := r.(message.EventJSON)
		if utils.AtSelf(h.RawMsg, h.Self) && h.MsgType == message.MSG_GROUP {
			go func(goh message.EventJSON) {
				memory.DefaultMes.Push(
					message.SendMsg(goh.MsgType, goh.GroupID,
						intelligence.GetTulingChat(
							0, utils.Fransferred(goh.RawMsg), fmt.Sprint(h.UserID), fmt.Sprint(h.GroupID), h.Sender.Nick, true), false, ""),
				)
			}(h)
		}
	}
	return nil
}

func Send() interface{} {
	r, err := memory.DefaultMes.Pop()
	if err != nil {
		log.Error("send", err)
	}
	return message.String2Interface(r)
}
