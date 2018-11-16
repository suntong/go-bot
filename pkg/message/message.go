package message

import (
	"encoding/json"
)

type sendJSON struct {
	Action string     `json:"action"`
	Params sendParams `json:"params"`
}

type sendParams struct {
	MsgType    string      `json:"message_type"`
	DiscussID  int64       `json:"discuss_id"`
	GroupID    int64       `json:"group_id"`
	UserID     int64       `json:"user_id"`
	Msg        interface{} `json:"message"`
	MsgID      int32       `json:"message_id"`
	AutoEscape bool        `json:"auto_escape"`
	Times      int32       `json:"times"`
	Reject     bool        `json:"reject_add_request"`
	Duration   int32       `json:"duration"`
	Flag       string      `json:"flag"`
	Enable     bool        `json:"enable"`
	Card       string      `json:"card"`
	Dismiss    bool        `json:"is_dismiss"`
	Special    string      `json:"special_title"`
	Approve    bool        `json:"approve"`
	Remark     string      `json:"remark"`
	SubType    string      `json:"sub_type"`
	Reason     string      `json:"reason"`
	NoCache    bool        `json:"no_cache"`
	GroupName  string      `json:"group_name"`
	Anonymous  struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Flag string `json:"flag"`
	} `json:"anonymous"`
}

const (
	MSG_PRIVATE = "private"
	MSG_GROUP   = "group"
	MSG_DISCUSS = "discuss"
)

func String2Interface(message string) interface{} {
	var result sendJSON
	json.Unmarshal([]byte(message), &result)
	return result
}

func SendMsg(msgType string, id int64, msg interface{}, auto_escape bool, sub string) string {
	m := sendJSON{
		Action: "send_msg",
		Params: sendParams{
			MsgType:    msgType,
			UserID:     id,
			GroupID:    id,
			DiscussID:  id,
			Msg:        msg,
			AutoEscape: auto_escape,
		},
	}

	if len(sub) > 0 {
		m.Params.SubType = sub
	}
	result, _ := json.Marshal(m)
	return string(result)
}

func DeleteMsg(msgID int32) string {
	m := sendJSON{
		Action: "delete_msg",
		Params: sendParams{
			MsgID: msgID,
		},
	}
	result, _ := json.Marshal(m)
	return string(result)
}

func SendLike(userID int64, times int32) string {
	m := sendJSON{
		Action: "send_like",
		Params: sendParams{
			UserID: userID,
			Times:  times,
		},
	}
	result, _ := json.Marshal(m)
	return string(result)
}

func SetGroupKick(groupID int64, userID int64, rejcet bool) string {

	m := sendJSON{
		Action: "set_group_kick",
		Params: sendParams{
			GroupID: groupID,
			UserID:  userID,
			Reject:  rejcet,
		},
	}
	result, _ := json.Marshal(m)
	return string(result)
}

func SetGroupBan(groupID int64, userID int64, duration int32) string {

	m := sendJSON{
		Action: "set_group_ban",
		Params: sendParams{
			GroupID:  groupID,
			UserID:   userID,
			Duration: duration,
		},
	}
	result, _ := json.Marshal(m)
	return string(result)
}

func SetGroupAnonymousBan(groupID int64, flag string, duration int32) string {

	m := sendJSON{
		Action: "set_group_anonymous_ban",
		Params: sendParams{
			GroupID:  groupID,
			Flag:     flag,
			Duration: duration,
		},
	}
	result, _ := json.Marshal(m)
	return string(result)
}

func SetGroupWholeBan(groupID int64, enable bool) string {

	m := sendJSON{
		Action: "set_group_whole_ban",
		Params: sendParams{
			GroupID: groupID,
			Enable:  enable,
		},
	}
	result, _ := json.Marshal(m)
	return string(result)
}

func GetGroupList() string {
	m := sendJSON{
		Action: "get_group_list",
		Params: sendParams{},
	}
	result, _ := json.Marshal(m)
	return string(result)

}

func GetGroupMemberList(group_id int64) string {
	m := sendJSON{
		Action: "get_group_member_list",
		Params: sendParams{
			GroupID: group_id,
		},
	}
	result, _ := json.Marshal(m)
	return string(result)

}

type EventJSON struct {
	Type       string      `json:"post_type"`
	NoticeType string      `json:"notice_type"`
	ReqType    string      `json:"request_type"`
	Comment    string      `json:"comment"`
	Flag       string      `json:"flag"`
	MsgType    string      `json:"message_type"`
	SubType    string      `json:"sub_type"`
	GroupID    int64       `json:"group_id"`
	DiscussID  int64       `json:"discuss_id"`
	MsgID      int32       `json:"message_id"`
	UserID     int64       `json:"user_id"`
	Msg        interface{} `json:"message"`
	RawMsg     string      `json:"raw_message"`
	Font       int32       `json:"font"`
	OperaID    int64       `json:"operator_id"`
	Self       int64       `json:"self_id"`
	Sender     struct {
		UserID int64  `json:"user_id"`
		Nick   string `json:"nickname"`
		Card   string `json:"card"`
		Sex    string `json:"sex"`
		Age    int32  `json:"age"`
	} `json:"sender"`
	Anonymous struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Flag string `json:"flag"`
	} `json:"anonymous"`
	File struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Size  int64  `json:"size"`
		Busid int64  `json:"busid"`
	} `json:"file"`

	// 添加事件
}
