package utils

import (
	"testing"
)

func TestAddMsg(t *testing.T) {
	m := NewMessage()
	m.AddMsg(AtQQ("1047439649"))
	m.AddMsg(AtQQ("1047439649"))
	t.Log(m.Message())
}
