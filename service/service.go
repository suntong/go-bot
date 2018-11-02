package service

import (
	"go-bot/handle"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

var mux sync.WaitGroup

func LoadService(addr string) {
	mux.Add(2)
	if err := eventService(addr); err != nil {
		// 记入
	}

	if err := apiService(addr); err != nil {
		// 记入
	}

	mux.Wait()
}

func eventService(addr string) error {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/event"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	// 启动服务
	go func(c *websocket.Conn) {
		defer c.Close()
		for {
			_, b, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = handle.Handle(b)
			if err != nil {
				break
			}
		}
		mux.Done()
	}(c)
	return nil
}

func apiService(addr string) error {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/api"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	// 启动服务
	go func(c *websocket.Conn) {
		defer c.Close()
		for {
			if result := handle.Send(); result != nil {
				err := c.WriteJSON(result)
				if err != nil {
					break
				}
			}
		}
		mux.Done()
	}(c)
	return nil
}
