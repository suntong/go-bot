package handle

import (
	"fmt"
	"go-bot/intelligence"
	"go-bot/pkg/memory"
	"go-bot/pkg/message"
	"go-bot/utils"
	"strings"
)

// 通用

func Command(s message.EventJSON) interface{} {
	if !utils.AtSelf(s.RawMsg, s.Self) {
		return s
	}

	raw := strings.Trim(utils.Fransferred(s.RawMsg), " ")
	l := strings.Split(raw, " ")
	if len(l) > 1 {
		switch l[0] {
		case "语音":
			go func(goh message.EventJSON, text string) {
				m := utils.NewMessage()
				m.AddMsg(utils.CQat(fmt.Sprint(goh.UserID)))
				m.AddMsg(utils.CQBase64record(intelligence.GetBaiduAudio(text), false))
				memory.DefaultMes.Push(
					message.SendMsg(goh.MsgType, goh.GroupID,
						m.Message(), false, ""),
				)
			}(s, strings.Join(l[1:], " "))
			return nil
		case "私聊":
			go func(goh message.EventJSON, text string) {
				m := utils.NewMessage()
				m.AddMsg(utils.CQtext(text))
				memory.DefaultMes.Push(
					message.SendMsg(message.MSG_PRIVATE, goh.UserID,
						m.Message(), false, goh.MsgType),
				)
			}(s, strings.Join(l[1:], " "))
			return nil
		default:
			return s
		}
	}

	return s
}

func handleCmd(cmd []string) (int, string) {
	return 0, ""
}
