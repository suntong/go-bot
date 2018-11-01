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
func (m *message) AddMsg(msg interface{}) {

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
}

type alone struct {
	Type string      `json:"type“`
	Data interface{} `json:"data"`
}

type more struct {
	Type string      `json:"type“`
	Data interface{} `json:"data"`
}

func AtQQ(qq string) more {
	return more{
		Type: "at",
		Data: struct {
			QQ string `json:"qq"`
		}{QQ: qq},
	}
}
