package handle

import (
	"encoding/json"
)

var testC = make(chan interface{}, 100)

func Handle(bytes []byte) error {
	// 返回错误关闭
	var e eventJSON
	if err := json.Unmarshal(bytes, &e); err != nil {
		return err
	}
	return nil
}

func Send() interface{} {
	return nil
}

type eventJSON struct {
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
	Anonymous  struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Flag string `json:"flag"`
	} `json:"anonymous"`
}
