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

	"github.com/spf13/viper"

	"github.com/lexkong/log"
)

func PrivateCmd(s message.EventJSON) interface{} {
	if s.MsgType != message.MSG_PRIVATE {
		return s
	}
	raw := strings.Trim(utils.Fransferred(s.RawMsg), " ")
	l := strings.Split(raw, " ")
	code, result := handleCmd(l)
	if code == -1 {
		return s
	}
	switch code {
	case 8:
		go func(j message.EventJSON) {
			m := utils.NewMessage()
			result, err := memory.GetKV(fmt.Sprintf("%s-%d", "qq", j.UserID)).GetKey()
			if err != nil {
				log.Error("监控列表", err)
				return
			}
			if len(result) == 0 {
				m.AddMsg(utils.CQtext("空"))
			} else {
				m.AddMsg(utils.CQtext("监控列表"))
			}
			for i := range result {
				m.AddMsg(utils.CQtext(fmt.Sprintf("\n%s", result[i])))
			}
			memory.DefaultMes.Push(
				message.SendMsg(message.MSG_PRIVATE, j.UserID,
					m.Message(), false, message.MSG_GROUP),
			)
		}(s)
		return nil
	case 10:
		go func(j message.EventJSON, key string) {
			m := utils.NewMessage()
			result, err := memory.GetKV(fmt.Sprintf("%s-%d", "qq", j.UserID)).Del(key)
			if err != nil {
				log.Error("删除监控", err)
				return
			}
			if result > 0 {
				m.AddMsg(utils.CQtext(fmt.Sprintf("删除[%s]成功!", key)))
			} else {
				m.AddMsg(utils.CQtext(fmt.Sprintf("删除[%s]失败，可能不存在。", key)))
			}
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.UserID,
					m.Message(), false, message.MSG_GROUP),
			)
		}(s, result[0])
		return nil
	case 11:
		go func(j message.EventJSON, id string, room string) {
			m := utils.NewMessage()
			m.AddMsg(utils.CQtext(
				fmt.Sprintf("监控[%s]频道", id),
			))
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.UserID,
					m.Message(), false, message.MSG_GROUP),
			)
			memory.GetLive("inform").Push(fmt.Sprintf("%s-%d", "qq", j.UserID))
			memory.GetKV(fmt.Sprintf("%s-%d", "qq", j.UserID)).Set(strings.Join([]string{room, id}, "-"), "false")
			memory.GetLive("liveRoom").Push(strings.Join([]string{room, id}, "-"))
		}(s, result[1], result[0])
		return nil
	}
	return s
}

// 通用

