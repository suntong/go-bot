package service

import (
	"go-bot/handle"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lexkong/log"
)

var mux sync.WaitGroup

func LoadService(addr string) {
	mux.Add(3)
	if err := eventService(addr); err != nil {
		log.Error("service", err)
	}

	if err := apiService(addr); err != nil {
		log.Error("service", err)
	}
	if err := inputService(); err != nil {
		log.Error("service", err)
	}

	mux.Wait()
}

func inputService() error {
	http.HandleFunc("/monitoring", func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		log.Infof("option=[%s] ip=[%s] path=[%s] body=[%s] parm=[%s] header=[%s]", r.Method, r.RemoteAddr, r.RequestURI, string(bytes), r.PostForm, r.Header)
		w.WriteHeader(http.StatusOK)
	})
	return http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
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
