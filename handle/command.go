package handle

import (
	"fmt"
	"go-bot/intelligence"
	"go-bot/pkg/memory"
	"go-bot/pkg/message"
	"go-bot/utils"
	"strconv"
	"strings"

	"github.com/lexkong/log"
)

// 通用

func Command(s message.EventJSON) interface{} {
	if !utils.AtSelf(s.RawMsg, s.Self) {
		return s
	}

	raw := strings.Trim(utils.Fransferred(s.RawMsg), " ")
	l := strings.Split(raw, " ")

	if len(l) == 3 && l[0] == "监控直播" && (l[1] == "斗鱼" || l[1] == "熊猫" || l[1] == "B站" || l[1] == "虎牙") {
		go func(goh message.EventJSON, text string, room string) {
			text = strings.Trim(text, " ")
			tmp, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				log.Error("监控", err)
				return
			}
			m := utils.NewMessage()
			m.AddMsg(utils.CQtext(
				fmt.Sprintf("监控[%s]频道", fmt.Sprintf("%d", tmp)),
			))
			memory.DefaultMes.Push(
				message.SendMsg(goh.MsgType, goh.GroupID,
					m.Message(), false, ""),
			)
			memory.GetLive(strings.Join([]string{room, fmt.Sprintf("%d", tmp)}, "-")).Push(goh.GroupID)
			memory.GetLive("liveRoom").Push(strings.Join([]string{room, fmt.Sprintf("%d", tmp)}, "-"))
		}(s, l[2], l[1])
		return nil
	}

	if len(l) > 1 {
		switch l[0] {
		case "语音":
			go func(goh message.EventJSON, text string) {
				m := utils.NewMessage()
				m.AddMsg(utils.CQat(fmt.Sprint(goh.UserID)))
				m.AddMsg(utils.CQBase64record(intelligence.GetRokidAudio(text), false))
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

// at 状态
func handleCmd(cmd []string) (int, string) {
	// 拆解cmd
	if len(cmd) == 0 {
		return -1, ""
	}
	switch len(cmd) {
	case 1:
	case 2:
	case 3:
	}
	// 不存在 原路返回
	return -1, strings.Join(cmd, " ")
}
