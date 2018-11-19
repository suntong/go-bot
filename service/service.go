package service

import (
	"go-bot/handle"
	"net/url"
	"sync"

	"github.com/lexkong/log"

	"github.com/gorilla/websocket"
)

var mux sync.WaitGroup

func LoadService(addr string) {
	mux.Add(2)
	if err := eventService(addr); err != nil {
		log.Error("service", err)
	}

	if err := apiService(addr); err != nil {
		log.Error("service", err)
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
				log.Error("read message", err)
				break
			}
			err = handle.Handle(b)
			if err != nil {
				log.Error("handle", err)
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
					log.Error("write JSON", err)
					break
				}
			}

			// 响应
			_, b, err := c.ReadMessage()
			if err != nil {
				log.Error("read message", err)
				break
			}
			err = handle.Handle(b)
			if err != nil {
				log.Error("handle", err)
				break
			}
		}
		mux.Done()
	}(c)
	return nil
}