func Command(s message.EventJSON) interface{} {
	if !utils.AtSelf(s.RawMsg, s.Self) {
		return s
	}

	if s.MsgType != message.MSG_GROUP {
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
		go func(j message.EventJSON, id string, room string) {
			m := utils.NewMessage()
			m.AddMsg(utils.CQat(fmt.Sprint(j.UserID)))
			m.AddMsg(utils.CQtext(
				fmt.Sprintf("监控[%s]频道", id),
			))
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
			memory.GetLive("inform").Push(fmt.Sprintf("%s-%d", "group", j.GroupID))
			memory.GetKV(fmt.Sprintf("%s-%d", "group", j.GroupID)).Set(strings.Join([]string{room, id}, "-"), "false")
			memory.GetLive("liveRoom").Push(strings.Join([]string{room, id}, "-"))
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
			t := rand.Intn(viper.GetInt("ban_time"))
			m := utils.NewMessage()
			m.AddMsg(utils.CQat(fmt.Sprint(j.UserID)))
			// m.AddMsg(utils.CQimage("https://ws1.sinaimg.cn/large/54d358dbly1fvbwx2kzc7g20e80e8wua.gif"))
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
	case 5:
		go func(j message.EventJSON) {
			t := rand.Intn(100)
			m := utils.NewMessage()
			m.AddMsg(utils.CQat(fmt.Sprint(j.UserID)))
			m.AddMsg(utils.CQtext(fmt.Sprintf("roll中了[%d]点！！", t)))
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
		}(s)
		return nil
	case 6:
		go func(j message.EventJSON) {
			memory.DefaultMes.Push(
				message.GetGroupMemberList(j.GroupID),
			)
		}(s)
		return nil
	case 7:
		go func(j message.EventJSON) {
			m := utils.NewMessage()
			result, err := memory.GetKV(fmt.Sprintf("%s-%d", "group", j.GroupID)).GetKey()
			if err != nil {
				log.Error("监控列表", err)
				return
			}
			m.AddMsg(utils.CQat(fmt.Sprintf("%d", j.UserID)))
			if len(result) == 0 {
				m.AddMsg(utils.CQtext("空"))
			}
			for i := range result {
				m.AddMsg(utils.CQtext(fmt.Sprintf("\n%s", result[i])))
			}
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
		}(s)
		return nil
	case 9:
		go func(j message.EventJSON, key string) {
			m := utils.NewMessage()
			result, err := memory.GetKV(fmt.Sprintf("%s-%d", "group", j.GroupID)).Del(key)
			if err != nil {
				log.Error("删除监控", err)
				return
			}
			m.AddMsg(utils.CQat(fmt.Sprintf("%d", j.UserID)))
			if result > 0 {
				m.AddMsg(utils.CQtext(fmt.Sprintf("\n删除[%s]成功!", key)))
			} else {
				m.AddMsg(utils.CQtext(fmt.Sprintf("\n删除[%s]失败，可能不存在。", key)))
			}
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
		}(s, result[0])
		return nil
	case 12:
		go func(j message.EventJSON) {
			m := utils.NewMessage()
			result, err := memory.GetLive(fmt.Sprintf("%s-%d", "draw", j.GroupID)).Push(j.UserID)
			if err != nil {
				log.Error("抽奖", err)
				return
			}
			m.AddMsg(utils.CQat(fmt.Sprintf("%d", j.UserID)))
			if result > 0 {
				m.AddMsg(utils.CQtext("报名成功"))
			}
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
		}(s)
		return nil
	case 13:
		go func(j message.EventJSON) {
			m := utils.NewMessage()
			result, err := memory.GetLive(fmt.Sprintf("%s-%d", "draw", j.GroupID)).Range()
			if err != nil {
				log.Error("抽奖", err)
				return
			}
			if len(result) == 0 {
				m.AddMsg("报名人数为空")
			} else {
				n := rand.Intn(len(result))
				if err != nil {
					log.Error("抽奖", err)
					return
				}
				m.AddMsg(utils.CQtext("恭喜")).AddMsg(utils.CQat(result[n])).AddMsg(utils.CQtext("获奖"))
			}
			memory.GetLive(fmt.Sprintf("%s-%d", "draw", j.GroupID)).Close()
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
		}(s)
		return nil
	case 14:
		go func(j message.EventJSON, ns string) {
			tmp, err := strconv.ParseInt(ns, 10, 64)
			if err != nil {
				log.Error("抽奖", err)
			}
			m := utils.NewMessage()
			result, err := memory.GetLive(fmt.Sprintf("%s-%d", "draw", j.GroupID)).Range()
			if err != nil {
				log.Error("抽奖", err)
				return
			}
			if len(result) == 0 {
				m.AddMsg(utils.CQtext("抽奖池为空"))
			} else {
				m.AddMsg(utils.CQtext("抽奖列表"))
			}
			for len(result) > 0 && tmp > 0 {
				n := rand.Intn(len(result))
				if err != nil {
					log.Error("抽奖", err)
					return
				}
				m.AddMsg(utils.CQtext("\n恭喜")).AddMsg(utils.CQat(result[n])).AddMsg(utils.CQtext("获奖"))
				if n < len(result)-1 {
					result = append(result[:n], result[n+1:]...)
				} else {
					result = result[:n]
				}
				tmp--
			}
			memory.GetLive(fmt.Sprintf("%s-%d", "draw", j.GroupID)).Close()
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
			)
		}(s, result[0])
		return nil
	case 15:
		go func(j message.EventJSON) {
			m := utils.NewMessage()
			result, err := memory.GetLive(fmt.Sprintf("%s-%d", "draw", j.GroupID)).Range()
			if err != nil {
				log.Error("抽奖列表", err)
				return
			}
			if len(result) == 0 {
				m.AddMsg(utils.CQtext("空"))
			} else {
				m.AddMsg(utils.CQtext("抽奖列表"))
			}
			for i := range result {
				m.AddMsg(utils.CQtext("\n")).AddMsg(utils.CQat(result[i]))
			}
			memory.DefaultMes.Push(
				message.SendMsg(j.MsgType, j.GroupID,
					m.Message(), false, ""),
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
		case "roll":
			return 5, []string{}
		case "点赞":
			return 1, []string{}
		case "我要自闭":
			fallthrough
		case "我禁我自己":
			fallthrough
		case "禁言抽奖":
			return 4, []string{}
		case "群信息":
			return 6, []string{}
		case "监控列表":
			return 7, []string{}
		case "私聊监控列表":
			return 8, []string{}
		case "抽奖报名":
			return 12, []string{}
		case "单人抽奖":
			return 13, []string{}
		case "抽奖列表":
			return 15, []string{}
		}
	case 2:
		switch c {
		case "删除监控":
			return 9, []string{cmd[1]}
		case "删除私聊监控":
			return 10, []string{cmd[1]}
		case "多人抽奖":
			tmp, err := strconv.ParseInt(cmd[1], 10, 64)
			if err != nil {
				log.Error("多人抽奖", err)
				return -1, nil
			}
			return 14, []string{fmt.Sprint(tmp)}
		}
	case 3:
		switch c {
		case "私聊监控":
			if supportLive(cmd[1]) {
				tmp, err := strconv.ParseInt(cmd[2], 10, 64)
				if err != nil {
					log.Error("监控", err)
					return -1, nil
				}
				return 11, []string{strings.ToUpper(cmd[1]), fmt.Sprintf("%d", tmp)}
			}
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
