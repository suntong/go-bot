package handle

import (
	"fmt"
	"go-bot/intelligence"
	"go-bot/pkg/memory"
	"go-bot/pkg/message"
	"go-bot/utils"
	"math/rand"
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
	code, result := handleCmd(l)
	if code == -1 {
		return s
	}

	switch code {
	case 1:
		go func(j message.EventJSON) {
			m := utils.NewMessage()
			m.AddMsg(utils.CQat(fmt.Sprint(j.UserID)))
			m.AddMsg(utils.CQtext("点赞成功"))
			memory.DefaultMes.Push(
				message.SendLike(j.UserID, 10),
			)
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
		}(s)
		return nil
	case 2:
		go func(j message.EventJSON, text string, room string) {
			text = strings.Trim(text, " ")
			tmp, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				log.Error("监控", err)
				return
			}
			m := utils.NewMessage()
			m.AddMsg(utils.CQat(fmt.Sprint(j.UserID)))
			m.AddMsg(utils.CQtext(
				fmt.Sprintf("监控[%s]频道", fmt.Sprintf("%d", tmp)),
			))
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
			memory.GetLive(strings.Join([]string{room, fmt.Sprintf("%d", tmp)}, "-")).Push(j.GroupID)
			memory.GetLive("liveRoom").Push(strings.Join([]string{room, fmt.Sprintf("%d", tmp)}, "-"))
		}(s, result[1], result[0])
		return nil
	case 3:
		go func(j message.EventJSON, text string) {
			m := utils.NewMessage()
			m.AddMsg(utils.CQat(fmt.Sprint(j.UserID)))
			m.AddMsg(utils.CQBase64record(intelligence.GetRokidAudio(text), false))
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
		}(s, strings.Join(result, " "))
		return nil
	case 4:
		go func(j message.EventJSON) {
			t := rand.Intn(300)
			m := utils.NewMessage()
			m.AddMsg(utils.CQat(fmt.Sprint(j.UserID)))
			m.AddMsg(utils.CQtext(fmt.Sprintf("恭喜您抽中了%d秒！！", t)))
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
			memory.DefaultMes.Push(
				message.SetGroupBan(j.GroupID, j.UserID, int32(t)),
			)
		}(s)
		return nil
	}
	return s
}

func handleCmd(cmd []string) (int, []string) {
	// 拆解cmd
	if len(cmd) == 0 {
		return -1, nil
	}

	c := cmd[0]
	if len(cmd) > 1 {
		switch c {
		case "语音":
			return 3, cmd[1:]
		}
	}

	switch len(cmd) {
	case 1:
		switch c {
		case "点赞":
			return 1, []string{}
		case "禁言抽奖":
			return 4, []string{}
		}
	case 2:
	case 3:
		switch c {
		case "监控":
			if supportLive(cmd[1]) {
				tmp, err := strconv.ParseInt(cmd[2], 10, 64)
				if err != nil {
					log.Error("监控", err)
					return -1, nil
				}
				return 2, []string{strings.ToUpper(cmd[1]), fmt.Sprintf("%d", tmp)}
			}
		}
	}
	// 不存在 原路返回
	return -1, cmd
}

func supportLive(live string) bool {
	live = strings.ToUpper(live)
	switch live {
	case "斗鱼":
		fallthrough
	case "熊猫":
		fallthrough
	case "B站":
		fallthrough
	case "虎牙":
		return true
	default:
		return false
	}
}
