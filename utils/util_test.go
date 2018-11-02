package utils

import (
	"encoding/json"
	"testing"
)

func TestAddMsg(t *testing.T) {
	m := NewMessage()
	m.AddMsg(CQat("1047439649"))
	m.AddMsg(CQat("1047439649"))
	d, _ := json.Marshal(m.Message())
	s := string(d)
	t.Log(s)
}

func TestFransferred(t *testing.T) {
	result := Fransferred("ss[QC:file,url:https://juejin.im]asdzsczasd[QC:file,url:https://juejin.im]zzzzz")
	if result != "ssasdzsczasdzzzzz" {
		t.Error(result)
	}
}
