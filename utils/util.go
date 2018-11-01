package utils

func NewMessage() *message {
	return &message{}
}

type message struct {
	data interface{}
}

func (m *message) Message() interface{} {
	return m.data
}
func (m *message) AddMsg(msg interface{}) *message {

	if v, ok := msg.(alone); ok {
		m.data = v
	} else if m.data != nil {
		s, ok := m.data.([]interface{})
		if !ok {
			m.data = []interface{}{m.data, msg}
		} else {
			m.data = append(s, msg)
		}
	} else {
		m.data = []interface{}{msg}
	}
	return m
}

type alone struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type more struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func CQtext(text string) more {
	return more{
		Type: "text",
		Data: struct {
			Text string `json:"text"`
		}{Text: text},
	}
}

func CQat(qq string) more {
	return more{
		Type: "at",
		Data: struct {
			QQ string `json:"qq"`
		}{QQ: qq},
	}
}

func CQface(id int32) more {
	return more{
		Type: "face",
		Data: struct {
			ID int32 `json:"id"`
		}{ID: id},
	}
}

func CQemoji(id int32) more {
	return more{
		Type: "emoji",
		Data: struct {
			ID int32 `json:"id"`
		}{ID: id},
	}
}

func CQbface(id int32) more {
	return more{
		Type: "bface",
		Data: struct {
			ID int32 `json:"id"`
		}{ID: id},
	}
}

func CQsface(id int32) more {
	return more{
		Type: "sface",
		Data: struct {
			ID int32 `json:"id"`
		}{ID: id},
	}
}

func CQimage(url string) more {
	return more{
		Type: "image",
		Data: struct {
			File string `json:"file"`
		}{File: url},
	}
}

func CQrecord(url string, magic bool) more {
	return more{
		Type: "record",
		Data: struct {
			File  string `json:"file"`
			Magic bool   `json:"magic"`
		}{File: url, Magic: magic},
	}
}

func CQrps() more {
	return more{
		Type: "rps",
		Data: struct {
		}{},
	}
}

func CQdice() more {
	return more{
		Type: "dice",
		Data: struct {
		}{},
	}
}

func CQshake() more {
	return more{
		Type: "shake",
		Data: struct {
		}{},
	}
}

func CQanonymous(ignore bool) more {
	return more{
		Type: "anonymous",
		Data: struct {
			Ignore bool `json:"ignore"`
		}{Ignore: ignore},
	}
}

func CQSendmusic(t string, id int64) alone {
	return alone{
		Type: "music",
		Data: struct {
			Type string `json:"type"`
			ID   int64  `json:"id"`
		}{Type: t, ID: id},
	}
}

func CQSharemusic(url string, audio string, title string, content string, image string) alone {
	return alone{
		Type: "music",
		Data: struct {
			Type    string `json:"type"`
			Url     string `json:"url"`
			Audio   string `json:"audio"`
			Title   string `json:"title"`
			Content string `json:"content"`
			Image   string `json:"image"`
		}{Type: "custom", Url: url, Audio: audio, Title: title, Content: content, Image: image},
	}
}

func CQshare(url string, title string, content string, image string) alone {
	return alone{
		Type: "share",
		Data: struct {
			Url     string `json:"url"`
			Title   string `json:"title"`
			Content string `json:"content"`
			Image   string `json:"image"`
		}{Url: url, Title: title, Content: content, Image: image},
	}
}
