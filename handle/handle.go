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
	if r := load(e, Command); r != nil {
		h := r.(message.EventJSON)
		if utils.AtSelf(h.RawMsg, h.Self) && h.MsgType == message.MSG_GROUP {
			go func(goh message.EventJSON) {
				m := utils.NewMessage()
				m.AddMsg(utils.CQat(fmt.Sprint(goh.UserID)))
				m.AddMsg(utils.CQtext(intelligence.GetTencentChat(
					utils.Fransferred(goh.RawMsg),
				)))
				memory.DefaultMes.Push(
					message.SendMsg(goh.MsgType, goh.GroupID,
						m.Message(), false, ""),
				)
			}(h)
		}
	}
	fmt.Println(string(data))
	return nil
}

func Send() interface{} {
	r, err := memory.DefaultMes.Pop()
	if err != nil {
		log.Error("send", err)
	}
	return message.String2Interface(r)
}
